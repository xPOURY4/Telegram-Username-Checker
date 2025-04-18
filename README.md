# Telegram Username Checker 

<div align="center">

![GitHub stars](https://img.shields.io/github/stars/xPOURY4/telegram-username-checker?style=for-the-badge&color=yellow)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![GitHub last commit](https://img.shields.io/github/last-commit/xPOURY4/telegram-username-checker?style=for-the-badge)
![Made with Love](https://img.shields.io/badge/Made%20with-❤️-red?style=for-the-badge)

</div>


This advanced tool allows you to check the availability of Telegram usernames using the Telegram API. It handles rate limiting, provides multiple modes of operation, and offers extensive customization options.

## ✨ Features

<div align="center">
  
![Multi-threaded](https://img.shields.io/badge/Multi--threaded-blue?style=flat-square&logo=gitbook&logoColor=white)
![Interactive Mode](https://img.shields.io/badge/Interactive%20Mode-blue?style=flat-square&logo=iterm2&logoColor=white)
![Adaptive Rate Limiting](https://img.shields.io/badge/Adaptive%20Rate%20Limiting-blue?style=flat-square&logo=speedtest&logoColor=white)
![Session Persistence](https://img.shields.io/badge/Session%20Persistence-blue?style=flat-square&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xNyA3aC0xVjZjMC0yLjgtMi4yLTUtNS01UzYgMy4yIDYgNnYxSDVjLTEuMSAwLTIgLjktMiAydjEwYzAgMS4xLjkgMiAyIDJoMTJjMS4xIDAgMi0uOSAyLTJWOWMwLTEuMS0uOS0yLTItMnptLTUgOGMtMS4xIDAtMi0uOS0yLTJzLjktMiAyLTIgMiAuOSAyIDItLjkgMi0yIDJ6bTMuMS04SDguOVY2YzAtMS43IDEuNC0zLjEgMy4xLTMuMSAxLjcgMCAzLjEgMS40IDMuMSAzLjF2MXoiLz48L3N2Zz4=&logoColor=white)
![Username Generation](https://img.shields.io/badge/Username%20Generation-blue?style=flat-square&logo=database&logoColor=white)

</div>

- ✅ Check multiple usernames in bulk or interactive mode
- 🚀 Multi-threaded checking with configurable worker count
- 🛡️ Adaptive rate limiting with automatic backoff
- 💾 Session persistence - resume checking after interruptions
- 🔄 Generate username combinations automatically
- 📊 Detailed progress tracking with visual progress bar
- 📝 Comprehensive results in both text and JSON formats
- 💻 Interactive mode for real-time username checking
- ⚙️ Customizable configuration via config file or command-line arguments
- 🏷️ Categorize usernames as available, taken, invalid, or purchasable
- 🔒 Graceful handling of API rate limits and errors
- 🌐 Proxy support for bypassing restrictions

## 📋 Prerequisites

- Go 1.16 or higher
- Telegram API credentials (API ID and API Hash)

## 🔧 Installation

1. Clone this repository:
```bash
git clone https://github.com/xPOURY4/telegram-username-checker.git
cd telegram-username-checker
```

2. Install dependencies:
```bash
go mod tidy
```

## ⚙️ Configuration

### Option 1: Config File (Recommended)

Create a file named `config.json` with the following structure:

```json
{
  "api_id": "YOUR_API_ID",
  "api_hash": "YOUR_API_HASH",
  "phone": "YOUR_PHONE_NUMBER",
  "base_delay_seconds": 3,
  "max_retries": 3,
  "workers": 1,
  "use_proxy": false,
  "proxy_url": ""
}
```

You can obtain your API ID and API Hash from the [Telegram API Development Tools](https://my.telegram.org/apps).

### Option 2: Environment Variables

Coming in future updates.

### Username Input

Create a file named `usernames.txt` in the same directory. Add the usernames you want to check, one per line.

## 🚀 Usage

### Basic Usage

Run the program with default settings:
```bash
go run main.go
```

### Command-line Arguments

The application supports various command-line arguments for customization:

```
--config        Path to config file (default: "config.json")
--input         Path to input file with usernames (default: "usernames.txt")
--output        Directory for output files (default: "results")
--state         Path to state file (default: "state.json")
--verbose       Enable verbose logging
--generate      Generate username combinations
--min-length    Minimum username length for generation (default: 3)
--max-length    Maximum username length for generation (default: 30)
--interactive   Run in interactive mode
```

### Examples

Check usernames with verbose logging:
```bash
go run main.go --verbose
```

Generate and check username combinations:
```bash
go run main.go --generate --min-length 3 --max-length 5
```

Run in interactive mode:
```bash
go run main.go --interactive
```

Use custom input and output paths:
```bash
go run main.go --input custom_usernames.txt --output custom_results
```

## 🔐 Authentication

When you run the program, you will be prompted to enter the authentication code sent to your Telegram account. After successful authentication, the program will start checking usernames.

## 📊 Results

Results will be saved to the specified output directory (default: "results"):

- `available_usernames.txt`: Usernames available for registration
- `taken_usernames.txt`: Usernames already in use
- `invalid_usernames.txt`: Usernames invalid according to Telegram's rules
- `purchasable_usernames.txt`: Usernames available for purchase
- `error_usernames.txt`: Usernames that couldn't be checked due to errors
- `detailed_results.json`: Complete results with additional metadata

## 🔄 Advanced Features

### State Persistence

The application maintains a state file that allows you to resume checking from where you left off if the program is interrupted. This feature is particularly useful when checking large lists of usernames.

### Adaptive Rate Limiting

The application implements an intelligent adaptive delay system to handle Telegram's rate limiting. If you encounter `FLOOD_WAIT` errors, the program will automatically wait and adjust its request rate.

### Username Generation

The `--generate` flag enables automatic generation of username variations based on the usernames provided in your input file. This can help you discover available variations of desired usernames.

## ❓ Troubleshooting

- If you encounter persistent `FLOOD_WAIT` errors, try increasing the `base_delay_seconds` in your config file
- Ensure your API credentials are correctly configured
- Check that your Telegram account is not limited or banned
- If the application crashes, it should automatically resume from the last saved state when restarted

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ⚠️ Disclaimer

This tool is not officially associated with Telegram. Use it at your own risk and ensure compliance with Telegram's Terms of Service.

## 📞 Contact

<div align="center">
  
[![Telegram](https://img.shields.io/badge/Telegram-xPOURY4-blue?style=for-the-badge&logo=telegram)](https://t.me/xPOURY4)
[![GitHub](https://img.shields.io/badge/GitHub-xPOURY4-black?style=for-the-badge&logo=github)](https://github.com/xPOURY4)

</div>

---



این ابزار پیشرفته به شما امکان می‌دهد با استفاده از API تلگرام، در دسترس بودن نام‌های کاربری تلگرام را بررسی کنید. این برنامه محدودیت‌های نرخ درخواست را مدیریت می‌کند، حالت‌های مختلف عملیاتی را ارائه می‌دهد و گزینه‌های سفارشی‌سازی گسترده‌ای را فراهم می‌کند.

</div>

- ✅ بررسی چندین نام کاربری به صورت دسته‌ای یا حالت تعاملی
- 🚀 بررسی چندنخی با تعداد قابل تنظیم کارگر
- 🛡️ مدیریت انطباقی محدودیت نرخ با عقب‌نشینی خودکار
- 💾 پایداری نشست - ادامه بررسی پس از وقفه‌ها
- 🔄 تولید خودکار ترکیبات نام کاربری
- 📊 ردیابی پیشرفت دقیق با نوار پیشرفت بصری
- 📝 نتایج جامع در هر دو فرمت متن و JSON
- 💻 حالت تعاملی برای بررسی نام کاربری در زمان واقعی
- ⚙️ پیکربندی قابل تنظیم از طریق فایل پیکربندی یا آرگومان‌های خط فرمان
- 🏷️ دسته‌بندی نام‌های کاربری به عنوان در دسترس، گرفته شده، نامعتبر یا قابل خرید
- 🔒 مدیریت مناسب محدودیت‌های نرخ API و خطاها
- 🌐 پشتیبانی از پراکسی برای دور زدن محدودیت‌ها

## 📋 پیش‌نیازها

- Go 1.16 یا بالاتر
- اعتبارنامه‌های API تلگرام (API ID و API Hash)

## 🔧 نصب

۱. کلون کردن این مخزن:
```bash
git clone https://github.com/xPOURY4/telegram-username-checker.git
cd telegram-username-checker
```

۲. نصب وابستگی‌ها:
```bash
go mod tidy
```

## ⚙️ پیکربندی

### گزینه ۱: فایل پیکربندی (توصیه شده)

یک فایل به نام `config.json` با ساختار زیر ایجاد کنید:

```json
{
  "api_id": "API_ID_شما",
  "api_hash": "API_HASH_شما",
  "phone": "شماره_تلفن_شما",
  "base_delay_seconds": 3,
  "max_retries": 3,
  "workers": 1,
  "use_proxy": false,
  "proxy_url": ""
}
```

شما می‌توانید API ID و API Hash خود را از [ابزارهای توسعه API تلگرام](https://my.telegram.org/apps) دریافت کنید.

### گزینه ۲: متغیرهای محیطی

در بروزرسانی‌های آینده اضافه خواهد شد.

### ورودی نام کاربری

یک فایل به نام `usernames.txt` در همان دایرکتوری ایجاد کنید. نام‌های کاربری که می‌خواهید بررسی کنید را، یکی در هر خط، اضافه کنید.

## 🚀 استفاده

### استفاده پایه

اجرای برنامه با تنظیمات پیش‌فرض:
```bash
go run main.go
```

### آرگومان‌های خط فرمان

برنامه از آرگومان‌های مختلف خط فرمان برای سفارشی‌سازی پشتیبانی می‌کند:

```
--config        مسیر فایل پیکربندی (پیش‌فرض: "config.json")
--input         مسیر فایل ورودی با نام‌های کاربری (پیش‌فرض: "usernames.txt")
--output        دایرکتوری برای فایل‌های خروجی (پیش‌فرض: "results")
--state         مسیر فایل وضعیت (پیش‌فرض: "state.json")
--verbose       فعال‌سازی لاگ تفصیلی
--generate      تولید ترکیبات نام کاربری
--min-length    حداقل طول نام کاربری برای تولید (پیش‌فرض: 3)
--max-length    حداکثر طول نام کاربری برای تولید (پیش‌فرض: 30)
--interactive   اجرا در حالت تعاملی
```

### مثال‌ها

بررسی نام‌های کاربری با لاگ تفصیلی:
```bash
go run main.go --verbose
```

تولید و بررسی ترکیبات نام کاربری:
```bash
go run main.go --generate --min-length 3 --max-length 5
```

اجرا در حالت تعاملی:
```bash
go run main.go --interactive
```

استفاده از مسیرهای ورودی و خروجی سفارشی:
```bash
go run main.go --input custom_usernames.txt --output custom_results
```

## 🔐 احراز هویت

هنگامی که برنامه را اجرا می‌کنید، از شما خواسته می‌شود کد احراز هویت ارسال شده به حساب تلگرام خود را وارد کنید. پس از احراز هویت موفق، برنامه شروع به بررسی نام‌های کاربری خواهد کرد.

## 📊 نتایج

نتایج در دایرکتوری خروجی مشخص شده (پیش‌فرض: "results") ذخیره خواهند شد:

- `available_usernames.txt`: نام‌های کاربری در دسترس برای ثبت‌نام
- `taken_usernames.txt`: نام‌های کاربری که در حال حاضر استفاده می‌شوند
- `invalid_usernames.txt`: نام‌های کاربری نامعتبر طبق قوانین تلگرام
- `purchasable_usernames.txt`: نام‌های کاربری در دسترس برای خرید
- `error_usernames.txt`: نام‌های کاربری که به دلیل خطا قابل بررسی نبودند
- `detailed_results.json`: نتایج کامل با متادیتای اضافی

## 🔄 ویژگی‌های پیشرفته

### پایداری وضعیت

برنامه یک فایل وضعیت را نگهداری می‌کند که به شما امکان می‌دهد در صورت قطع برنامه، بررسی را از جایی که متوقف شده‌اید ادامه دهید. این ویژگی به ویژه هنگام بررسی لیست‌های بزرگ نام‌های کاربری مفید است.

### محدودیت نرخ انطباقی

برنامه یک سیستم تأخیر انطباقی هوشمند را برای مدیریت محدودیت نرخ تلگرام پیاده‌سازی می‌کند. اگر با خطاهای `FLOOD_WAIT` مواجه شوید، برنامه به طور خودکار منتظر می‌ماند و نرخ درخواست خود را تنظیم می‌کند.

### تولید نام کاربری

پرچم `--generate` امکان تولید خودکار نسخه‌های مختلف نام کاربری بر اساس نام‌های کاربری ارائه شده در فایل ورودی شما را فراهم می‌کند. این می‌تواند به شما کمک کند تا نسخه‌های در دسترس از نام‌های کاربری مورد نظر را کشف کنید.

## ❓ عیب‌یابی

- اگر به طور مداوم با خطاهای `FLOOD_WAIT` مواجه می‌شوید، سعی کنید `base_delay_seconds` را در فایل پیکربندی خود افزایش دهید
- اطمینان حاصل کنید که اعتبارنامه‌های API شما به درستی پیکربندی شده‌اند
- بررسی کنید که حساب تلگرام شما محدود یا مسدود نشده باشد
- اگر برنامه خراب شود، هنگام راه‌اندازی مجدد باید به طور خودکار از آخرین وضعیت ذخیره شده ادامه دهد

## 🤝 مشارکت

مشارکت‌ها مورد استقبال قرار می‌گیرند! لطفاً برای ارسال Pull Request تردید نکنید.

## ⚠️ سلب مسئولیت

این ابزار به طور رسمی با تلگرام مرتبط نیست. از آن با مسئولیت خود استفاده کنید و اطمینان حاصل کنید که با شرایط خدمات تلگرام مطابقت دارد.


</div>

</div>
