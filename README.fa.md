# EdgePulse؛ اسکنر IPهای CDN و ابزار سنجش تأخیر و سرعت

[English](README.md) | [فارسی](README.fa.md)

ابزار **EdgePulse** یک برنامه خط فرمان نوشته‌شده با Go است که برای اسکن IPv4، IPv6 و بازه‌های CIDR طراحی شده است. برنامه تأخیر TCP یا HTTP، Packet Loss، سرعت دانلود، Region مربوط به CDN و نتیجه کامل CSV را بررسی می‌کند.

همه پیام‌های اجرایی برنامه فقط انگلیسی هستند.

## قابلیت‌ها

- سنجش هم‌زمان تأخیر TCP.
- سنجش تأخیر HTTP با Status Code قابل قبول.
- پشتیبانی از IPv4، IPv6، ورودی ترکیبی، CIDR و IP تکی.
- انتخاب یک IP از هر بازه یا اسکن کامل IPv4 با `-allip`.
- سنجش سرعت دانلود با URL و زمان قابل تنظیم.
- فیلتر تأخیر، Packet Loss، سرعت، HTTP Status و Region.
- تشخیص Region از Headerهای رایج CDN.
- خروجی CSV با Headerهای انگلیسی.
- اسکریپت‌های اختیاری برای hosts، DNS، dnsmasq و proxy.
- Release استاتیک Linux ARM64 برای محیط‌های سازگار با Termux.

## شروع سریع

فایل Release را دریافت و Extract کنید و سپس اجرا کنید:

```bash
./edgepulse
```

برای آزمایش IP یا بازه مشخص:

```bash
./edgepulse -ip 1.1.1.1,1.0.0.1,2606:4700::/32
```

برای فیلتر و سنجش سرعت دانلود:

```bash
./edgepulse -tl 200 -tlr 0.25 -sl 1 -dn 10
```

برای دیدن همه گزینه‌ها:

```bash
./edgepulse -h
```

## روند کار

1. خواندن بازه‌ها از `ip.txt` یا گزینه `-ip`.
2. انتخاب یک IP از هر بازه یا همه IPv4ها با `-allip`.
3. سنجش هم‌زمان تأخیر TCP.
4. ارسال درخواست HTTP و بررسی Status Code با `-httping`.
5. فیلتر تأخیر، Packet Loss و Region.
6. سنجش سرعت دانلود مگر اینکه `-dd` فعال باشد.
7. مرتب‌سازی نتایج و نوشتن `result.csv`.

## گزینه‌ها

```text
-n 200          تعداد Workerهای سنجش تأخیر
-t 4            تعداد Probe برای هر IP
-dn 10          تعداد مقصدهای سنجش دانلود
-dt 10          حداکثر زمان دانلود هر IP برحسب ثانیه
-tp 443         پورت آزمایش
-url URL        URL آزمایش HTTP و دانلود
-httping        فعال‌کردن حالت HTTP
-httping-code N کد HTTP قابل قبول
-cfcolo CODES   فیلتر Region با ویرگول
-tl MS          حداکثر تأخیر میانگین
-tll MS         حداقل تأخیر میانگین
-tlr RATIO     حداکثر نسبت Loss از صفر تا یک
-sl MBPS        حداقل سرعت دانلود
-p N            تعداد ردیف‌های نمایش‌داده‌شده
-f FILE         فایل ورودی IP و CIDR
-ip RANGES      IPها و CIDRهای جداشده با ویرگول
-o FILE         فایل CSV؛ مقدار خالی خروجی فایل را خاموش می‌کند
-dd             غیرفعال‌کردن سنجش دانلود
-allip          آزمایش همه IPv4های هر بازه
-debug          نمایش اطلاعات عیب‌یابی
-v              نمایش نسخه و بررسی Update
-h              نمایش راهنما
```

## حالت TCP و HTTP

حالت TCP امکان بازشدن اتصال به پورت مشخص را بررسی می‌کند. حالت HTTP، Hostname و URL انتخابی را از طریق هر IP آزمایش می‌کند و در صورت وجود Header مناسب، Region را می‌خواند.

```bash
./edgepulse -httping -url https://example.com/
./edgepulse -httping -tp 80 -url http://example.com/
```

حالت HTTP نوعی اسکن شبکه است. روی VPS مقدار `-n` را پایین بگذارید و قوانین سرویس مقصد را رعایت کنید. فقط شبکه‌ها و URLهایی را آزمایش کنید که اجازه بررسی آن‌ها را دارید.

## IPv4، IPv6 و بازه‌ها

فایل ورودی می‌تواند هر دو نوع IP را داشته باشد:

```text
1.1.1.1
1.0.0.1/24
2606:4700::/32
```

برای IP تکی نیازی به Mask نیست. به‌صورت پیش‌فرض از هر بازه یک IP انتخاب می‌شود. گزینه `-allip` را فقط برای بازه‌های مجاز و قابل مدیریت استفاده کنید.

## فیلتر Region

در حالت HTTP، اگر CDN اطلاعات Region را در Header قرار دهد، برنامه آن را می‌خواند. این مقدار ممکن است کد فرودگاه، کشور یا شهر باشد:

```bash
./edgepulse -httping -cfcolo FRA,LAX,SJC
```

مسیریابی Region می‌تواند با توجه به شبکه، Provider و زمان تغییر کند.

## خروجی

Headerهای جدول ترمینال انگلیسی هستند:

```text
Address           Sent  Received  Loss %  Average Latency  Download Speed (MB/s)  Region
```

نتیجه کامل در `result.csv` ذخیره می‌شود. برای خاموش‌کردن خروجی فایل از `-o ""` استفاده کنید.

## Termux و Android

فایل `edgepulse-linux-arm64.tar.gz` برای محیط‌های سازگار Linux ARM64 ساخته شده است. اگر باینری Linux در Userspace اندروید Termux سازگار نبود، سورس را با Go خود Termux بسازید:

```bash
pkg update
pkg install golang git
git clone https://github.com/mbm110/edgepulse.git
cd edgepulse
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o edgepulse main.go
./edgepulse -h
```

## Release

GitHub Actions نسخه‌های Linux AMD64، Linux ARM64 استاتیک، Windows AMD64، macOS Intel، macOS Apple Silicon و Source را می‌سازد. نسخه Linux ARM64 با `CGO_ENABLED=0` و `-ldflags="-s -w"` تولید می‌شود.

## Donate

- **Ethereum:** `0xA9e165Db72a6f1B3261dfcFb7cc9ADD725Ea2f25`
- **Bitcoin:** `bc1q4ntxszj643n0ujj3uuncljk8uyprug4swprxlg`
- **Tron:** `TN58EHrD9AvtVEx4FSka5kitwpgXdXbri9`
- **Solana:** `3wPEfTm1bCbrU9aqie7xMgAXJK1qXGktRn2yjtjb6LPt`

## License

EdgePulse تحت مجوز GPLv3 منتشر می‌شود. متن کامل در [LICENSE](LICENSE) قرار دارد.
