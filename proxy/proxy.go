package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"pdrive/pipedrive-test-api/config"
	"pdrive/pipedrive-test-api/metrics"
)

// Request forwards the request to the PipeDrive API and returns its response
func Request(rw http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	path := strings.TrimPrefix(req.URL.Path, "/deals")
	targetURL := config.BaseURL + path + config.APITokenParam

	log.Printf("Proxy request handling started for %s %s", req.Method, req.URL.Path)
	proxyReq, err := http.NewRequest(req.Method, targetURL, req.Body)
	if err != nil {
		http.Error(rw, "Failed to create request", http.StatusInternalServerError)
		log.Printf("Failed to create request for %s: %v\n", targetURL, err)
		return
	}

	proxyReq.Header = req.Header.Clone()

	startLatency := time.Now()
	proxyResp, err := client.Do(proxyReq)

	if err != nil {
		http.Error(rw, "Failed to get response from target URL", http.StatusInternalServerError)
		log.Printf("Failed to get response from %s: %v\n", targetURL, err)
		return
	}

	latency := time.Since(startLatency).Milliseconds()

	defer proxyResp.Body.Close()

	for name, values := range proxyResp.Header {
		for _, value := range values {
			rw.Header().Add(name, value)
		}
	}

	rw.WriteHeader(proxyResp.StatusCode)

	_, err = io.Copy(rw, proxyResp.Body)
	if err != nil {
		http.Error(rw, "Failed to copy response body", http.StatusInternalServerError)
		log.Printf("Failed to copy response body for %s: %v\n", targetURL, err)
	}

	metrics.UpdateLatencyAndTotalRequestsMetrics(req.Method, latency)
	log.Printf("%s %s request successfully completed with status %d \n", req.Method, req.URL.Path, proxyResp.StatusCode)
}
