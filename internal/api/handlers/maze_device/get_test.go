package maze_device

import (
	"context"
	"encoding/json"
	"errors"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/maze_device"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Mock service for GET testing
type mockMazeDeviceGetService struct {
	readManyFunc       func(int, int, context.Context) ([]*models.MazeDeviceStatus, error)
	readByDeviceIDFunc func(string, context.Context) ([]*models.MazeDeviceStatus, error)
}

func (m *mockMazeDeviceGetService) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	return nil
}

func (m *mockMazeDeviceGetService) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceGetService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	if m.readManyFunc != nil {
		return m.readManyFunc(page, rowsPerPage, ctx)
	}
	return nil, nil
}

func (m *mockMazeDeviceGetService) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	if m.readByDeviceIDFunc != nil {
		return m.readByDeviceIDFunc(deviceID, ctx)
	}
	return nil, nil
}

func (m *mockMazeDeviceGetService) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDeviceGetService) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDeviceGetService) ValidateStatus(status *models.MazeDeviceStatus) error {
	return nil
}

func TestGetHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedStatuses := []*models.MazeDeviceStatus{
		{
			ID:              1,
			DeviceID:        "ESP32_001",
			AlarmActive:     true,
			MazeCompleted:   false,
			HallSensorValue: true,
			BatteryLevel:    95,
			Timestamp:       "2024-01-15T10:00:00Z",
		},
		{
			ID:              2,
			DeviceID:        "ESP32_002",
			AlarmActive:     false,
			MazeCompleted:   true,
			HallSensorValue: false,
			BatteryLevel:    87,
			Timestamp:       "2024-01-15T11:00:00Z",
		},
	}

	mockService := &mockMazeDeviceGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			return expectedStatuses, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.MazeDeviceStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(response))
	}
}

func TestGetHandlerWithPagination(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDeviceGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			// Verify pagination parameters are passed correctly
			if page != 2 || rowsPerPage != 20 {
				t.Errorf("Expected page=2, rowsPerPage=20, got page=%d, rowsPerPage=%d", page, rowsPerPage)
			}
			return []*models.MazeDeviceStatus{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status?page=2&rows_per_page=20", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetHandlerByDeviceID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedStatuses := []*models.MazeDeviceStatus{
		{
			ID:              1,
			DeviceID:        "ESP32_TEST",
			AlarmActive:     true,
			MazeCompleted:   false,
			HallSensorValue: true,
			BatteryLevel:    100,
			Timestamp:       "2024-01-15T10:00:00Z",
		},
	}

	mockService := &mockMazeDeviceGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			if deviceID != "ESP32_TEST" {
				t.Errorf("Expected device_id ESP32_TEST, got %s", deviceID)
			}
			return expectedStatuses, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status?device_id=ESP32_TEST", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.MazeDeviceStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 status, got %d", len(response))
	}

	if response[0].DeviceID != "ESP32_TEST" {
		t.Errorf("Expected DeviceID ESP32_TEST, got %s", response[0].DeviceID)
	}
}

func TestGetHandlerByDeviceIDError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDeviceGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			return nil, maze_device.MazeDeviceStatusError{Message: "invalid device_id"}
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status?device_id=INVALID", nil)
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

	mockService := &mockMazeDeviceGetService{
		readByDeviceIDFunc: func(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			return nil, errors.New("database connection error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status?device_id=ESP32_001", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDeviceGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			return nil, errors.New("database error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetHandlerEmptyResult(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDeviceGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			return []*models.MazeDeviceStatus{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.MazeDeviceStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected empty array, got %d items", len(response))
	}
}

func TestGetHandlerInvalidPaginationParameters(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDeviceGetService{
		readManyFunc: func(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
			// Invalid params should be converted to 0
			if page != 0 || rowsPerPage != 0 {
				t.Errorf("Expected page=0, rowsPerPage=0 for invalid params, got page=%d, rowsPerPage=%d", page, rowsPerPage)
			}
			return []*models.MazeDeviceStatus{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/device/status?page=invalid&rows_per_page=also_invalid", nil)
	w := httptest.NewRecorder()

	GetHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
