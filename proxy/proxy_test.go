package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyGetRequestResponseEqualToDirect(t *testing.T) {
	proxyReq := httptest.NewRequest(http.MethodGet, "/deals", nil)
	proxyRr := httptest.NewRecorder()
	handler := http.HandlerFunc(Request)
	handler.ServeHTTP(proxyRr, proxyReq)

	directReq, err := http.NewRequest(http.MethodGet, baseURL+apiTokenParam, nil)

	if err != nil {
		t.Fatalf("Failed to get direct response: %v", err)
	}

	directResp, err := http.DefaultClient.Do(directReq)
	if err != nil {
		t.Fatalf("Failed to get direct response: %v", err)
	}

	defer directResp.Body.Close()
	directBody, err := io.ReadAll(directResp.Body)
	if err != nil {
		t.Fatalf("Failed to read direct response body: %v", err)
	}

	if proxyRr.Body.String() != string(directBody) {
		t.Errorf("Proxy handler returned unexpected body for GET: got %v want %v", proxyRr.Body.String(), directResp)
	}

}

// func TestProxyPostRequestResponseEqualToDirect(t *testing.T) {
// 	reqBody := strings.NewReader(`{"title": "proxyPostTest"}`)
// 	proxyReq := httptest.NewRequest(http.MethodPost, "/deals", reqBody)
// 	proxyReq.Header.Set("Content-Type", "application/json")
// 	proxyRr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Request)
// 	handler.ServeHTTP(proxyRr, proxyReq)

// 	directReq, err := http.NewRequest(http.MethodPost, baseURL+apiTokenParam, reqBody)
// 	directReq.Header.Set("Content-Type", "application/json")

// 	if err != nil {
// 		t.Fatalf("Failed to get direct response: %v", err)
// 	}

// 	directResp, err := http.DefaultClient.Do(directReq)
// 	if err != nil {
// 		t.Fatalf("Failed to get direct response: %v", err)
// 	}

// 	defer directResp.Body.Close()
// 	directBody, err := io.ReadAll(directResp.Body)
// 	if err != nil {
// 		t.Fatalf("Failed to read direct response body: %v", err)
// 	}

// 	if proxyRr.Body.String() != string(directBody) {
// 		t.Errorf("Proxy handler returned unexpected body for POST: got %v want %v", proxyRr.Body.String(), directResp)
// 	}

// 	var response map[string]interface{}
// 	err = json.Unmarshal(proxyRr.Body.Bytes(), &response)

// 	if err != nil {
// 		t.Fatalf("Failed to parse response body: %v", err)
// 	}

// 	if !response["success"].(bool) {
// 		t.Fatalf("API response indicated failure: %v", proxyRr.Body.String())

// 	}

// 	data := response["data"].(map[string]interface{})
// 	id := int(data["id"].(float64))

// 	err = utils.DeleteCreatedResourceInTests(fmt.Sprintf("%d", id))

// 	if err != nil {
// 		t.Fatalf("Failed to delete created resource in proxy POST test: %v", err)
// 	}

// }
