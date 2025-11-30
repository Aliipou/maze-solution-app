package device_config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/device_config"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Mock service for PUT testing
type mockDeviceConfigPutService struct {
	updateFunc func(*models.DeviceConfig, context.Context) (int64, error)
}

func (m *mockDeviceConfigPutService) Create(config *models.DeviceConfig, ctx context.Context) error {
	return nil
}

func (m *mockDeviceConfigPutService) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigPutService) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigPutService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigPutService) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	if m.updateFunc != nil {
		return m.updateFunc(config, ctx)
	}
	return 0, nil
}

func (m *mockDeviceConfigPutService) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigPutService) ValidateConfig(config *models.DeviceConfig) error {
	return nil
}

func TestPutHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 1, nil // 1 row affected
		},
	}

	config := models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_001",
		AlarmTimeout:     600,
		SensitivityLevel: 9,
		UpdatedAt:        "2024-01-15T11:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.DeviceConfig
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}
}

func TestPutHandlerInvalidJSON(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigPutService{}

	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid request data") {
		t.Errorf("Expected error message about invalid request data, got %s", w.Body.String())
	}
}

func TestPutHandlerMissingID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigPutService{}

	config := models.DeviceConfig{
		// ID is 0 (missing)
		DeviceID:         "ESP32_001",
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	expected := `{"error": "ID is required for update."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestPutHandlerValidationError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, device_config.DeviceConfigError{Message: "device_id is required"}
		},
	}

	config := models.DeviceConfig{
		ID:               1,
		DeviceID:         "", // Invalid
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "device_id is required") {
		t.Errorf("Expected error about device_id, got %s", w.Body.String())
	}
}

func TestPutHandlerNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, nil // 0 rows affected
		},
	}

	config := models.DeviceConfig{
		ID:               999,
		DeviceID:         "ESP32_NOTFOUND",
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Device config not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestPutHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, errors.New("database connection error")
		},
	}

	config := models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_001",
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestPutHandlerEmptyBody(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigPutService{}

	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPutHandlerBoundaryValues(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			// Verify boundary values
			if config.AlarmTimeout != 0 || config.SensitivityLevel != 0 {
				t.Errorf("Expected AlarmTimeout=0, SensitivityLevel=0, got AlarmTimeout=%d, SensitivityLevel=%d",
					config.AlarmTimeout, config.SensitivityLevel)
			}
			return 1, nil
		},
	}

	config := models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_001",
		AlarmTimeout:     0, // Boundary value
		SensitivityLevel: 0, // Boundary value
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPutHandlerMaxValues(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigPutService{
		updateFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 1, nil
		},
	}

	config := models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_001",
		AlarmTimeout:     3600, // Max value
		SensitivityLevel: 10,   // Max value
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPut, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
