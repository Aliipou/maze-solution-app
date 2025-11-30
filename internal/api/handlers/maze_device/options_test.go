package maze_device

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/device/status", nil)
	w := httptest.NewRecorder()

	OptionsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body, got %s", w.Body.String())
	}
}

func TestOptionsHandlerCORS(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/device/status", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()

	OptionsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestOptionsHandlerMultipleRequests(t *testing.T) {
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodOptions, "/device/status", nil)
		w := httptest.NewRecorder()

		OptionsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i, w.Code)
		}
	}
}

func TestOptionsHandlerWithQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/device/status?page=1&rows_per_page=10", nil)
	w := httptest.NewRecorder()

	OptionsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestOptionsHandlerDifferentPaths(t *testing.T) {
	paths := []string{
		"/device/status",
		"/device/status/1",
		"/device/status?device_id=ESP32_001",
	}

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		w := httptest.NewRecorder()

		OptionsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Path %s: Expected status 200, got %d", path, w.Code)
		}
	}
}
