#!/usr/bin/env python3
"""EdgePulse: concurrent endpoint latency and throughput measurement."""
from __future__ import annotations

import argparse
import csv
import ipaddress
import socket
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, TextIO


@dataclass
class Endpoint:
    address: str
    sent: int = 0
    received: int = 0
    loss: float = 1.0
    latency_ms: float = 0.0
    speed_mbps: float = 0.0
    region: str = "N/A"
    error: str = ""


def parse_targets(value: str) -> list[ipaddress._BaseNetwork | ipaddress._BaseAddress]:
    result = []
    for item in value.split(","):
        item = item.strip()
        if not item:
            continue
        try:
            result.append(ipaddress.ip_network(item, strict=False) if "/" in item else ipaddress.ip_address(item))
        except ValueError as exc:
            raise ValueError(f"invalid target '{item}': {exc}") from exc
    return result


def load_targets(path: str | Path) -> list[ipaddress._BaseNetwork | ipaddress._BaseAddress]:
    values = []
    for raw in Path(path).read_text(encoding="utf-8").splitlines():
        line = raw.split("#", 1)[0].strip()
        if line:
            values.extend(parse_targets(line))
    return values


def expand_targets(values: Iterable, all_addresses: bool, limit: int = 4096) -> list[str]:
    addresses: list[str] = []
    for value in values:
        if isinstance(value, (ipaddress.IPv4Address, ipaddress.IPv6Address)):
            addresses.append(str(value))
            continue
        hosts = list(value.hosts()) if all_addresses else [next(value.hosts(), value.network_address)]
        addresses.extend(str(host) for host in hosts)
        if len(addresses) >= limit:
            break
    return addresses[:limit]


def _socket_family(address: str) -> int:
    return socket.AF_INET6 if ":" in address else socket.AF_INET


def tcp_probe(address: str, port: int, attempts: int, timeout: float) -> Endpoint:
    row = Endpoint(address=address, sent=attempts)
    samples = []
    for _ in range(attempts):
        started = time.perf_counter()
        try:
            with socket.socket(_socket_family(address), socket.SOCK_STREAM) as sock:
                sock.settimeout(timeout)
                sock.connect((address, port))
            samples.append((time.perf_counter() - started) * 1000)
        except OSError as exc:
            row.error = str(exc)
    row.received = len(samples)
    row.loss = (attempts - row.received) / attempts if attempts else 1.0
    row.latency_ms = sum(samples) / len(samples) if samples else 0.0
    return row


def http_probe(address: str, url: str, port: int, attempts: int, timeout: float) -> Endpoint:
    row = Endpoint(address=address, sent=attempts)
    host = urllib.parse.urlsplit(url).hostname if "//" in url else None
    if not host:
        row.error = "invalid URL"
        return row
    samples = []
    for _ in range(attempts):
        started = time.perf_counter()
        try:
            request = urllib.request.Request(url, method="HEAD", headers={"Host": host})
            opener = _pinned_opener(address, url, port)
            with opener.open(request, timeout=timeout) as response:
                if 200 <= response.status < 400:
                    samples.append((time.perf_counter() - started) * 1000)
                    row.region = response.headers.get("X-Region", response.headers.get("CF-Ray", "N/A"))
                else:
                    row.error = f"HTTP {response.status}"
        except (OSError, urllib.error.URLError) as exc:
            row.error = str(exc.reason if isinstance(exc, urllib.error.URLError) else exc)
    row.received = len(samples)
    row.loss = (attempts - row.received) / attempts if attempts else 1.0
    row.latency_ms = sum(samples) / len(samples) if samples else 0.0
    return row


class _PinnedAddressHandler(urllib.request.HTTPHandler):
    def __init__(self, address: str, port: int, secure: bool):
        super().__init__()
        self.address, self.port, self.secure = address, port, secure

    def _open(self, request):
        sock = socket.socket(_socket_family(self.address), socket.SOCK_STREAM)
        sock.settimeout(request.timeout if hasattr(request, "timeout") else None)
        sock.connect((self.address, self.port))
        if self.secure:
            context = ssl.create_default_context()
            sock = context.wrap_socket(sock, server_hostname=urllib.parse.urlsplit(request.full_url).hostname)
        return sock

    def http_open(self, request):
        request.host = urllib.parse.urlsplit(request.full_url).netloc
        return self.do_open(lambda host, timeout=None, **kwargs: self._open(request), request)


class _PinnedHTTPSHandler(urllib.request.HTTPSHandler):
    def __init__(self, address: str, port: int):
        super().__init__()
        self.address, self.port = address, port

    def https_open(self, request):
        request.host = urllib.parse.urlsplit(request.full_url).netloc
        sock = socket.socket(_socket_family(self.address), socket.SOCK_STREAM)
        sock.settimeout(getattr(request, "timeout", None))
        sock.connect((self.address, self.port))
        context = ssl.create_default_context()
        sock = context.wrap_socket(sock, server_hostname=urllib.parse.urlsplit(request.full_url).hostname)
        return self.do_open(lambda host, timeout=None, **kwargs: sock, request)


def _pinned_opener(address: str, url: str, port: int):
    secure = urllib.parse.urlsplit(url).scheme.lower() == "https"
    handler = _PinnedHTTPSHandler(address, port) if secure else _PinnedAddressHandler(address, port, False)
    return urllib.request.build_opener(handler)


