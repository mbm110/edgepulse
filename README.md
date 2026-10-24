# EdgePulse IP Range Scanner & CDN Latency and Speed Tester

EdgePulse is a cross-platform command-line utility for scanning authorized IPv4 and IPv6 addresses and CIDR ranges, comparing endpoint quality, and finding responsive low-latency candidates. It measures TCP reachability, can send HTTP `HEAD` probes, calculates packet loss and average latency, optionally measures download throughput, filters weak endpoints, sorts the survivors, and exports a complete CSV report.

The executable and all terminal messages are English only. This README contains complete English and Persian documentation.

> **Authorized use:** scan only addresses, networks, and URLs that you own or are explicitly permitted to test. HTTP mode sends requests to the selected URL, and a large range scan can create noticeable traffic.

## What EdgePulse does

- Scans a single IP address, a list of addresses, or CIDR ranges.
- Scans IPv4 and IPv6 targets in one run.
- Uses one representative address per range by default.
- Expands a range with `--all` when a full scan is authorized.
- Runs concurrent TCP connection tests.
- Runs HTTP `HEAD` checks against a selected URL.
- Counts sent probes, successful probes, and packet loss.
- Calculates average latency in milliseconds.
- Measures download throughput for reachable endpoints.
- Filters by maximum latency, maximum loss, and minimum speed.
- Sorts by latency or download speed.
- Prints the strongest results in the terminal.
- Saves the complete filtered dataset as UTF-8 CSV.
- Runs on Linux, Windows, macOS, and other systems with Python 3.10+.

EdgePulse is a network measurement and endpoint selection tool. It does not exploit systems, bypass authentication, or perform credential attacks.

## Why scan IP ranges?

A CDN or distributed service may publish several possible IP addresses. The best endpoint depends on your ISP, geography, routing, congestion, and time of day. An address that performs well for one network can be slow or unreachable from another.

Scanning from the network that will actually use the service gives you a practical comparison of the available endpoints. Results describe the route at test time and are not a permanent guarantee that one address will always remain fastest.

## Installation

### Release archive

