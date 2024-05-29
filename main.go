package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MethodMetrics struct {
	TotalRequests int64   `json:"total_requests"`
	MeanDuration  float64 `json:"mean_duration"`
	TotalDuration int64   `json:"total_duration"`
	MeanLatency   float64 `json:"mean_latency"`
	TotalLatency  int64   `json:"total_latency"`
}

type Metrics struct {
	GET  MethodMetrics `json:"GET"`
	POST MethodMetrics `json:"POST"`
	PUT  MethodMetrics `json:"PUT"`
}

var (
	metrics       Metrics
	mu            sync.Mutex
	apiToken      = "863be942d8456f146e61026f7cf69dc78efda801"
	apiTokenParam = "?api_token=" + apiToken
	baseUrl       = "https://api.pipedrive.com/v1/deals/"
)

func main() {
	http.HandleFunc("/deals", requestMetricsMiddleware(handler))
	http.HandleFunc("/deals/", requestMetricsMiddleware(handler))
	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyRequest(rw, req)
	default:
		http.Error(rw, "Method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
	}
}

func proxyRequest(rw http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	path := strings.TrimPrefix(req.URL.Path, "/deals")
	targetURL := baseUrl + path + apiTokenParam

	proxyReq, err := http.NewRequest(req.Method, targetURL, req.Body)
	if err != nil {
		http.Error(rw, "Failed to create request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = req.Header.Clone()

	startLatency := time.Now()
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(rw, "Failed to get response from target URL", http.StatusInternalServerError)
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
	}

	updateMetrics(req.Method, latency)

}

func updateMetrics(method string, latency int64) {
	mu.Lock()
	defer mu.Unlock()
	var methodMetrics *MethodMetrics
	switch method {
	case http.MethodGet:
		methodMetrics = &metrics.GET
	case http.MethodPost:
		methodMetrics = &metrics.POST
	case http.MethodPut:
		methodMetrics = &metrics.PUT
	}

	if methodMetrics != nil {
		methodMetrics.TotalRequests++
		methodMetrics.TotalLatency += latency
		methodMetrics.MeanLatency = float64(methodMetrics.TotalLatency) / float64(methodMetrics.TotalRequests)
	}
}

func requestMetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		duration := time.Since(start).Milliseconds()

		mu.Lock()
		defer mu.Unlock()

		var methodMetrics *MethodMetrics
		switch req.Method {
		case http.MethodGet:
			methodMetrics = &metrics.GET
		case http.MethodPost:
			methodMetrics = &metrics.POST
		case http.MethodPut:
			methodMetrics = &metrics.PUT
		}

		if methodMetrics != nil {
			methodMetrics.TotalDuration += duration
			methodMetrics.MeanDuration = float64(methodMetrics.TotalDuration) / float64(methodMetrics.TotalRequests)
		}
	}
}

func metricsHandler(rw http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(metrics)
}
