package device_config

import (
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

// Mock service for GET testing
type mockDeviceConfigGetService struct {
	readManyFunc       func(int, int, context.Context) ([]*models.DeviceConfig, error)
	readByDeviceIDFunc func(string, context.Context) (*models.DeviceConfig, error)
}

func (m *mockDeviceConfigGetService) Create(config *models.DeviceConfig, ctx context.Context) error {
	return nil
}

func (m *mockDeviceConfigGetService) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigGetService) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	if m.readByDeviceIDFunc != nil {
		return m.readByDeviceIDFunc(deviceID, ctx)
	}
	return nil, nil
}

func (m *mockDeviceConfigGetService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	if m.readManyFunc != nil {
		return m.readManyFunc(page, rowsPerPage, ctx)
	}
	return nil, nil
}

func (m *mockDeviceConfigGetService) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigGetService) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigGetService) ValidateConfig(config *models.DeviceConfig) error {
	return nil
}

func TestGetHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedConfigs := []*models.DeviceConfig{
		{
			ID:               1,
			DeviceID:         "ESP32_001",
			AlarmTimeout:     300,
			SensitivityLevel: 7,
			UpdatedAt:        "2024-01-15T10:00:00Z",
		},
		{
			ID:               2,
			DeviceID:         "ESP32_002",
			AlarmTimeout:     600,
			SensitivityLevel: 5,
			UpdatedAt:        "2024-01-15T11:00:00Z",
		},
	}

	mockService := &mockDeviceConfigGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
			return expectedConfigs, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.DeviceConfig
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 configs, got %d", len(response))
	}
}

func TestGetHandlerWithPagination(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
			if page != 2 || rowsPerPage != 20 {
				t.Errorf("Expected page=2, rowsPerPage=20, got page=%d, rowsPerPage=%d", page, rowsPerPage)
			}
			return []*models.DeviceConfig{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?page=2&rows_per_page=20", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetHandlerByDeviceID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedConfig := &models.DeviceConfig{
		ID:               1,
		DeviceID:         "ESP32_TEST",
		AlarmTimeout:     300,
		SensitivityLevel: 7,
		UpdatedAt:        "2024-01-15T10:00:00Z",
	}

	mockService := &mockDeviceConfigGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
			if deviceID != "ESP32_TEST" {
				t.Errorf("Expected device_id ESP32_TEST, got %s", deviceID)
			}
			return expectedConfig, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?device_id=ESP32_TEST", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.DeviceConfig
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.DeviceID != "ESP32_TEST" {
		t.Errorf("Expected DeviceID ESP32_TEST, got %s", response.DeviceID)
	}
}

func TestGetHandlerByDeviceIDNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
			return nil, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?device_id=NOTFOUND", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Device config not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestGetHandlerByDeviceIDValidationError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
			return nil, device_config.DeviceConfigError{Message: "invalid device_id"}
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?device_id=INVALID", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid device_id") {
		t.Errorf("Expected error message about invalid device_id, got %s", w.Body.String())
	}
}

func TestGetHandlerByDeviceIDInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
			return nil, errors.New("database connection error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?device_id=ESP32_001", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
			return nil, errors.New("database error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetHandlerEmptyResult(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
			return []*models.DeviceConfig{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.DeviceConfig
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected empty array, got %d items", len(response))
	}
}

func TestGetHandlerInvalidPaginationParameters(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockDeviceConfigGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
			// Invalid params should be converted to 0
			if page != 0 || rowsPerPage != 0 {
				t.Errorf("Expected page=0, rowsPerPage=0 for invalid params, got page=%d, rowsPerPage=%d", page, rowsPerPage)
			}
			return []*models.DeviceConfig{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/config?page=invalid&rows_per_page=also_invalid", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
