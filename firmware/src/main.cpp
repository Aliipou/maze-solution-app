#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <NimBLEDevice.h>
#include "../include/config.h"

// ==================== Global Objects ====================
Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, -1);
NimBLEServer* pServer = nullptr;
NimBLECharacteristic* pStatusCharacteristic = nullptr;
NimBLECharacteristic* pTimerCharacteristic = nullptr;

// ==================== Game State ====================
enum GameState {
    STATE_IDLE,
    STATE_READY,
    STATE_PLAYING,
    STATE_COMPLETED,
    STATE_ERROR
};

struct MazeGame {
    GameState state;
    unsigned long startTime;
    unsigned long elapsedTime;
    bool alarmActive;
    bool mazeCompleted;
    bool hallSensorStart;
    bool hallSensorFinish;
    int batteryLevel;
    int sessionId;
    bool bleConnected;
};

MazeGame game = {STATE_IDLE, 0, 0, false, false, false, false, 100, -1, false};

// ==================== Function Prototypes ====================
void setupWiFi();
void setupOLED();
void setupBLE();
void setupSensors();
int readBatteryLevel();
void updateDisplay();
void handleGameLogic();
void sendStatusToAPI();
void playTone(int frequency, int duration);
String getCurrentTimestamp();
void handleSensorStart();
void handleSensorFinish();

// ==================== BLE Callbacks ====================
class ServerCallbacks: public NimBLEServerCallbacks {
    void onConnect(NimBLEServer* pServer) {
        game.bleConnected = true;
        Serial.println("BLE Client connected");
    }

    void onDisconnect(NimBLEServer* pServer) {
        game.bleConnected = false;
        Serial.println("BLE Client disconnected");
        NimBLEDevice::startAdvertising();
    }
};

// ==================== Setup ====================
void setup() {
    Serial.begin(115200);
    Serial.println("\n\n🎮 Maze Challenge ESP32");
    Serial.println("========================");

    // Initialize pins
    pinMode(HALL_SENSOR_START_PIN, INPUT);
    pinMode(HALL_SENSOR_FINISH_PIN, INPUT);
    pinMode(LED_PIN, OUTPUT);
    pinMode(BUZZER_PIN, OUTPUT);
    pinMode(BUTTON_PIN, INPUT_PULLUP);

    // Initialize I2C for OLED
    Wire.begin(OLED_SDA_PIN, OLED_SCL_PIN);

    // Setup components
    setupOLED();
    setupWiFi();
    setupBLE();
    setupSensors();

    Serial.println("✅ Setup complete!");
    game.state = STATE_READY;
    updateDisplay();
    playTone(1000, 100);
}

// ==================== Main Loop ====================
void loop() {
    static unsigned long lastAPIUpdate = 0;
    static unsigned long lastBatteryCheck = 0;
    static unsigned long lastDisplayUpdate = 0;

    unsigned long now = millis();

    // Read sensors
    handleSensorStart();
    handleSensorFinish();

    // Update game logic
    handleGameLogic();

    // Update display (10 FPS)
    if (now - lastDisplayUpdate >= 100) {
        updateDisplay();
        lastDisplayUpdate = now;
    }

    // Send status to API
    if (now - lastAPIUpdate >= API_UPDATE_INTERVAL_MS) {
        sendStatusToAPI();
        lastAPIUpdate = now;
    }

    // Check battery
    if (now - lastBatteryCheck >= BATTERY_CHECK_INTERVAL_MS) {
        game.batteryLevel = readBatteryLevel();
        lastBatteryCheck = now;
    }

    // Update BLE characteristics
    if (game.bleConnected && pTimerCharacteristic) {
        String timerStr = String(game.elapsedTime / 1000);
        pTimerCharacteristic->setValue(timerStr.c_str());
        pTimerCharacteristic->notify();
    }

    delay(10);
}

// ==================== WiFi Setup ====================
void setupWiFi() {
    Serial.print("Connecting to WiFi");
    display.clearDisplay();
    display.setCursor(0, 20);
    display.println("Connecting WiFi...");
    display.display();

    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    unsigned long startAttempt = millis();
    while (WiFi.status() != WL_CONNECTED &&
           millis() - startAttempt < WIFI_TIMEOUT_MS) {
        Serial.print(".");
        delay(500);
    }

    if (WiFi.status() == WL_CONNECTED) {
        Serial.println("\n✅ WiFi connected!");
        Serial.print("IP: ");
        Serial.println(WiFi.localIP());
    } else {
        Serial.println("\n❌ WiFi connection failed!");
        game.state = STATE_ERROR;
    }
}