def download_probe(row: Endpoint, url: str, port: int, timeout: float) -> Endpoint:
    started = time.perf_counter()
    total = 0
    try:
        request = urllib.request.Request(url, headers={"Range": "bytes=0-10485759"})
        with _pinned_opener(row.address, url, port).open(request, timeout=timeout) as response:
            row.region = response.headers.get("X-Region", response.headers.get("CF-Ray", row.region))
            while total < 10 * 1024 * 1024:
                chunk = response.read(65536)
                if not chunk:
                    break
                total += len(chunk)
    except (OSError, urllib.error.URLError) as exc:
        row.error = str(exc.reason if isinstance(exc, urllib.error.URLError) else exc)
    elapsed = max(time.perf_counter() - started, 0.001)
    row.speed_mbps = total / elapsed / 1024 / 1024
    return row


def format_table(rows: list[Endpoint]) -> str:
    headers = ["Address", "Sent", "Received", "Loss %", "Latency (ms)", "Speed (MB/s)", "Region"]
    lines = ["  ".join(headers), "  ".join("-" * len(h) for h in headers)]
    lines.extend("  ".join([r.address, str(r.sent), str(r.received), f"{r.loss * 100:.2f}", f"{r.latency_ms:.2f}", f"{r.speed_mbps:.2f}", r.region]) for r in rows)
    return "\n".join(lines)


def write_csv(rows: Iterable[Endpoint], target: TextIO | str | Path) -> None:
    close = False
    if hasattr(target, "write"):
        stream = target
    else:
        stream = open(target, "w", newline="", encoding="utf-8")
        close = True
    try:
        writer = csv.writer(stream)
        writer.writerow(["Address", "Sent", "Received", "Loss %", "Latency (ms)", "Speed (MB/s)", "Region"])
        for row in rows:
            writer.writerow([row.address, row.sent, row.received, f"{row.loss * 100:.2f}", f"{row.latency_ms:.2f}", f"{row.speed_mbps:.2f}", row.region])
    finally:
        if close:
            stream.close()


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="edgepulse", description="Measure endpoint latency and download throughput.")
    parser.add_argument("-f", "--file", default="ip.txt", help="Read addresses and CIDR ranges from this file.")
    parser.add_argument("-i", "--ip", help="Comma-separated addresses and CIDR ranges.")
    parser.add_argument("-p", "--port", type=int, default=443, help="TCP port (default: 443).")
    parser.add_argument("-n", "--workers", type=int, default=32, help="Concurrent workers (default: 32).")
    parser.add_argument("-t", "--attempts", type=int, default=3, help="Probes per endpoint (default: 3).")
    parser.add_argument("--timeout", type=float, default=3.0, help="Timeout in seconds (default: 3).")
    parser.add_argument("--http", action="store_true", help="Use HTTP HEAD probes instead of TCP probes.")
    parser.add_argument("--url", default="https://example.com/", help="URL used by HTTP and download probes.")
    parser.add_argument("--download", action="store_true", help="Measure download throughput after latency probes.")
    parser.add_argument("--all", action="store_true", help="Expand every address in each CIDR range.")
    parser.add_argument("--max-latency", type=float, default=float("inf"), help="Keep results at or below this latency in ms.")
    parser.add_argument("--max-loss", type=float, default=1.0, help="Keep results at or below this loss ratio.")
    parser.add_argument("--min-speed", type=float, default=0.0, help="Keep results at or above this speed in MB/s.")
    parser.add_argument("-o", "--output", default="result.csv", help="CSV output path; use an empty value to disable.")
    parser.add_argument("--version", action="version", version="EdgePulse 1.0.0")
    return parser


def main(argv: list[str] | None = None) -> int:
    args = build_parser().parse_args(argv)
    try:
        values = parse_targets(args.ip) if args.ip else load_targets(args.file)
        addresses = expand_targets(values, args.all)
    except (OSError, ValueError) as exc:
        print(f"Error: {exc}", file=sys.stderr)
        return 2
    if not addresses:
        print("Error: no valid targets were found.", file=sys.stderr)
        return 2
    probe = (lambda address: http_probe(address, args.url, args.port, args.attempts, args.timeout)) if args.http else (lambda address: tcp_probe(address, args.port, args.attempts, args.timeout))
    print(f"Testing {len(addresses)} endpoints with {args.workers} workers...")
    rows = []
    with ThreadPoolExecutor(max_workers=max(1, args.workers)) as pool:
        futures = {pool.submit(probe, address): address for address in addresses}
        for future in as_completed(futures):
            rows.append(future.result())
    rows = [r for r in rows if r.received and r.latency_ms <= args.max_latency and r.loss <= args.max_loss]
    if args.download:
        with ThreadPoolExecutor(max_workers=max(1, min(args.workers, 16))) as pool:
            rows = list(pool.map(lambda row: download_probe(row, args.url, args.port, args.timeout), rows))
        rows = [r for r in rows if r.speed_mbps >= args.min_speed]
    rows.sort(key=lambda r: (-r.speed_mbps, r.latency_ms) if args.download else (r.latency_ms, r.loss))
    print(format_table(rows[:20]) if rows else "No endpoint matched the selected criteria.")
    if args.output:
        write_csv(rows, args.output)
        print(f"Full results written to {args.output}.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
