# Uptime Monitoring System

A comprehensive uptime monitoring system for websites built with Go and Fiber framework. This system provides real-time monitoring, historical data tracking, and performance analytics for web services.

## Features

- ✅ **Automated Website Monitoring** - Continuous uptime checking with configurable intervals
- ✅ **Suspended Site Detection** - Intelligent detection of suspended websites
- ✅ **RESTful API** - Complete API for integration and data access
- ✅ **Historical Data & Logs** - Comprehensive logging and historical tracking
- ✅ **High Performance** - Optimized database queries with advanced indexing
- ✅ **Configurable Settings** - Flexible configuration via environment variables
- ✅ **Graceful Shutdown** - Clean shutdown handling
- ✅ **Health Check Endpoints** - Built-in health monitoring
- ✅ **Parallel Processing** - Concurrent monitoring with worker pools
- ✅ **Error Handling** - Robust error handling and recovery

## Prerequisites

- **Go 1.23+**
- **MySQL 8.0+**
- **Git**

## Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd uptime
```

### 2. Copy Environment Configuration
```bash
cp .env.example .env
```

### 3. Configure Environment Variables
Edit the `.env` file with your settings:
```env
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/ms-uptime?charset=utf8mb4&parseTime=True&loc=Local
UPTIME_API_KEY=your_api_key_here
PORT=3000
CHECK_INTERVAL=1m
REQUEST_TIMEOUT=60s
MAX_WORKERS=50
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Database Optimization (Required)
**Before running the application**, optimize your database structure and create performance indexes:
```bash
go run optimize.go
```

This script will:
- Fix table structure (AUTO_INCREMENT, PRIMARY KEY)
- Create 8 performance indexes for optimal query speed
- Analyze tables for better performance
- Optimize database for handling large datasets

### 6. Run the Application
```bash
go run cmd/main.go
```

The server will start on the configured port (default: 3000) and begin monitoring websites every minute.

## Configuration

All settings are configurable via environment variables or the `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `MYSQL_DSN` | `root:@tcp(127.0.0.1:3306)/ms-uptime?charset=utf8mb4&parseTime=True&loc=Local` | MySQL database connection string |
| `PORT` | `3000` | Server port |
| `UPTIME_API_KEY` | - | API authentication key |
| `CHECK_INTERVAL` | `1m` | Monitoring interval (e.g., 30s, 1m, 5m) |
| `REQUEST_TIMEOUT` | `60s` | HTTP request timeout |
| `MAX_WORKERS` | `50` | Maximum concurrent workers |

## API Documentation

### Health Check
```http
GET /health
```
Returns server health status.

### Nodes Management
```http
GET    /api/nodes              # List all monitored nodes
POST   /api/nodes              # Create new node
GET    /api/nodes/:id          # Get specific node
PUT    /api/nodes/:id          # Update node
DELETE /api/nodes/:id          # Delete node
```

### Reports & Analytics
```http
GET    /api/report/get                 # Get URL report
GET    /api/report/get-smart-query     # Smart query report
POST   /api/report/bulk-url/get        # Bulk URL reports
GET    /api/report/all-from-history    # Complete history
GET    /api/report/last                # Latest reports
```

### Node Logs
```http
GET    /api/node-logs/:id             # Get node logs
GET    /api/uptime/:id                # Get uptime statistics
```

## Performance Features

### Database Optimization
- **Advanced Indexing**: 8 specialized indexes for optimal query performance
- **Composite Indexes**: Multi-column indexes for complex queries
- **Query Optimization**: Optimized for handling millions of records
- **Parallel Processing**: Concurrent database operations

### Monitoring Performance
- **Worker Pools**: Configurable concurrent monitoring
- **Intelligent Timeouts**: Adaptive timeout handling
- **Error Recovery**: Automatic retry mechanisms
- **Memory Efficient**: Optimized memory usage for large datasets

## Architecture

```
├── cmd/                    # Application entry points
│   └── main.go            # Main server
├── config/                # Configuration management
├── controllers/           # HTTP request handlers
├── database/              # Database connection & setup
├── models/                # Data models (Node, NodeLog, History)
├── repositories/          # Data access layer
├── routes/                # API route definitions
├── services/              # Business logic layer
├── uptime/                # Uptime monitoring system
├── utils/                 # Utility functions
├── optimize.go            # Database optimization script
└── README.md             # Documentation
```

## Database Schema

### Tables
- **nodes**: Website/service definitions
- **node_logs**: Monitoring results and logs
- **histories**: Historical data and statistics

### Indexes (Created by optimize.go)
- `idx_node_logs_node_id_created_at`: Node + timestamp queries
- `idx_node_logs_node_id_id`: Node + ID sorting
- `idx_node_logs_status`: Status filtering
- `idx_node_logs_node_id_status`: Node + status queries
- `idx_node_logs_up`: Uptime calculations
- `idx_node_logs_composite`: Complex multi-field queries
- `idx_node_logs_created_at`: Time-based queries
- `idx_node_logs_uptime`: Uptime statistics

## Development

### Running in Development Mode
```bash
# Run server
go run cmd/main.go

# Run with custom port
PORT=8080 go run cmd/main.go

# Run with debug logging
DEBUG=true go run cmd/main.go
```

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Building for Production
```bash
# Build binary
go build -o uptime cmd/main.go

# Run binary
./uptime
```

## Deployment

### Docker Deployment
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o uptime cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/uptime .
CMD ["./uptime"]
```

### Systemd Service
```ini
[Unit]
Description=Uptime Monitoring Service
After=network.target

[Service]
Type=simple
User=uptime
WorkingDirectory=/opt/uptime
ExecStart=/opt/uptime/uptime
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Monitoring & Maintenance

### Database Maintenance
```bash
# Re-optimize database after data import
go run optimize.go

# Check database performance
mysql -e "SHOW INDEX FROM node_logs" ms-uptime
```

### Log Management
- Logs are automatically managed by the application
- Historical data is preserved for analytics
- Configure log retention via database settings

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Find process using port
   netstat -ano | findstr :3000
   # Kill process
   taskkill /PID <PID> /F
   ```

2. **Database Connection Issues**
   - Verify MySQL is running
   - Check connection string in `.env`
   - Ensure database exists

3. **Performance Issues**
   - Run `optimize.go` script
   - Check database indexes
   - Monitor worker pool size

### Performance Tuning
- Adjust `MAX_WORKERS` based on server capacity
- Optimize `CHECK_INTERVAL` for your needs
- Monitor database query performance
- Use connection pooling for high load

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the troubleshooting section

---

**Made with ❤️ using Go and Fiber**
