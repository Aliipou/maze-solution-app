# Installation Guide

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [Method 1: From Source](#method-1-from-source)
  - [Method 2: Using Docker](#method-2-using-docker)
  - [Method 3: Pre-built Binary](#method-3-pre-built-binary)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Server](#running-the-server)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required

- **Go 1.21 or higher** - [Download](https://golang.org/dl/)
- **SQLite3** - Usually pre-installed on Linux/macOS
- **Git** - For cloning the repository

### Optional

- **Docker & Docker Compose** - For containerized deployment
- **Make** - For using Makefile commands
- **golangci-lint** - For code quality checks

### System Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 1 core | 2+ cores |
| RAM | 512 MB | 1+ GB |
| Disk | 50 MB | 100+ MB |
| OS | Linux, macOS, Windows | Linux/macOS |

---

## Installation Methods

### Method 1: From Source

#### Step 1: Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/maze-solution-api.git
cd maze-solution-api
```

#### Step 2: Install Dependencies

```bash
go mod download
go mod verify
```

#### Step 3: Build the Application

```bash
go build -o api ./cmd/api/main.go
```

**For production builds with optimizations:**

```bash
CGO_ENABLED=1 go build -ldflags="-w -s" -o api ./cmd/api/main.go
```

#### Step 4: Configure Environment

```bash
cp .env.example .env
# Edit .env with your preferred editor
nano .env
```

#### Step 5: Run the Server

```bash
./api
```

Server will start on `http://localhost:8080`

---

### Method 2: Using Docker

#### Step 1: Clone Repository

```bash
git clone https://github.com/YOUR_USERNAME/maze-solution-api.git
cd maze-solution-api
```

#### Step 2: Using Docker Compose (Recommended)

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

#### Step 3: Using Docker Directly

```bash
# Build image
docker build -t maze-api:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  -e BASIC_AUTH_USERNAME=admin \
  -e BASIC_AUTH_PASSWORD=password \
  --name maze-api \
  maze-api:latest

# View logs
docker logs -f maze-api

# Stop container
docker stop maze-api
```

---

### Method 3: Pre-built Binary

#### Step 1: Download Binary

Download the latest release from [GitHub Releases](https://github.com/YOUR_USERNAME/maze-solution-api/releases):

```bash
# Linux AMD64
wget https://github.com/YOUR_USERNAME/maze-solution-api/releases/download/v1.0.0/api-linux-amd64

# macOS ARM64 (M1/M2)
wget https://github.com/YOUR_USERNAME/maze-solution-api/releases/download/v1.0.0/api-darwin-arm64

# Windows AMD64
wget https://github.com/YOUR_USERNAME/maze-solution-api/releases/download/v1.0.0/api-windows-amd64.exe
```

#### Step 2: Make Executable (Linux/macOS)

```bash
chmod +x api-linux-amd64
mv api-linux-amd64 api
```

#### Step 3: Run

```bash
./api
```

---

## Configuration

### Environment Variables

Create a `.env` file or set environment variables:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DB_PATH=production.db

# Logging
LOG_FILE=production.log
LOG_LEVEL=info

# Authentication
BASIC_AUTH_USERNAME=admin
BASIC_AUTH_PASSWORD=your_secure_password_here

# CORS
CORS_ALLOWED_ORIGINS=*
```

### Configuration File

Alternatively, you can modify the configuration in code:

Edit `cmd/api/main.go` or create a `config.yaml` file (if implementing config file support).

---

## Database Setup

### Automatic Setup

The database is created automatically on first run:

```bash
./api
# Creates production.db with all required tables
```

### Manual Setup (Optional)

If you want to pre-create the database:

```bash
sqlite3 production.db < schema.sql
```

### Database Tables

The following tables are created automatically:

1. **maze_device_status** - Game session tracking
2. **device_config** - Device configuration
3. **data** - Generic data storage

### Migrations

Currently, the schema is created on startup. For future migrations:

```bash
# TODO: Add migration tool (e.g., golang-migrate)
```

---

## Running the Server

### Development Mode

```bash
# Using go run
go run ./cmd/api/main.go

# Using Make
make run

# With live reload (install air first)
air
```

### Production Mode

```bash
# Using binary
./api

# Using systemd (Linux)
sudo systemctl start maze-api

# Using Docker
docker-compose up -d
```

### Running in Background

#### Linux/macOS (using nohup)

```bash
nohup ./api > server.log 2>&1 &
```

#### Using screen

```bash
screen -S maze-api
./api
# Press Ctrl+A then D to detach
```

#### Using systemd (Linux)

Create `/etc/systemd/system/maze-api.service`:

```ini
[Unit]
Description=Maze Solution API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/maze-api
ExecStart=/opt/maze-api/api
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable maze-api
sudo systemctl start maze-api
sudo systemctl status maze-api
```

---

## Verification

### Health Check

```bash
curl http://localhost:8080/data
```

Expected response: `200 OK` (may require authentication)

### Test API Endpoints

```bash
# Test device status endpoint
curl -u admin:password \
  -H "Content-Type: application/json" \
  http://localhost:8080/device/status

# Create a test status
curl -X POST http://localhost:8080/device/status \
  -u admin:password \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "TEST001",
    "alarm_active": false,
    "maze_completed": false,
    "hall_sensor_value": false,
    "battery_level": 100,
    "timestamp": "2024-01-15T10:00:00Z"
  }'
```

### Check Logs

```bash
# View log file
tail -f production.log

# Docker logs
docker-compose logs -f
```

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use a different port
SERVER_PORT=8081 ./api
```

### Permission Denied (Database)

```bash
# Fix permissions
chmod 755 .
chmod 644 production.db

# Or run with sudo (not recommended)
sudo ./api
```

### Build Errors

```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build -v ./cmd/api/main.go
```

### CGO Errors (SQLite)

```bash
# Install build tools (Ubuntu/Debian)
sudo apt-get install build-essential

# Install build tools (macOS)
xcode-select --install

# Build with CGO
CGO_ENABLED=1 go build ./cmd/api/main.go
```

### Connection Refused

```bash
# Check if server is running
ps aux | grep api

# Check firewall
sudo ufw status
sudo ufw allow 8080

# Check if listening on correct interface
netstat -tlnp | grep 8080
```

### Authentication Issues

```bash
# Verify credentials
echo -n "admin:password" | base64
# Use the output in Authorization header

# Test without auth requirement (temporarily disable in code)
```

---

## Next Steps

After successful installation:

1. Read the [API Documentation](../docs/openapi.yaml)
2. Review [Architecture](../README.md#architecture)
3. Check [Development Guide](../CONTRIBUTING.md)
4. Set up [ESP32 Firmware](../firmware/README.md) (when available)
5. Deploy [Mobile Apps](../mobile/README.md) (when available)

---

## Getting Help

- 📖 [Documentation](../README.md)
- 🐛 [Report Issues](https://github.com/YOUR_USERNAME/maze-solution-api/issues)
- 💬 [Discussions](https://github.com/YOUR_USERNAME/maze-solution-api/discussions)

---

**Happy Coding! 🚀**
