package device_config

import (
	"context"
	"errors"
	"goapi/internal/api/repository/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Mock service for delete testing
type mockDeviceConfigDeleteService struct {
	deleteFunc func(*models.DeviceConfig, context.Context) (int64, error)
}

func (m *mockDeviceConfigDeleteService) Create(config *models.DeviceConfig, ctx context.Context) error {
	return nil
}

func (m *mockDeviceConfigDeleteService) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigDeleteService) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigDeleteService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	return nil, nil
}

func (m *mockDeviceConfigDeleteService) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDeviceConfigDeleteService) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return m.deleteFunc(config, ctx)
}

func (m *mockDeviceConfigDeleteService) ValidateConfig(config *models.DeviceConfig) error {
	return nil
}

func TestDeleteHandlerInvalidID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigDeleteService{}

	req := httptest.NewRequest(http.MethodDelete, "/device/config/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	expected := `{"error": "Invalid ID format."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDeleteHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigDeleteService{
		deleteFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, errors.New("database error")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/config/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestDeleteHandlerNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigDeleteService{
		deleteFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, nil // 0 rows affected
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/config/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Device config not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDeleteHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigDeleteService{
		deleteFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 1, nil // 1 row affected
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/config/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"message": "Device config deleted successfully."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDeleteHandlerNegativeID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockDeviceConfigDeleteService{
		deleteFunc: func(config *models.DeviceConfig, ctx context.Context) (int64, error) {
			return 0, nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/config/-1", nil)
	req.SetPathValue("id", "-1")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