// ==================== OLED Setup ====================
void setupOLED() {
    if (!display.begin(SSD1306_SWITCHCAPVCC, OLED_ADDRESS)) {
        Serial.println("❌ OLED initialization failed!");
        return;
    }

    display.clearDisplay();
    display.setTextSize(1);
    display.setTextColor(SSD1306_WHITE);
    display.setCursor(0, 0);
    display.println("Maze Challenge");
    display.println("Initializing...");
    display.display();

    Serial.println("✅ OLED initialized");
}

// ==================== BLE Setup ====================
void setupBLE() {
    Serial.println("Setting up BLE...");

    NimBLEDevice::init(BLE_DEVICE_NAME);
    pServer = NimBLEDevice::createServer();
    pServer->setCallbacks(new ServerCallbacks());

    NimBLEService* pService = pServer->createService(BLE_SERVICE_UUID);

    pStatusCharacteristic = pService->createCharacteristic(
        BLE_CHAR_STATUS_UUID,
        NIMBLE_PROPERTY::READ | NIMBLE_PROPERTY::NOTIFY
    );

    pTimerCharacteristic = pService->createCharacteristic(
        BLE_CHAR_TIMER_UUID,
        NIMBLE_PROPERTY::READ | NIMBLE_PROPERTY::NOTIFY
    );

    NimBLECharacteristic* pControlCharacteristic = pService->createCharacteristic(
        BLE_CHAR_CONTROL_UUID,
        NIMBLE_PROPERTY::WRITE
    );

    pService->start();

    NimBLEAdvertising* pAdvertising = NimBLEDevice::getAdvertising();
    pAdvertising->addServiceUUID(BLE_SERVICE_UUID);
    pAdvertising->start();

    Serial.println("✅ BLE started");
}

// ==================== Sensor Setup ====================
void setupSensors() {
    Serial.println("Testing sensors...");

    bool startSensor = digitalRead(HALL_SENSOR_START_PIN);
    bool finishSensor = digitalRead(HALL_SENSOR_FINISH_PIN);

    Serial.printf("Start Sensor: %d, Finish Sensor: %d\n", startSensor, finishSensor);
    Serial.println("✅ Sensors ready");
}

// ==================== Sensor Handlers ====================
void handleSensorStart() {
    static unsigned long lastDebounce = 0;
    static bool lastState = LOW;

    bool reading = digitalRead(HALL_SENSOR_START_PIN);

    if (reading != lastState) {
        lastDebounce = millis();
    }

    if ((millis() - lastDebounce) > DEBOUNCE_DELAY_MS) {
        if (reading == HIGH && game.hallSensorStart == false) {
            game.hallSensorStart = true;

            if (game.state == STATE_READY) {
                // Start game
                game.state = STATE_PLAYING;
                game.startTime = millis();
                game.alarmActive = true;
                game.mazeCompleted = false;

                Serial.println("🎮 Game Started!");
                playTone(1500, 200);
                digitalWrite(LED_PIN, HIGH);
            }
        } else if (reading == LOW) {
            game.hallSensorStart = false;
        }
    }

    lastState = reading;
}

void handleSensorFinish() {
    static unsigned long lastDebounce = 0;
    static bool lastState = LOW;

    bool reading = digitalRead(HALL_SENSOR_FINISH_PIN);

    if (reading != lastState) {
        lastDebounce = millis();
    }

    if ((millis() - lastDebounce) > DEBOUNCE_DELAY_MS) {
        if (reading == HIGH && game.hallSensorFinish == false) {
            game.hallSensorFinish = true;

            if (game.state == STATE_PLAYING) {
                // Complete game
                game.state = STATE_COMPLETED;
                game.elapsedTime = millis() - game.startTime;
                game.mazeCompleted = true;
                game.alarmActive = false;

                Serial.printf("🏆 Game Completed! Time: %lu ms\n", game.elapsedTime);
                playTone(2000, 500);
                digitalWrite(LED_PIN, LOW);

                sendStatusToAPI();
            }
        } else if (reading == LOW) {
            game.hallSensorFinish = false;
        }
    }

    lastState = reading;
}

