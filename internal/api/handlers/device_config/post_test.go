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
	"testing"
)

// Mock service for testing
type mockDeviceConfigService struct {
	createFunc func(*models.DeviceConfig, context.Context) error
}

func (m *mockDeviceConfigService) Create(config *models.DeviceConfig, ctx context.Context) error {
	return m.createFunc(config, ctx)
}

func (m *mockDeviceConfigService) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigService) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigService) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigService) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigService) ValidateConfig(config *models.DeviceConfig) error {
	return nil
}

func TestPostHandlerInvalidJSON(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigService{}

	req := httptest.NewRequest(http.MethodPost, "/device/config", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostHandlerValidationError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigService{
		createFunc: func(config *models.DeviceConfig, ctx context.Context) error {
			return device_config.DeviceConfigError{Message: "device_id is required"}
		},
	}

	config := models.DeviceConfig{
		DeviceID:         "",
		AlarmTimeout:     300,
		SensitivityLevel: 5,
		UpdatedAt:        "2024-01-15T10:30:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPost, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostHandlerUniqueConstraintError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigService{
		createFunc: func(config *models.DeviceConfig, ctx context.Context) error {
			return errors.New("UNIQUE constraint failed: device_config.device_id")
		},
	}

	config := models.DeviceConfig{
		DeviceID:         "ARD001",
		AlarmTimeout:     300,
		SensitivityLevel: 5,
		UpdatedAt:        "2024-01-15T10:30:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPost, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestPostHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigService{
		createFunc: func(config *models.DeviceConfig, ctx context.Context) error {
			config.ID = 1
			return nil
		},
	}

	config := models.DeviceConfig{
		DeviceID:         "ARD001",
		AlarmTimeout:     300,
		SensitivityLevel: 5,
		UpdatedAt:        "2024-01-15T10:30:00Z",
	}

	jsonData, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPost, "/device/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.DeviceConfig
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}
}
