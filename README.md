# Maze Challenge IoT System

Complete IoT system for maze-based challenge tracking with ESP32, mobile app, and web dashboard.

![Dashboard](./screenshot_1.jpg)

## Quick Start

### Backend API
```bash
go run ./cmd/api/main.go
# Server: http://localhost:8080
```

### Web Dashboard
```bash
cd web_new
npm run dev
# Dashboard: http://localhost:5173
```

### Mobile App
```bash
cd mobile
npm install
npx react-native run-android
```

### ESP32 Firmware
```bash
cd firmware
pio run --target upload
```

## API Routes

### Device Status
- `GET /device/status` - List all device statuses
- `GET /device/status/{id}` - Get specific status
- `GET /device/status?device_id=ESP32_001` - Filter by device
- `POST /device/status` - Create new status
- `PUT /device/status` - Update status
- `DELETE /device/status/{id}` - Delete status

### Device Configuration
- `GET /device/config` - List all configs
- `GET /device/config/{id}` - Get specific config
- `GET /device/config?device_id=ESP32_001` - Filter by device
- `POST /device/config` - Create config
- `PUT /device/config` - Update config
- `DELETE /device/config/{id}` - Delete config

### General Data
- `GET /data` - List data
- `GET /data/{id}` - Get specific data
- `POST /data` - Create data
- `PUT /data` - Update data
- `DELETE /data/{id}` - Delete data

## Authentication

All API endpoints require Basic Authentication:
- Username: `admin`
- Password: `password`

Example:
```bash
curl http://localhost:8080/device/status -u admin:password
```

## Project Structure

```
.
├── cmd/api/              # Backend API entry point
├── internal/             # Backend core logic
│   ├── api/handlers/    # HTTP handlers
│   ├── api/middleware/  # Auth & CORS
│   └── api/service/     # Business logic
├── firmware/             # ESP32 C++ code
├── mobile/               # React Native app
├── web_new/              # React dashboard (Vite)
├── demo.html             # Standalone dashboard
└── production.db         # SQLite database
```

## Technology Stack

- **Backend**: Go 1.25, SQLite
- **Frontend**: React 18, TypeScript, Material-UI, Vite
- **Mobile**: React Native 0.73
- **Firmware**: C++, ESP32, PlatformIO
- **Testing**: Go testing framework (100+ tests)

## Features

### Backend
- RESTful API with 12 endpoints
- SQLite database
- Basic authentication
- CORS support
- Input validation
- 80-87% test coverage

### Web Dashboard
- Real-time device monitoring
- Live statistics
- Interactive controls
- Auto-refresh (5s)
- Material-UI design

### Mobile App
- BLE connectivity
- Real-time timer
- Statistics & charts
- Push notifications

### ESP32 Firmware
- WiFi & BLE
- Hall effect sensors
- OLED display
- Alarm system
- Battery monitoring

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

## Development

### Backend
```bash
go run ./cmd/api/main.go
```

### Web
```bash
cd web_new
npm install
npm run dev
```

### Mobile
```bash
cd mobile
npm install
npm start
```

### Firmware
```bash
cd firmware
pio run
```

## License

MIT