// ==================== Game Logic ====================
void handleGameLogic() {
    if (game.state == STATE_PLAYING) {
        game.elapsedTime = millis() - game.startTime;

        // Check timeout
        if (game.elapsedTime >= ALARM_TIMEOUT_MS) {
            Serial.println("⏰ Timeout!");
            game.state = STATE_IDLE;
            game.alarmActive = false;
            digitalWrite(LED_PIN, LOW);
        }
    }

    // Button press to reset
    if (digitalRead(BUTTON_PIN) == LOW && game.state == STATE_COMPLETED) {
        delay(50); // Debounce
        if (digitalRead(BUTTON_PIN) == LOW) {
            Serial.println("🔄 Resetting...");
            game.state = STATE_READY;
            game.elapsedTime = 0;
            game.mazeCompleted = false;
            playTone(1000, 100);
        }
    }
}

// ==================== Display Update ====================
void updateDisplay() {
    display.clearDisplay();
    display.setTextSize(1);
    display.setCursor(0, 0);

    // Title
    display.println("MAZE CHALLENGE");
    display.drawLine(0, 10, SCREEN_WIDTH, 10, SSD1306_WHITE);

    // State
    display.setCursor(0, 15);
    switch (game.state) {
        case STATE_IDLE:
            display.println("Status: IDLE");
            break;
        case STATE_READY:
            display.println("Status: READY");
            display.println("\nPlace ball at START");
            break;
        case STATE_PLAYING:
            display.println("Status: PLAYING");
            display.println();
            display.setTextSize(2);
            display.printf("%02lu:%02lu",
                game.elapsedTime / 60000,
                (game.elapsedTime / 1000) % 60);
            display.setTextSize(1);
            break;
        case STATE_COMPLETED:
            display.println("Status: COMPLETED!");
            display.println();
            display.printf("Time: %lu.%03lu s\n",
                game.elapsedTime / 1000,
                game.elapsedTime % 1000);
            display.println("\nPress button to reset");
            break;
        case STATE_ERROR:
            display.println("Status: ERROR");
            break;
    }

    // Bottom info
    display.setCursor(0, 56);
    display.printf("Bat:%d%% ", game.batteryLevel);
    if (WiFi.status() == WL_CONNECTED) display.print("WiFi");
    if (game.bleConnected) display.print(" BLE");

    display.display();
}

// ==================== API Communication ====================
void sendStatusToAPI() {
    if (WiFi.status() != WL_CONNECTED) {
        Serial.println("❌ No WiFi connection");
        return;
    }

    HTTPClient http;
    String url = String(API_BASE_URL) + "/device/status";

    http.begin(url);
    http.setAuthorization(API_USERNAME, API_PASSWORD);
    http.addHeader("Content-Type", "application/json");

    // Create JSON payload
    JsonDocument doc;
    doc["device_id"] = DEVICE_ID;
    doc["alarm_active"] = game.alarmActive;
    doc["maze_completed"] = game.mazeCompleted;
    doc["hall_sensor_value"] = game.hallSensorFinish;
    doc["battery_level"] = game.batteryLevel;
    doc["timestamp"] = getCurrentTimestamp();

    String payload;
    serializeJson(doc, payload);

    Serial.println("📤 Sending to API: " + payload);

    int httpCode = http.POST(payload);

    if (httpCode > 0) {
        Serial.printf("✅ API Response: %d\n", httpCode);
        String response = http.getString();
        Serial.println(response);
    } else {
        Serial.printf("❌ API Error: %s\n", http.errorToString(httpCode).c_str());
    }

    http.end();
}

String getCurrentTimestamp() {
    // Simple ISO 8601 format (would need NTP for accurate time)
    time_t now = time(nullptr);
    struct tm timeinfo;
    gmtime_r(&now, &timeinfo);

    char buffer[30];
    strftime(buffer, sizeof(buffer), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
    return String(buffer);
}

// ==================== Battery Monitoring ====================
int readBatteryLevel() {
    int rawValue = analogRead(BATTERY_PIN);
    float voltage = (rawValue / 4095.0) * 3.3 * 2;  // Voltage divider

    int percentage = ((voltage - BATTERY_MIN_VOLTAGE) /
                      (BATTERY_MAX_VOLTAGE - BATTERY_MIN_VOLTAGE)) * 100;

    return constrain(percentage, 0, 100);
}

// ==================== Audio Feedback ====================
void playTone(int frequency, int duration) {
    tone(BUZZER_PIN, frequency, duration);
    delay(duration);
    noTone(BUZZER_PIN);
}