Download an archive from [Releases](https://github.com/mbm110/edgepulse/releases), extract it, and run the included program. Release packages are created only by GitHub Actions.

### Run from source

```bash
git clone https://github.com/mbm110/edgepulse.git
cd edgepulse
python3 edgepulse.py --help
```

EdgePulse uses only the Python standard library and requires no third-party package.

## Quick start

### Scan individual addresses

```bash
python3 edgepulse.py --ip 1.1.1.1,1.0.0.1
```

### Scan every address in a small CIDR range

```bash
python3 edgepulse.py --ip 1.1.1.0/29 --all
```

### Read targets from a file

Put one address or CIDR range on each line of `ip.txt`:

```text
1.1.1.1
1.0.0.1
2606:4700:4700::1111
2606:4700::/120
```

Run the scan:

```bash
python3 edgepulse.py --file ip.txt
```

By default, EdgePulse selects one representative host from each range. Use `--all` for every host. Range expansion is capped at 4096 addresses per run to reduce accidental traffic and protect small VPS machines.

### Apply latency and loss filters

```bash
python3 edgepulse.py --file ip.txt --max-latency 200 --max-loss 0.25
```

This keeps endpoints with average latency up to 200 ms and packet loss up to 25 percent.

### Measure download throughput

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/test.bin \
  --min-speed 1.0
```

Use a large test file or a dedicated endpoint that you control. A tiny file may finish before a stable throughput sample is collected.

## How a scan works

1. EdgePulse reads individual addresses and CIDR ranges.
2. It expands ranges according to the selected mode.
3. It runs concurrent TCP probes, or HTTP `HEAD` probes when `--http` is enabled.
4. It removes unreachable targets and applies latency and loss filters.
5. With `--download`, it reads up to 10 MiB from each remaining endpoint.
6. It applies the minimum speed filter.
7. It sorts the survivors.
8. It prints a terminal table and writes `result.csv`.

The terminal shows up to 20 rows. The CSV contains every row remaining after filtering.

## Command-line options

Run `python3 edgepulse.py --help` for the live list.

| Option | Purpose | Default |
|---|---|---:|
| `-f`, `--file FILE` | Read addresses and CIDR ranges from a file. | `ip.txt` |
| `-i`, `--ip VALUE` | Comma-separated addresses and ranges. | empty |
| `-p`, `--port PORT` | TCP service port. | `443` |
| `-n`, `--workers N` | Concurrent workers. | `32` |
| `-t`, `--attempts N` | Probes per endpoint. | `3` |
| `--timeout SECONDS` | Timeout for one connection or request. | `3` |
| `--http` | Use HTTP `HEAD` probes instead of raw TCP. | off |
| `--url URL` | URL used for HTTP and download checks. | `https://example.com/` |
| `--download` | Measure throughput after reachability tests. | off |
| `--all` | Test every host in each CIDR range. | off |
| `--max-latency MS` | Keep results at or below this average latency. | unlimited |
| `--max-loss RATIO` | Keep results at or below a ratio from `0` to `1`. | `1.0` |
| `--min-speed MB/S` | Keep download results at or above this speed. | `0` |
| `-o`, `--output FILE` | Write the complete CSV report. | `result.csv` |
| `--version` | Print the program version. | — |

## TCP mode and HTTP mode

### TCP mode

TCP mode is the default. It opens a TCP connection to the selected port and closes it after success. It is lightweight and does not require a URL:

```bash
python3 edgepulse.py --ip 1.1.1.1 --port 443
```

TCP success proves that a connection to the address and port was possible from the current machine. It does not prove that a hostname, certificate, HTTP route, or application will work through that address.

### HTTP mode

HTTP mode sends `HEAD` requests to the URL supplied with `--url`:

```bash
python3 edgepulse.py \
  --ip 1.1.1.1,1.0.0.1 \
  --http \
  --port 443 \
  --url https://example.com/
```

HTTP mode is a form of network scanning. Lower `--workers`, use a reasonable timeout, and respect the policy of the service being tested.

## IPv4, IPv6, and CIDR input

IPv4 and IPv6 can be mixed in one command:

```bash
python3 edgepulse.py \
  --ip 1.1.1.1,2606:4700:4700::1111,2001:db8::/120
```

A single address does not need `/32` or `/128`. A CIDR range is fully expanded only with `--all`, which prevents an unexpectedly large scan.

Comments beginning with `#` and blank lines are ignored:

```text
# IPv4 candidates
1.1.1.0/29

# IPv6 candidates
2606:4700::/120
```

If every IPv6 target fails while IPv4 works, check the machine's IPv6 address, route, firewall, and provider support.

## Filtering and sorting

### Latency ceiling

```bash
python3 edgepulse.py --file ip.txt --max-latency 200
```

Only endpoints at or below 200 ms remain.

### Loss ceiling

The loss value is a ratio. `0` means no loss and `0.25` means up to 25 percent loss:

```bash
python3 edgepulse.py --file ip.txt --max-loss 0
python3 edgepulse.py --file ip.txt --max-loss 0.25
```

### Minimum speed

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/large-file.bin \
  --min-speed 5
```

If no endpoint reaches the requested speed, first run without `--min-speed` and choose a realistic threshold from the measured range.

### Combined filters

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/large-file.bin \
  --max-latency 200 \
  --max-loss 0.10 \
  --min-speed 2
```

Every selected filter must pass for an endpoint to remain in the final output.

Without download mode, results are sorted by latency and then loss. With download mode, speed is the primary sort key and latency is secondary.

## Understanding the output

```text
Address  Sent  Received  Loss %  Latency (ms)  Speed (MB/s)  Region
-------  ----  --------  ------  ------------  ------------  ------
1.1.1.1  3     3         0.00    18.42         0.00          N/A
```

- **Address:** tested IPv4 or IPv6 address.
- **Sent:** probes attempted.
- **Received:** successful probes.
- **Loss %:** failed probes as a percentage.
- **Latency (ms):** average successful response time.
- **Speed (MB/s):** throughput when `--download` is enabled.
- **Region:** value read from response headers when available; otherwise `N/A`.

## CSV output

The default report is `result.csv`:

```csv
Address,Sent,Received,Loss %,Latency (ms),Speed (MB/s),Region
1.1.1.1,3,3,0.00,18.42,0.00,N/A
```

Choose another path:

```bash
python3 edgepulse.py --file ip.txt --output reports/edgepulse.csv
```

Disable file output:

```bash
python3 edgepulse.py --ip 1.1.1.1 --output ""
```

The CSV is UTF-8 with English headers and no Chinese or Persian runtime text.

## Custom ports and URLs

For a non-default port:

```bash
python3 edgepulse.py --ip 203.0.113.10 --port 8443
```

For HTTP on port 80, select an `http://` URL that does not force a redirect to HTTPS:

```bash
python3 edgepulse.py \
  --ip 203.0.113.10 \
  --http \
  --port 80 \
  --url http://example.com/
```

For other ports, confirm that the selected service accepts traffic on that port. A hostname and IP should belong to the same service; mixing unrelated CDN addresses commonly causes certificate or HTTP failures.

## VPS and resource guidance

EdgePulse is small and dependency-free. A conservative VPS command is:

```bash
python3 edgepulse.py \
  --file ip.txt \
  --workers 8 \
  --attempts 2 \
  --timeout 5 \
  --output result.csv
```

The default worker count is 32. On a small VPS, use 8 or 16. More workers consume more sockets and may trigger remote rate limits. Avoid `--all` on large ranges unless the scan is authorized and planned.

## Troubleshooting

### No valid targets were found

Check that the input file exists and each line contains a valid IPv4 address, IPv6 address, or CIDR range:

```bash
python3 edgepulse.py --ip 1.1.1.1
```

### Every target fails

Possible causes include a closed port, outbound firewall rules, missing IPv4 or IPv6 routing, a short timeout, or IPs that belong to a different service. Try:

```bash
python3 edgepulse.py --ip 1.1.1.1 --timeout 8 --workers 4
```

### HTTP mode returns no results

Confirm the URL, hostname, port, and ownership of the tested IPs. Test the URL first with a browser or `curl -I`.

### Download speed is zero

Possible causes are an unavailable URL, a file that is too small, an IP that does not serve the hostname, direct-IP blocking, certificate mismatch, provider limits, or a timeout. Use a permitted large test file and a longer timeout. Do not disable certificate verification to hide a mismatch.

### Results differ between runs

This is normal. Routing, congestion, load balancing, and CDN selection change. Repeat the test and compare CSV files before making a long-term routing decision.

## Windows, Linux, and macOS

PowerShell:

```powershell
py edgepulse.py --ip 1.1.1.1,1.0.0.1
py edgepulse.py --file ip.txt --download --url https://example.com/test.bin
```

Linux and macOS:

```bash
python3 edgepulse.py --help
python3 edgepulse.py --file ip.txt
```

All terminal output is English only on every supported platform.

## Build and release

Packages are built only by GitHub Actions. The workflow runs tests, checks that runtime files contain no CJK text, creates lightweight archives, and publishes assets when a `v*.*.*` tag is pushed.

Local validation:

```bash
python3 -m unittest discover -s tests -v
python3 -m py_compile edgepulse.py
python3 edgepulse.py --help
```

## Brand and logo

The new project logo is [EDGEPLUSE_LOGO.svg](EDGEPLUSE_LOGO.svg). The blue and cyan pulse mark represents a network signal crossing multiple endpoints.

## Donate

- **Ethereum:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **Bitcoin:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **Tron:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **Solana:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

## License

EdgePulse is distributed under the MIT License. See [LICENSE](LICENSE).

---

# EdgePulse؛ اسکنر بازه‌های IP و ابزار سنجش تأخیر و سرعت CDN

ابزار **EdgePulse** یک برنامه خط فرمان چندسکویی برای اسکن نشانی‌های IPv4 و IPv6 و بازه‌های CIDR است. برنامه کیفیت مسیر شبکه را مقایسه می‌کند، اتصال TCP را می‌سنجد، در حالت HTTP درخواست `HEAD` می‌فرستد، Packet Loss و میانگین تأخیر را محاسبه می‌کند، سرعت دریافت را می‌سنجد، مقصدهای ضعیف را فیلتر می‌کند و گزارش کامل CSV می‌سازد.

تمام پیام‌های اجرایی در Linux، Windows، macOS و سایر سیستم‌های پشتیبانی‌شده فقط انگلیسی هستند. این بخش، راهنمای کامل فارسی پروژه است.

> **استفاده مجاز:** فقط IPها، شبکه‌ها و URLهایی را اسکن کنید که مالک آن‌ها هستید یا اجازه بررسی‌شان را دارید. اسکن HTTP درخواست واقعی ارسال می‌کند و اسکن بازه بزرگ می‌تواند ترافیک قابل توجهی ایجاد کند.

## پروژه چه کاری انجام می‌دهد؟

- اسکن یک IP، چند IP یا بازه CIDR.
- اسکن هم‌زمان IPv4 و IPv6.
- انتخاب یک IP نماینده از هر بازه به‌صورت پیش‌فرض.
- اسکن کامل بازه با گزینه `--all`.
- اجرای آزمایش‌های هم‌زمان TCP.
- اجرای درخواست HTTP `HEAD` برای بررسی سرویس واقعی.
- شمارش Probeهای ارسال‌شده، موفق و ناموفق.
- محاسبه Packet Loss و میانگین تأخیر برحسب میلی‌ثانیه.
- سنجش Download Throughput برای مقصدهای قابل دسترس.
- فیلتر بر اساس Maximum Latency، Maximum Loss و Minimum Speed.
- مرتب‌سازی بر اساس تأخیر یا سرعت دانلود.
- نمایش بهترین نتایج در ترمینال.
- ذخیره همه نتایج فیلترشده در CSV با ستون‌های انگلیسی.
- اجرا روی Linux، Windows، macOS و سیستم‌های دارای Python 3.10 یا جدیدتر.

EdgePulse ابزار نفوذ، دورزدن احراز هویت یا حمله به Credentialها نیست؛ کاربرد آن اندازه‌گیری شبکه و انتخاب Endpoint مناسب است.

## چرا بازه‌های IP را اسکن کنیم؟

یک CDN یا سرویس توزیع‌شده ممکن است چندین IP در اختیار کاربر بگذارد. بهترین IP به ISP، موقعیت جغرافیایی، Routing، شلوغی شبکه و زمان آزمایش بستگی دارد. IP مناسب برای یک شبکه ممکن است از شبکه دیگر کند یا غیرقابل دسترس باشد.

اسکن از همان شبکه‌ای که قرار است سرویس را مصرف کند، مقایسه واقعی‌تری ارائه می‌دهد. نتیجه وضعیت مسیر در زمان آزمایش است و تضمین نمی‌کند یک IP همیشه سریع‌ترین باقی بماند.

## نصب و اجرا

### دریافت Release

از بخش [Releases](https://github.com/mbm110/edgepulse/releases) فایل مناسب را دریافت و Extract کنید. بسته‌ها فقط با GitHub Actions ساخته می‌شوند.

### اجرا از Source

```bash
git clone https://github.com/mbm110/edgepulse.git
cd edgepulse
python3 edgepulse.py --help
```

پروژه به Package خارجی Python نیاز ندارد و فقط از Standard Library استفاده می‌کند.

## شروع سریع

### اسکن چند IP

```bash
python3 edgepulse.py --ip 1.1.1.1,1.0.0.1
```

### اسکن کامل یک بازه کوچک

```bash
python3 edgepulse.py --ip 1.1.1.0/29 --all
```

### استفاده از فایل

هر IP یا CIDR را در یک خط فایل `ip.txt` قرار دهید:

```text
1.1.1.1
1.0.0.1
2606:4700:4700::1111
2606:4700::/120
```

سپس اجرا کنید:

```bash
python3 edgepulse.py --file ip.txt
```

به‌صورت پیش‌فرض از هر بازه یک Host انتخاب می‌شود. برای اسکن همه Hostها `--all` را اضافه کنید. برای کاهش مصرف VPS و جلوگیری از اسکن ناخواسته، بازکردن بازه در هر اجرا به ۴۰۹۶ نشانی محدود شده است.

### فیلتر تأخیر و Loss

```bash
python3 edgepulse.py --file ip.txt --max-latency 200 --max-loss 0.25
```

این دستور مقصدهایی را نگه می‌دارد که تأخیرشان حداکثر ۲۰۰ میلی‌ثانیه و Loss آن‌ها حداکثر ۲۵ درصد است.

### سنجش سرعت دانلود

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/test.bin \
  --min-speed 1.0
```

از فایل بزرگ یا Endpoint اختصاصی و مجاز استفاده کنید؛ فایل کوچک ممکن است نمونه پایداری از سرعت ارائه ندهد.

## روند کار برنامه

1. خواندن IPها و CIDRها.
2. بازکردن بازه‌ها بر اساس گزینه انتخابی.
3. اجرای Probeهای TCP یا HTTP.
4. حذف مقصدهای غیرقابل دسترس و اعمال فیلترها.
5. دریافت حداکثر ۱۰ مگابایت در حالت `--download`.
6. اعمال فیلتر Minimum Speed.
7. مرتب‌سازی نتایج.
8. نمایش جدول و نوشتن `result.csv`.

در ترمینال حداکثر ۲۰ ردیف نمایش داده می‌شود و CSV همه ردیف‌های نهایی را نگه می‌دارد.

## گزینه‌های خط فرمان

برای فهرست به‌روز گزینه‌ها دستور `python3 edgepulse.py --help` را اجرا کنید.

| گزینه | کاربرد | پیش‌فرض |
|---|---|---:|
| `-f` یا `--file FILE` | خواندن IP و CIDR از فایل | `ip.txt` |
| `-i` یا `--ip VALUE` | فهرست IP و CIDR با ویرگول | خالی |
| `-p` یا `--port PORT` | پورت TCP | `443` |
| `-n` یا `--workers N` | تعداد Worker هم‌زمان | `32` |
| `-t` یا `--attempts N` | تعداد Probe برای هر IP | `3` |
| `--timeout SECONDS` | Timeout هر اتصال یا درخواست | `3` |
| `--http` | استفاده از HTTP HEAD | خاموش |
| `--url URL` | URL آزمایش HTTP و دانلود | `https://example.com/` |
| `--download` | اندازه‌گیری Throughput | خاموش |
| `--all` | اسکن همه IPهای هر CIDR | خاموش |
| `--max-latency MS` | حداکثر تأخیر مجاز | بدون محدودیت |
| `--max-loss RATIO` | حداکثر نسبت Loss از صفر تا یک | `1.0` |
| `--min-speed MB/S` | حداقل سرعت دانلود | `0` |
| `-o` یا `--output FILE` | مسیر CSV | `result.csv` |
| `--version` | نمایش نسخه | — |

## حالت TCP و HTTP

### TCP

حالت پیش‌فرض یک اتصال TCP به پورت مشخص باز می‌کند و آن را می‌بندد:

```bash
python3 edgepulse.py --ip 1.1.1.1 --port 443
```

موفقیت TCP فقط نشان می‌دهد اتصال به IP و Port ممکن بوده است؛ درباره Hostname، Certificate یا سرویس HTTP تضمین نمی‌دهد.

### HTTP

حالت HTTP به URL تعیین‌شده با `--url` درخواست `HEAD` می‌فرستد:

```bash
python3 edgepulse.py \
  --ip 1.1.1.1,1.0.0.1 \
  --http \
  --port 443 \
  --url https://example.com/
```

HTTP نوعی Network Scanning است؛ تعداد Worker را پایین نگه دارید و قوانین سرویس مقصد را رعایت کنید.

## IPv4، IPv6 و CIDR

```bash
python3 edgepulse.py \
  --ip 1.1.1.1,2606:4700:4700::1111,2001:db8::/120
```

برای یک IP نیازی به `/32` یا `/128` نیست. CIDR فقط با `--all` کامل باز می‌شود تا اسکن بزرگ ناخواسته رخ ندهد.

خطوط خالی و توضیحاتی که با `#` شروع می‌شوند نادیده گرفته می‌شوند. اگر IPv6ها شکست می‌خورند ولی IPv4 کار می‌کند، Route، Firewall و پشتیبانی Provider را بررسی کنید.

## فیلترها و مرتب‌سازی

### حداکثر تأخیر

```bash
python3 edgepulse.py --file ip.txt --max-latency 200
```

### حداکثر Loss

صفر یعنی بدون Loss و `0.25` یعنی حداکثر ۲۵ درصد:

```bash
python3 edgepulse.py --file ip.txt --max-loss 0
python3 edgepulse.py --file ip.txt --max-loss 0.25
```

### حداقل سرعت

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/large-file.bin \
  --min-speed 5
```

اگر هیچ مقصدی به مقدار تعیین‌شده نرسید، ابتدا بدون `--min-speed` آزمایش کنید تا مقدار واقعی مشخص شود.

### ترکیب فیلترها

```bash
python3 edgepulse.py \
  --file ip.txt \
  --download \
  --url https://example.com/large-file.bin \
  --max-latency 200 \
  --max-loss 0.10 \
  --min-speed 2
```

همه فیلترهای انتخاب‌شده باید برای باقی‌ماندن مقصد برقرار باشند.

## خروجی و CSV

نمونه خروجی انگلیسی:

```text
Address  Sent  Received  Loss %  Latency (ms)  Speed (MB/s)  Region
-------  ----  --------  ------  ------------  ------------  ------
1.1.1.1  3     3         0.00    18.42         0.00          N/A
```

فایل پیش‌فرض `result.csv` است. برای تغییر مسیر:

```bash
python3 edgepulse.py --file ip.txt --output reports/edgepulse.csv
```

برای حذف فایل خروجی:

```bash
python3 edgepulse.py --ip 1.1.1.1 --output ""
```

CSV با UTF-8 و Headerهای انگلیسی نوشته می‌شود و در Runtime هیچ متن فارسی یا چینی تولید نمی‌شود.

## پورت و URL سفارشی

```bash
python3 edgepulse.py --ip 203.0.113.10 --port 8443
```

برای HTTP روی پورت ۸۰:

```bash
python3 edgepulse.py \
  --ip 203.0.113.10 \
  --http \
  --port 80 \
  --url http://example.com/
```

مطمئن شوید سرویس و URL انتخابی واقعاً پورت مورد نظر را پشتیبانی می‌کنند. ترکیب Hostname یک CDN با IP سرویس دیگر معمولاً خطای Certificate یا HTTP ایجاد می‌کند.

## رفع خطا

### مقصد معتبر پیدا نشد

فایل ورودی و قالب IP یا CIDR را بررسی کنید:

```bash
python3 edgepulse.py --ip 1.1.1.1
```

### همه مقصدها شکست می‌خورند

پورت بسته، Firewall، نبود Route، Timeout کم یا تعلق IPها به سرویس دیگر از علت‌های رایج هستند:

```bash
python3 edgepulse.py --ip 1.1.1.1 --timeout 8 --workers 4
```

### حالت HTTP نتیجه ندارد

URL، Hostname، Port و تعلق IPها به سرویس را بررسی کنید و URL را با `curl -I` آزمایش کنید.

### سرعت دانلود صفر است

URL نامعتبر، فایل کوچک، IP نامرتبط، جلوگیری از اتصال مستقیم، خطای Certificate، محدودیت Provider یا Timeout می‌تواند علت باشد. برای مخفی‌کردن ناسازگاری، Certificate Verification را خاموش نکنید.

### نتایج تغییر می‌کنند

این موضوع طبیعی است؛ Routing، Congestion، Load Balancing و انتخاب CDN تغییر می‌کنند. چند بار آزمایش و CSVها را مقایسه کنید.

## اجرای بهینه روی VPS

```bash
python3 edgepulse.py \
  --file ip.txt \
  --workers 8 \
  --attempts 2 \
  --timeout 5 \
  --output result.csv
```

برای VPS کوچک ۸ یا ۱۶ Worker مناسب‌تر است. بازه‌ها را محدود کنید، از `--all` غیرضروری استفاده نکنید و اسکن‌های زمان‌بندی‌شده را با فاصله منطقی اجرا کنید.

## Windows، Linux و macOS

در PowerShell:

```powershell
py edgepulse.py --ip 1.1.1.1,1.0.0.1
```

در Linux و macOS:

```bash
python3 edgepulse.py --help
python3 edgepulse.py --file ip.txt
```

خروجی ترمینال در همه سیستم‌ها فقط انگلیسی است.

## Build و Release

ساخت بسته‌ها فقط با GitHub Actions انجام می‌شود. Workflow تست‌ها را اجرا می‌کند، نبود متن CJK را در فایل‌های Runtime بررسی می‌کند، Archive سبک می‌سازد و با Push تگ `v*.*.*` Release را منتشر می‌کند.

تست محلی:

```bash
python3 -m unittest discover -s tests -v
python3 -m py_compile edgepulse.py
python3 edgepulse.py --help
```

## برند و لوگو

لوگوی جدید پروژه در فایل [EDGEPLUSE_LOGO.svg](EDGEPLUSE_LOGO.svg) قرار دارد. موج آبی و فیروزه‌ای نماد عبور سیگنال شبکه از چند Endpoint است.

## Donate

- **Ethereum:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **Bitcoin:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **Tron:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **Solana:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

## License

EdgePulse تحت مجوز MIT منتشر می‌شود. متن کامل در فایل [LICENSE](LICENSE) قرار دارد.
