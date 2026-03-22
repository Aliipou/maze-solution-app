# Maze Solution App

**IoT maze game platform** — ESP32 hardware with Hall-effect sensors and OLED display, Go REST API backend, React dashboard, and React Native mobile companion.

Players navigate a physical maze while the platform tracks movement via Hall-effect sensors, displays real-time feedback on an OLED screen, and streams data over BLE to the backend. The dashboard visualizes game sessions, device status, and player statistics.

## Architecture

```
┌──────────────┐        BLE        ┌──────────────┐       HTTP       ┌──────────────┐
│  ESP32 Maze  │ ───────────────▶  │   Go REST    │ ◀─────────────▶ │    React     │
│   Hardware   │                   │     API      │                  │  Dashboard   │
│              │                   │  (SQLite)    │                  │              │
│ Hall sensors │                   │  :8080       │                  │  :3000       │
│ OLED SSD1306 │                   └──────────────┘                  └──────────────┘
│ NimBLE       │                          ▲
└──────────────┘                          │ HTTP
                                   ┌──────────────┐
                                   │ React Native │
                                   │   Mobile     │
                                   └──────────────┘
```

## Tech Stack

| Layer     | Technology                                              |
|-----------|---------------------------------------------------------|
| Firmware  | ESP32 (Arduino/PlatformIO), Hall sensors, SSD1306 OLED, NimBLE |
| Backend   | Go 1.22, net/http, SQLite, Basic Auth                   |
| Dashboard | React, TypeScript                                       |
| Mobile    | React Native                                            |
| Infra     | Docker, GitHub Actions CI/CD                            |

## Project Structure

```
maze-solution-app/
├── cmd/api/                  # API entrypoint
├── internal/
│   └── api/
│       ├── handlers/         # HTTP handlers (data, device_config, maze_device)
│       ├── middleware/        # Auth middleware
│       ├── repository/       # SQLite data access
│       ├── server/           # Server setup
│       └── service/          # Business logic
├── firmware/
│   ├── src/main.cpp          # ESP32 firmware
│   └── platformio.ini        # PlatformIO config
├── web/                      # React dashboard (TypeScript)
│   └── src/
│       ├── components/       # Layout, shared components
│       ├── pages/            # Dashboard, Devices, Settings, Statistics
│       └── services/         # API client
├── mobile/                   # React Native companion app
│   └── src/
│       ├── screens/
│       ├── services/
│       └── store/
├── docker-compose.yml        # Single-command deployment
├── Dockerfile
├── Makefile
└── .github/workflows/        # CI and release pipelines
```

## Getting Started

### Prerequisites

- **Go 1.22+** for the API
- **Node.js 18+** for the dashboard
- **PlatformIO** for firmware development
- **Docker** (optional) for containerized deployment

### Run the API

```bash
go build -o api ./cmd/api/main.go
./api
```

The API starts on port **8080** by default.

### Run with Docker

```bash
docker compose up --build
```

### Run the Dashboard

```bash
cd web
npm install
npm start
```

### Flash the Firmware

```bash
cd firmware
pio run --target upload
pio device monitor
```

## API Endpoints

The API exposes RESTful endpoints for three resource types:

| Resource         | Endpoints                          |
|------------------|------------------------------------|
| `/data`          | GET, GET/:id, POST, PUT/:id, DELETE/:id |
| `/device_config` | GET, GET/:id, POST, PUT/:id, DELETE/:id |
| `/maze_device`   | GET, GET/:id, POST, PUT/:id, DELETE/:id |

All endpoints require **Basic Authentication**.

> **Note:** Default credentials (`admin`/`password`) are for development only. Configure real credentials via environment variables `BASIC_AUTH_USERNAME` and `BASIC_AUTH_PASSWORD` in production.

## License

MIT
