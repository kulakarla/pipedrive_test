package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"pdrive/pipedrive-test-api/config"
	"testing"
)

// This tests checks that for a request to the proxy API, the response body-s are equal when using the proxy (GET request check) or directly calling the PipeDrive API and that the
// header key-s are the same
func TestProxyRequestResponseEqualToDirect(t *testing.T) {
	proxyReq := httptest.NewRequest(http.MethodGet, "/deals", nil)
	proxyRr := httptest.NewRecorder()
	handler := http.HandlerFunc(Request)
	handler.ServeHTTP(proxyRr, proxyReq)

	directReq, err := http.NewRequest(http.MethodGet, config.BaseURL+config.APITokenParam, nil)

	if err != nil {
		t.Fatalf("Failed to create direct request: %v", err)
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

	for key, values := range directResp.Header {
		proxyValues := proxyRr.Header().Values(key)
		if len(values) != len(proxyValues) {
			t.Errorf("Proxy handler returned unexpected number of header values for %s: got %v want %v", key, proxyValues, values)
		}
	}

	for key, values := range proxyRr.Header() {
		if _, ok := directResp.Header[key]; !ok {
			t.Errorf("Proxy handler returned unexpected header %s: got %v", key, values)
		}
	}

}
