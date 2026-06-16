package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


// =======================
// MOCK REPOSITORY
// =======================
type MockRepo struct{}

func (m *MockRepo) Create(pkg *Package) error {
	return nil
}

func (m *MockRepo) SaveOutbox(event string, data string) error {
	return nil
}


// =======================
// TEST START SORT SUCCESS
// =======================
func TestStartSort_AllValid(t *testing.T) {

	repo := &MockRepo{}
	service := NewSortingService(repo)
	handler := NewSortingHandler(service, repo)

	body := `{"resi":"123","warehouse_zone":"A1","status":"sorting"}`
	req := httptest.NewRequest("POST", "/sort", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.StartSort(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}


// =======================
// TEST ERROR CASES
// =======================
func TestStartSort_AllErrorPaths(t *testing.T) {

	repo := &MockRepo{}
	service := NewSortingService(repo)
	handler := NewSortingHandler(service, repo)

	tests := []string{
		`invalid-json`,
		`{"resi":"","warehouse_zone":"A1","status":"sorting"}`,
		`{"resi":"123","warehouse_zone":"","status":"sorting"}`,
		`{"resi":"123","warehouse_zone":"A1","status":""}`,
		`{"resi":"123","warehouse_zone":"A1","status":"pending"}`,
	}

	for _, body := range tests {
		req := httptest.NewRequest("POST", "/sort", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.StartSort(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}
}


// =======================
// TEST HEALTH
// =======================
func TestHealth_OK(t *testing.T) {

	repo := &MockRepo{}
	service := NewSortingService(repo)
	handler := NewSortingHandler(service, repo)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
