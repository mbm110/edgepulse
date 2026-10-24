# EdgePulse CDN IP Scanner and Latency / Speed Tester

[English](README.md) | [فارسی](README.fa.md)

EdgePulse is a cross-platform Go command-line tool for scanning authorized IPv4 and IPv6 addresses and CIDR ranges. It measures TCP or HTTP latency, calculates packet loss, tests download throughput, filters candidates, reads CDN region codes, and exports complete CSV results.

The terminal and all program messages are English only.

## Capabilities

- Concurrent TCP latency testing.
- HTTP HEAD latency testing with accepted status codes.
- IPv4, IPv6, mixed input, CIDR ranges, and individual addresses.
- One representative address per range or full IPv4 range testing.
- Download speed testing with configurable URL and duration.
- Latency, loss, speed, HTTP status, and region filters.
- Region detection for common CDN response headers.
- CSV output with English headers.
- Optional hosts, DNS, dnsmasq, and proxy helper scripts.
- Static Linux ARM64 release for Termux-compatible Linux environments.

## Quick start

Download and extract a release, then run the binary from its directory:

```bash
./edgepulse
```

Test explicit addresses or ranges:

```bash
./edgepulse -ip 1.1.1.1,1.0.0.1,2606:4700::/32
```

Filter candidates and test download speed:

```bash
./edgepulse -tl 200 -tlr 0.25 -sl 1 -dn 10
```

View every option:

```bash
./edgepulse -h
```

## Workflow

1. Read ranges from `ip.txt` or `-ip`.
2. Select one address from each range, or all IPv4 addresses with `-allip`.
3. Measure latency concurrently using TCP by default.
4. Use HTTP requests and accepted status codes with `-httping`.
5. Filter by latency, packet loss, and region.
6. Measure download speed unless `-dd` is supplied.
7. Sort results and write `result.csv`.

## Options

```text
-n 200          latency worker count (default 200, maximum 1000)
-t 4            latency probes per endpoint
-dn 10          number of download candidates
-dt 10          maximum download time per endpoint in seconds
-tp 443         test port
-url URL        URL used for HTTP and download tests
-httping        use HTTP latency mode
-httping-code N accepted HTTP status code in HTTP mode
-cfcolo CODES   filter region codes, comma-separated
-tl MS          maximum average latency
-tll MS         minimum average latency
-tlr RATIO      maximum loss ratio from 0.00 to 1.00
-sl MBPS        minimum download speed
-p N            number of displayed rows; zero hides the table
-f FILE         IP range input file (default ip.txt)
-ip RANGES      comma-separated addresses or CIDR ranges
-o FILE         CSV output path; empty disables file output
-dd             disable download testing
-allip          test every IPv4 address in each range
-debug          print diagnostic details
-v              print version and check for updates
-h              print help
```

## TCP and HTTP modes

TCP mode checks whether a connection can be opened to the selected port. HTTP mode checks the selected hostname and URL through each candidate IP and can read a region code from response headers.

```bash
./edgepulse -httping -url https://example.com/
./edgepulse -httping -tp 80 -url http://example.com/
```

HTTP mode is network scanning. Lower `-n` on a VPS and follow the destination service policy. Test only networks and URLs you are authorized to measure.

## IPv4, IPv6, and ranges

The input file can mix address families:

```text
1.1.1.1
1.0.0.1/24
2606:4700::/32
```

One address is accepted without a mask. By default, one address is selected from each range. Use `-allip` only for authorized and manageable ranges.

## Region filtering

HTTP mode reads region data from response headers when the CDN provides it. Depending on the provider, the code may be an IATA airport code, country code, or city code:

```bash
./edgepulse -httping -cfcolo FRA,LAX,SJC
```

Region routing can change by network, provider, and time. Treat it as a measurement from the current machine.

## Output

The terminal uses English headers:

```text
Address           Sent  Received  Loss %  Average Latency  Download Speed (MB/s)  Region
```

The complete result is saved in `result.csv`. Use `-o ""` to disable file output.

## Helper scripts

The `script/` directory contains optional examples for hosts, DNS, dnsmasq, and proxy updates. Review scripts before use and keep API credentials outside source control. The main program does not modify hosts files or DNS records automatically.

## Termux and Android

Use the static `edgepulse-linux-arm64.tar.gz` asset on compatible Linux ARM64 environments. For Termux installations where a Linux binary is not compatible with Android's userspace, use the source archive and build it with the Termux Go toolchain:

```bash
pkg update
pkg install golang git
git clone https://github.com/mbm110/edgepulse.git
cd edgepulse
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o edgepulse main.go
./edgepulse -h
```

## Releases

GitHub Actions builds Linux AMD64, static Linux ARM64, Windows AMD64, macOS Intel, macOS Apple Silicon, and a source archive. Linux ARM64 is built with `CGO_ENABLED=0` and `-ldflags="-s -w"`.

## Donate

- **Ethereum:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **Bitcoin:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **Tron:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **Solana:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

## License

EdgePulse is distributed under GPLv3. See [LICENSE](LICENSE).
