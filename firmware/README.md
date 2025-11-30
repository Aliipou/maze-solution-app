# ESP32 Maze Challenge Firmware

Production-ready firmware for the Interactive Maze Challenge ESP32 device.

## Features

✅ **Hall Sensor Detection** - Start and finish sensors
✅ **OLED Display** - Real-time game status and timer
✅ **WiFi Connectivity** - API integration
✅ **BLE Server** - Mobile app communication
✅ **Audio Feedback** - Buzzer for game events
✅ **Battery Monitoring** - ADC-based level checking
✅ **State Machine** - Robust game logic
✅ **Auto-reconnect** - WiFi and BLE resilience

## Hardware Requirements

| Component | Model | Quantity |
|-----------|-------|----------|
| ESP32 DevKit | ESP32-WROOM-32 | 1 |
| Hall Sensors | A3144 | 2 |
| OLED Display | SSD1306 0.96" I2C | 1 |
| Passive Buzzer | 5V | 1 |
| LED | 5mm | 1 |
| Push Button | Tactile | 1 |
| Resistors | 220Ω, 10kΩ | Pack |
| Neodymium Magnets | 5mm | 5 |

## Wiring Diagram

```
ESP32 Pin Layout:
┌─────────────────────────────────┐
│  GPIO 7  → Hall Sensor START    │
│  GPIO 8  → Hall Sensor FINISH   │
│  GPIO 2  → LED (built-in)       │
│  GPIO 4  → Buzzer                │
│  GPIO 5  → Button (with pullup) │
│  GPIO 21 → OLED SDA             │
│  GPIO 22 → OLED SCL             │
│  GPIO 36 → Battery ADC          │
│  3V3     → Sensors VCC          │
│  GND     → Common Ground        │
└─────────────────────────────────┘
```

## Installation

### 1. Install PlatformIO

```bash
# VS Code Extension (Recommended)
# Install "PlatformIO IDE" from extensions

# OR Command Line
pip install platformio
```

### 2. Configure Device

Edit `include/config.h`:

```cpp
// WiFi
#define WIFI_SSID "YOUR_WIFI_NAME"
#define WIFI_PASSWORD "YOUR_PASSWORD"

// API
#define API_BASE_URL "http://192.168.1.100:8080"
#define API_USERNAME "admin"
#define API_PASSWORD "password"
#define DEVICE_ID "ESP32_MAZE_001"
```

### 3. Build and Upload

```bash
# Build
pio run

# Upload
pio run --target upload

# Monitor serial output
pio device monitor
```

## Game Flow

```
1. READY → Place ball at START sensor
2. PLAYING → Timer starts, alarm active
3. Navigate maze with ball (magnet)
4. COMPLETED → Ball reaches FINISH sensor
5. Stats sent to API
6. Press button to reset
```

## API Integration

### Device Status POST

```json
{
  "device_id": "ESP32_MAZE_001",
  "alarm_active": true,
  "maze_completed": false,
  "hall_sensor_value": false,
  "battery_level": 85,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Endpoint

```
POST http://YOUR_API:8080/device/status
Authorization: Basic base64(admin:password)
```

## BLE Services

### Service UUID
`4fafc201-1fb5-459e-8fcc-c5c9c331914b`

### Characteristics

| UUID | Type | Description |
|------|------|-------------|
| `beb5483e-36e1-4688-b7f5-ea07361b26a8` | READ/NOTIFY | Game status |
| `beb5483e-36e1-4688-b7f5-ea07361b26a9` | READ/NOTIFY | Timer value |
| `beb5483e-36e1-4688-b7f5-ea07361b26aa` | WRITE | Control commands |

## Testing

```bash
# Serial monitor
pio device monitor

# Expected output:
🎮 Maze Challenge ESP32
========================
✅ OLED initialized
Connecting to WiFi...
✅ WiFi connected!
IP: 192.168.1.50
✅ BLE started
✅ Sensors ready
✅ Setup complete!
```

## Troubleshooting

### WiFi Won't Connect

```cpp
// Check credentials in config.h
// Ensure 2.4GHz network (ESP32 doesn't support 5GHz)
// Try increasing WIFI_TIMEOUT_MS
```

### OLED Not Working

```bash
# Check I2C address (default 0x3C)
# Verify SDA/SCL connections
# Try i2c scanner sketch
```

### Sensors Not Triggering

```bash
# Test magnet strength (neodymium recommended)
# Check sensor polarity
# Verify GPIO pins match config
# Adjust DEBOUNCE_DELAY_MS
```

## Configuration Options

```cpp
// Timeouts
ALARM_TIMEOUT_MS = 300000      // 5 minutes
API_UPDATE_INTERVAL_MS = 5000  // 5 seconds
BATTERY_CHECK_INTERVAL_MS = 60000  // 1 minute

// Debounce
DEBOUNCE_DELAY_MS = 50  // Sensor stabilization
```

## Performance

- **Startup Time**: ~3 seconds
- **WiFi Connect**: ~5 seconds
- **Sensor Response**: <50ms
- **API Call**: ~100-300ms
- **Display Update**: 10 FPS
- **Battery Life**: ~8 hours (typical)

## Development

### Build Commands

```bash
make firmware-build     # Build firmware
make firmware-upload    # Upload to device
make firmware-monitor   # Serial monitor
make firmware-clean     # Clean build
```

### Debug Mode

Enable verbose logging in `platformio.ini`:

```ini
build_flags =
    -D CORE_DEBUG_LEVEL=5
```

## Production Deployment

### OTA Updates

```cpp
// TODO: Add OTA update capability
// using ESP32 OTA library
```

### Factory Reset

Hold button for 5 seconds on startup.

## License

MIT License - See main project LICENSE

## Support

Issues: [GitHub Issues](https://github.com/YOUR_USERNAME/maze-solution-api/issues)
