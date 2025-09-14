# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.io,direct
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o uptime ./cmd/main.go
RUN go build -o optimize ./optimize.go
RUN go build -o log_checker ./cmd/log_checker/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/uptime ./uptime
COPY --from=builder /app/optimize ./optimize
COPY --from=builder /app/log_checker ./log_checker
EXPOSE 3000
ENV PORT=3000
RUN chmod 755 ./uptime ./optimize ./log_checker
CMD ["/bin/sh", "-c", "./uptime"]
