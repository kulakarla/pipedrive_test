package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyGetRequestResponseEqualToOriginal(t *testing.T) {
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
