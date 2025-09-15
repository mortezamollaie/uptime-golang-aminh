# Uptime Monitoring API

A website uptime monitoring system built with Go and Fiber.

## Features
- Website status monitoring and log recording
- Fiber-based RESTful API
- Swagger documentation
- MySQL data storage
- Report management and analytics'


amin h

## Prerequisites
- Docker & Docker Compose
- MySQL (Docker or local)

## Quick Start with Docker

### 1. Build and Run Services

```bash
docker-compose up --build -d
```

### 2. Run Database Optimization Script (Recommended)

```bash
docker exec -it uptime-app ./optimize
```

### 3. View Logs
```bash
docker-compose logs -f
```

### 4. Stop Services
```bash
docker-compose down
```

## Environment Variables (.env)

Example:
```
UPTIME_API_KEY=123456789
MYSQL_DSN=root:@tcp(mysql:3306)/ms-uptime?charset=utf8mb4&parseTime=True&loc=Local
```

- Set `MYSQL_DSN` according to your database configuration.

## Project Structure
```
cmd/           # Main application entry (main.go)
config/        # Project configuration
controllers/   # API controllers
models/        # Database models
repositories/  # Database repositories
routes/        # API route definitions
services/      # Business logic
uptime/        # Uptime checking logic
optimize.go    # Database optimization script
Dockerfile     # Docker image build file
```

## Important Notes
- MySQL database must be ready and accessible before use.
- To fix errors related to the `id` column in `node_logs`, be sure to run optimize.go.
- Swagger documentation is available at `/swagger/index.html` (if enabled).

## Local Development & Run

```bash
go run cmd/main.go
```
Or for database optimization:
```bash
go run optimize.go
```

---

For questions, please open an issue or contact the developer.
