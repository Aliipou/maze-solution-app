package maze_device

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
type mockMazeDeviceDeleteService struct {
	deleteFunc func(*models.MazeDeviceStatus, context.Context) (int64, error)
}

func (m *mockMazeDeviceDeleteService) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	return nil
}

func (m *mockMazeDeviceDeleteService) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceDeleteService) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceDeleteService) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return nil, nil
}

func (m *mockMazeDeviceDeleteService) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockMazeDeviceDeleteService) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return m.deleteFunc(status, ctx)
}

func (m *mockMazeDeviceDeleteService) ValidateStatus(status *models.MazeDeviceStatus) error {
	return nil
}

func TestDeleteHandlerInvalidID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceDeleteService{}

	req := httptest.NewRequest(http.MethodDelete, "/device/status/invalid", nil)
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
	mockService := &mockMazeDeviceDeleteService{
		deleteFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, errors.New("database error")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/status/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestDeleteHandlerNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceDeleteService{
		deleteFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, nil // 0 rows affected
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/status/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	expected := `{"error": "Maze device status not found."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDeleteHandlerSuccess(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceDeleteService{
		deleteFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 1, nil // 1 row affected
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/status/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"message": "Maze device status deleted successfully."}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestDeleteHandlerZeroID(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mockService := &mockMazeDeviceDeleteService{
		deleteFunc: func(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
			return 0, nil // Not found
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/device/status/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	DeleteHandler(w, req, logger, mockService)

	// Zero is a valid integer, so it should proceed with deletion
	// The behavior depends on whether the service finds it
	if w.Code != http.StatusNotFound && w.Code != http.StatusOK {
		t.Errorf("Expected status 404 or 200, got %d", w.Code)
	}
}
