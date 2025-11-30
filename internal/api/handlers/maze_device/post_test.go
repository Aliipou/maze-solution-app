package maze_device

import (
	"bytes"
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/maze_device"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Mock service for testing
type mockMazeDeviceStatusService struct {
	createFunc func(*models.MazeDeviceStatus, context.Context) error
}

func (m *mockMazeDeviceStatusService) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	return m.createFunc(status, ctx)
}

func (m *mockMazeDeviceStatusService) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceStatusService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceStatusService) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceStatusService) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDeviceStatusService) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDeviceStatusService) ValidateStatus(status *models.MazeDeviceStatus) error {
	return nil
}

func TestPostHandlerInvalidJSON(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceStatusService{}

	req := httptest.NewRequest(http.MethodPost, "/device/status", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostHandlerValidationError(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceStatusService{
		createFunc: func(status *models.MazeDeviceStatus, ctx context.Context) error {
			return maze_device.MazeDeviceStatusError{Message: "device_id is required"}
		},
	}

	status := models.MazeDeviceStatus{
		DeviceID:       "",
		AlarmActive:    true,
		MazeCompleted:  false,
		HallSensorValue: false,
		BatteryLevel:   85,
		Timestamp:      "2024-01-15T10:30:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPost, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceStatusService{
		createFunc: func(status *models.MazeDeviceStatus, ctx context.Context) error {
			status.ID = 1
			return nil
		},
	}

	status := models.MazeDeviceStatus{
		DeviceID:        "ARD001",
		AlarmActive:     true,
		MazeCompleted:   false,
		HallSensorValue: false,
		BatteryLevel:    85,
		Timestamp:       "2024-01-15T10:30:00Z",
	}

	jsonData, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPost, "/device/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostHandler(w, req, logger, mockService)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.MazeDeviceStatus
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}
}
