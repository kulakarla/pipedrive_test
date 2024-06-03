package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/deals", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPostHandler(t *testing.T) {
	reqBody := strings.NewReader(`{"title": "nomdea24414"}`)
	req := httptest.NewRequest(http.MethodPost, "/deals", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusCreated, rr.Body.String())
	}
}
