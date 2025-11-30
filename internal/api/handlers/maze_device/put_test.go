package maze_device

import (
	"bytes"
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

// Mock service for PUT testing
type mockMazeDevicePutService struct {
	updateFunc func(*models.MazeDeviceStatus, context.Context) (int64, error)
}

func (m *mockMazeDevicePutService) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	return nil
}

func (m *mockMazeDevicePutService) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDevicePutService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDevicePutService) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDevicePutService) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	if m.updateFunc != nil {
		return m.updateFunc(status, ctx)
	}
	return 0, nil
}

func (m *mockMazeDevicePutService) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDevicePutService) ValidateStatus(status *models.MazeDeviceStatus) error {
	return nil
}

func TestPutHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDevicePutService{
		updateFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 1, nil // 1 row affected
		},
	}

	status := models.MazeDeviceStatus{
		ID:              1,
		DeviceID:        "ESP32_001",
		AlarmActive:     false,
		MazeCompleted:   true,
		HallSensorValue: true,
		BatteryLevel:    85,
		Timestamp:       "2024-01-15T11:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.MazeDeviceStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}
}

func TestPutHandlerInvalidJSON(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDevicePutService{}

	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBufferString("{invalid json}"))
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
	mockService := &mockMazeDevicePutService{}

	status := models.MazeDeviceStatus{
		// ID is 0 (missing)
		DeviceID:        "ESP32_001",
		AlarmActive:     true,
		MazeCompleted:   false,
		HallSensorValue: true,
		BatteryLevel:    95,
		Timestamp:       "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
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

	mockService := &mockMazeDevicePutService{
		updateFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, maze_device.MazeDeviceStatusError{Message: "device_id is required"}
		},
	}

	status := models.MazeDeviceStatus{
		ID:              1,
		DeviceID:        "", // Invalid
		AlarmActive:     true,
		MazeCompleted:   false,
		HallSensorValue: true,
		BatteryLevel:    95,
		Timestamp:       "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
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

	mockService := &mockMazeDevicePutService{
		updateFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, nil // 0 rows affected
		},
	}

	status := models.MazeDeviceStatus{
		ID:              999,
		DeviceID:        "ESP32_NOTFOUND",
		AlarmActive:     true,
		MazeCompleted:   false,
		HallSensorValue: true,
		BatteryLevel:    95,
		Timestamp:       "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Maze device status not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestPutHandlerInternalError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDevicePutService{
		updateFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, errors.New("database connection error")
		},
	}

	status := models.MazeDeviceStatus{
		ID:              1,
		DeviceID:        "ESP32_001",
		AlarmActive:     true,
		MazeCompleted:   false,
		HallSensorValue: true,
		BatteryLevel:    95,
		Timestamp:       "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestPutHandlerEmptyBody(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDevicePutService{}

	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPutHandlerBoundaryValues(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mockService := &mockMazeDevicePutService{
		updateFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			// Verify boundary values
			if status.BatteryLevel != 0 {
				t.Errorf("Expected BatteryLevel 0, got %d", status.BatteryLevel)
			}
			return 1, nil
		},
	}

	status := models.MazeDeviceStatus{
		ID:              1,
		DeviceID:        "ESP32_001",
		AlarmActive:     false,
		MazeCompleted:   true,
		HallSensorValue: false,
		BatteryLevel:    0, // Boundary value
		Timestamp:       "2024-01-15T10:00:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PutHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
