package device_config

import (
	"context"
	"encoding/json"
	"errors"
	"goapi/internal/api/repository/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Mock service for GetByID testing
type mockDeviceConfigGetByIDService struct {
	readOneFunc func(int, context.Context) (*models.DeviceConfig, error)
}

func (m *mockDeviceConfigGetByIDService) Create(config *models.DeviceConfig, ctx context.Context) error {
	return nil
}

func (m *mockDeviceConfigGetByIDService) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	if m.readOneFunc != nil {
		return m.readOneFunc(id, ctx)
	}
	return nil, nil
}

func (m *mockDeviceConfigGetByIDService) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigGetByIDService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigGetByIDService) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigGetByIDService) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigGetByIDService) ValidateConfig(config *models.DeviceConfig) error {
	return nil
}

func TestGetByIDHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedConfig := &models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_001",
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			if id != 1 {
				t.Errorf("Expected ID 1, got %d", id)
			}
			return expectedConfig, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

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

	if response.DeviceID != "ESP32_001" {
		t.Errorf("Expected DeviceID ESP32_001, got %s", response.DeviceID)
	}
}

func TestGetByIDHandlerInvalidID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigGetByIDService{}

	req := httptest.NewRequest(http.MethodGet, "/device/config/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	expected := `{"error": "Invalid ID format."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestGetByIDHandlerNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			return nil, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Device config not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestGetByIDHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			return nil, errors.New("database connection error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetByIDHandlerNegativeID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			if id != -1 {
				t.Errorf("Expected ID -1, got %d", id)
			}
			return nil, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/-1", nil)
	req.SetPathValue("id", "-1")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetByIDHandlerZeroID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			if id != 0 {
				t.Errorf("Expected ID 0, got %d", id)
			}
			return nil, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetByIDHandlerLargeID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetByIDService{
		readOneFunc: func(id int, ctx context.Context) (*models.DeviceConfig, error) {
			if id != 999999 {
				t.Errorf("Expected ID 999999, got %d", id)
			}
			return nil, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config/999999", nil)
	req.SetPathValue("id", "999999")
	w := httptest.NewRecorder()

	GetByIDHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
