package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pdrive/pipedrive-test-api/metrics"
	"pdrive/pipedrive-test-api/utils"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/deals", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPostHandler(t *testing.T) {
	reqBody := strings.NewReader(`{"title": "testDeal"}`)
	req := httptest.NewRequest(http.MethodPost, "/deals", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v, body: %v", status, http.StatusCreated, rr.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)

	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !response["success"].(bool) {
		t.Fatalf("API response indicated failure: %v", rr.Body.String())

	}

	data := response["data"].(map[string]interface{})
	id := int(data["id"].(float64))
	title := data["title"].(string)

	if title != "testDeal" {
		t.Fatalf("Created resource title: got %v want %v", title, "testDeal")
	}

	err = utils.DeleteCreatedResourceInTests(fmt.Sprintf("%d", id))

	if err != nil {
		t.Fatalf("Failed to delete created resource in POST test: %v", err)
	}

}

func TestPutHandler(t *testing.T) {
	reqBodyPost := strings.NewReader(`{"title": "testDealPut"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/deals", reqBodyPost)
	reqPost.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, reqPost)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v, body: %v", status, http.StatusCreated, rr.Body.String())
	}

	var responsePost map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &responsePost)

	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !responsePost["success"].(bool) {
		t.Fatalf("API response indicated failure for POST: %v", rr.Body.String())

	}

	dataPost := responsePost["data"].(map[string]interface{})
	idPost := int(dataPost["id"].(float64))

	reqBodyPut := strings.NewReader(`{"currency": "NZD"}`)
	reqPut := httptest.NewRequest(http.MethodPut, "/deals/"+fmt.Sprintf("%d", idPost), reqBodyPut)
	reqPut.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, reqPut)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	var responsePut map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsePut)

	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !responsePut["success"].(bool) {
		t.Fatalf("API response indicated failure for PUT: %v", rr.Body.String())

	}

	dataPut := responsePut["data"].(map[string]interface{})
	idPut := int(dataPut["id"].(float64))
	titlePut := dataPut["title"].(string)
	currencyPut := dataPut["currency"].(string)

	if idPost != idPut && titlePut != "testDealPut" && currencyPut != "NZD" {
		t.Fatalf("PUT method for deal failed, deal id: %d, currency got: %s, want: %s", idPost, currencyPut, "NZD")
	}

	err = utils.DeleteCreatedResourceInTests(fmt.Sprintf("%d", idPost))

	if err != nil {
		t.Fatalf("Failed to delete created resource in POST test: %v", err)
	}

}

func TestMetricsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	metricsHandler := http.HandlerFunc(MetricsHandler)
	dealsHandler := http.HandlerFunc(Handler)
	metrics.ResetMetrics()

	metricsHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v body: %v", status, http.StatusOK, rr.Body.String())
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	for i := 0; i < 2; i++ {
		reqGet := httptest.NewRequest(http.MethodGet, "/deals", nil)
		dealsHandler.ServeHTTP(rr, reqGet)

	}

	reqPost := httptest.NewRequest(http.MethodPost, "/deals", nil)
	dealsHandler.ServeHTTP(rr, reqPost)

	reqPut := httptest.NewRequest(http.MethodPut, "/deals/4", nil)
	dealsHandler.ServeHTTP(rr, reqPut)

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr = httptest.NewRecorder()
	metricsHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code on GET /metrics: got %v want %v", status, http.StatusOK)
	}

	var metricsResponse map[string]map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &metricsResponse)

	if err != nil {
		t.Fatalf("Handler returned invalid JSON: %v", err)
	}

	if metricsResponse["GET"]["total_requests"].(float64) != 2 {
		t.Errorf("Metrics count mismatch for GET: got %v want %v", metricsResponse["GET"]["total_requests"], 2)
	}

	if metricsResponse["POST"]["total_requests"].(float64) != 1 {
		t.Errorf("Metrics count mismatch for POST: got %v want %v", metricsResponse["POST"]["total_requests"], 1)
	}

	if metricsResponse["PUT"]["total_requests"].(float64) != 1 {
		t.Errorf("Metrics count mismatch for PUT: got %v want %v", metricsResponse["PUT"]["total_requests"], 1)
	}

}
