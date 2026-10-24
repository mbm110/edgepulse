# EdgePulse

EdgePulse is a compact, cross-platform network measurement utility for comparing endpoint latency, packet loss, and download throughput across IP ranges.

The command-line program prints English only. This README contains both English and Persian documentation.

## Highlights

- TCP latency checks with bounded concurrency.
- Optional HTTP HEAD checks.
- Optional download throughput measurement.
- IPv4, IPv6, individual addresses, and CIDR ranges.
- Filters for latency, loss, and download speed.
- UTF-8 CSV reports with English column names.
- Release binaries built only by GitHub Actions.

## Quick start

Download the archive for your operating system from **Releases**, extract it, and run the executable from the extracted directory.

```bash
./edgepulse --ip 1.1.1.1,1.0.0.1 --port 443
```

For a larger input, place one address or CIDR range per line in `ip.txt` and run:

```bash
./edgepulse --file ip.txt --download --url https://example.com/test.bin
```

Use `./edgepulse --help` to see every option.

## Examples

```bash
# Test every address in a small CIDR range
./edgepulse --ip 1.1.1.0/29 --all

# Use HTTP HEAD probes
./edgepulse --ip 1.1.1.1 --http --url https://example.com/

# Keep only responsive endpoints
./edgepulse --file ip.txt --max-latency 200 --max-loss 0.25

# Disable CSV output
./edgepulse --ip 1.1.1.1 --output ""
```

## Responsible use

Only test networks and URLs that you are authorized to measure. HTTP mode sends requests to the URL you select. Keep the worker count reasonable and avoid stressing services you do not operate.

## Support

If EdgePulse helps your work, support maintenance through any of these networks:

- **Ethereum:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **Bitcoin:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **Tron:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **Solana:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

---

# راهنمای فارسی

## معرفی

ابزار **EdgePulse** برای مقایسه تأخیر، درصد از دست رفتن بسته‌ها و سرعت دریافت در نشانی‌های شبکه و بازه‌های IP طراحی شده است. این برنامه از IPv4، IPv6، بازه‌های CIDR، آزمایش TCP و آزمایش HTTP پشتیبانی می‌کند.

خروجی برنامه در خط فرمان همیشه انگلیسی است و فایل گزارش CSV نیز ستون‌های انگلیسی دارد.

## شروع سریع

فایل اجرایی مناسب سیستم‌عامل خود را از بخش **Releases** دریافت کنید، از حالت فشرده خارج کنید و در همان پوشه اجرا کنید:

```bash
./edgepulse --ip 1.1.1.1,1.0.0.1 --port 443
```

برای آزمایش چندین نشانی، هر نشانی یا بازه CIDR را در یک خط از فایل `ip.txt` قرار دهید:

```bash
./edgepulse --file ip.txt --download --url https://example.com/test.bin
```

برای مشاهده همه گزینه‌ها، دستور `./edgepulse --help` را اجرا کنید.

## نمونه‌ها

```bash
# آزمایش همه نشانی‌های یک بازه کوچک
./edgepulse --ip 1.1.1.0/29 --all

# استفاده از آزمایش HTTP
./edgepulse --ip 1.1.1.1 --http --url https://example.com/

# نگه‌داشتن مسیرهای سریع‌تر و پایدارتر
./edgepulse --file ip.txt --max-latency 200 --max-loss 0.25

# غیرفعال‌کردن فایل خروجی CSV
./edgepulse --ip 1.1.1.1 --output ""
```

## استفاده مسئولانه

فقط شبکه‌ها و نشانی‌هایی را آزمایش کنید که اجازه بررسی آن‌ها را دارید. حالت HTTP به نشانی انتخابی درخواست ارسال می‌کند. تعداد هم‌زمانی را منطقی انتخاب کنید تا فشاری به سرویس‌ها وارد نشود.

## حمایت مالی

برای حمایت از نگهداری EdgePulse می‌توانید از شبکه‌های زیر استفاده کنید:

- **اتریوم:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **بیت‌کوین:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **ترون:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **سولانا:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

## مجوز

EdgePulse با مجوز MIT منتشر می‌شود. متن کامل مجوز در فایل [LICENSE](LICENSE) قرار دارد.
