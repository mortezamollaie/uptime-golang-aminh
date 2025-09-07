# Uptime Monitoring System

یک سیستم مانیتورینگ uptime برای وب‌سایت‌ها که با Go و Fiber framework نوشته شده است.

## ویژگی‌ها

- ✅ مانیتورینگ خودکار وب‌سایت‌ها
- ✅ تشخیص سایت‌های suspended
- ✅ API RESTful کامل
- ✅ ذخیره تاریخچه و لاگ‌ها
- ✅ پیکربندی قابل تنظیم
- ✅ Graceful shutdown
- ✅ Health check endpoint

## نصب و راه‌اندازی

### پیش‌نیازها

- Go 1.23+
- MySQL 8.0+

### مراحل نصب

1. کلون کردن پروژه:
```bash
git clone <repository-url>
cd uptime
```

2. کپی کردن فایل تنظیمات:
```bash
cp .env.example .env
```

3. ویرایش فایل `.env` و تنظیم مقادیر:
```env
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/uptime_db?charset=utf8mb4&parseTime=True&loc=Local
UPTIME_API_KEY=your_api_key_here
PORT=3000
CHECK_INTERVAL=30s
REQUEST_TIMEOUT=30s
MAX_WORKERS=50
```

4. نصب dependencies:
```bash
go mod tidy
```

5. اجرای برنامه:
```bash
go run cmd/main.go
```

## Configuration

تمام تنظیمات از طریق environment variables یا فایل `.env` قابل تغییر هستند:

| Variable | Default | Description |
|----------|---------|-------------|
| `MYSQL_DSN` | `root:@tcp(127.0.0.1:3306)/uptime_db?charset=utf8mb4&parseTime=True&loc=Local` | رشته اتصال به دیتابیس |
| `PORT` | `3000` | پورت سرور |
| `UPTIME_API_KEY` | - | کلید API |
| `CHECK_INTERVAL` | `30s` | فاصله زمانی چک کردن |
| `REQUEST_TIMEOUT` | `30s` | تایم‌اوت درخواست‌ها |
| `MAX_WORKERS` | `50` | حداکثر worker های همزمان |

## API Endpoints

### Health Check
```
GET /health
```

### Nodes Management
```
GET    /api/nodes              # لیست همه nodes
POST   /api/nodes              # ایجاد node جدید
GET    /api/nodes/:id          # دریافت یک node
PUT    /api/nodes/:id          # به‌روزرسانی node
DELETE /api/nodes/:id          # حذف node
```

### Reports
```
GET    /api/report/get                 # گزارش یک URL
GET    /api/report/get-smart-query     # گزارش هوشمند
POST   /api/report/bulk-url/get        # گزارش چند URL
GET    /api/report/all-from-history    # همه تاریخچه
GET    /api/report/last                # آخرین گزارش‌ها
```

## Commands

### اجرای سرور اصلی
```bash
go run cmd/main.go
```

### همگام‌سازی nodes
```bash
go run cmd/node_sync/main.go
```

### پاک‌سازی لاگ‌های قدیمی
```bash
go run cmd/log_checker/main.go
```

## Architecture

```
├── cmd/                    # Entry points
│   ├── main.go            # سرور اصلی
│   ├── node_sync/         # همگام‌سازی nodes
│   ├── log_checker/       # پاک‌سازی لاگ‌ها
│   └── starter/           # راه‌انداز اولیه
├── config/                # پیکربندی
├── controllers/           # کنترلرها
├── database/              # اتصال دیتابیس
├── models/                # مدل‌های دیتا
├── repositories/          # لایه دیتا
├── routes/                # روت‌ها
├── services/              # منطق تجاری
├── uptime/                # سیستم uptime check
└── logger/                # سیستم لاگ
```

## Development

برای توسعه پروژه:

1. Fork کردن پروژه
2. ایجاد branch جدید
3. Commit کردن تغییرات
4. Push به branch
5. ایجاد Pull Request

## License

MIT License
