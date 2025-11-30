#ifndef CONFIG_H
#define CONFIG_H

// WiFi Configuration
#define WIFI_SSID "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define WIFI_TIMEOUT_MS 20000

// API Configuration
#define API_BASE_URL "http://192.168.1.100:8080"
#define API_USERNAME "admin"
#define API_PASSWORD "password"
#define DEVICE_ID "ESP32_MAZE_001"

// GPIO Pin Configuration
#define HALL_SENSOR_START_PIN 7
#define HALL_SENSOR_FINISH_PIN 8
#define LED_PIN 2
#define BUZZER_PIN 4
#define BUTTON_PIN 5

// OLED Configuration (I2C)
#define OLED_SDA_PIN 21
#define OLED_SCL_PIN 22
#define OLED_ADDRESS 0x3C
#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64

// BLE Configuration
#define BLE_DEVICE_NAME "MazeChallenge"
#define BLE_SERVICE_UUID "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
#define BLE_CHAR_STATUS_UUID "beb5483e-36e1-4688-b7f5-ea07361b26a8"
#define BLE_CHAR_TIMER_UUID "beb5483e-36e1-4688-b7f5-ea07361b26a9"
#define BLE_CHAR_CONTROL_UUID "beb5483e-36e1-4688-b7f5-ea07361b26aa"

// Game Configuration
#define DEBOUNCE_DELAY_MS 50
#define ALARM_TIMEOUT_MS 300000  // 5 minutes
#define BATTERY_CHECK_INTERVAL_MS 60000  // 1 minute
#define API_UPDATE_INTERVAL_MS 5000  // 5 seconds

// Battery Monitoring (ADC)
#define BATTERY_PIN 36  // VP pin (ADC1_CH0)
#define BATTERY_MIN_VOLTAGE 3.0
#define BATTERY_MAX_VOLTAGE 4.2

#endif // CONFIG_H
