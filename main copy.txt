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
	http.HandleFunc("/", invalidPathHandler)
	log.Println("Server listening on localhost:8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling %s request for %s\n", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyRequest(rw, req)
	default:
		http.Error(rw, "Method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method %s not allowed for %s\n", req.Method, req.URL.Path)
	}
}

func proxyRequest(rw http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	path := strings.TrimPrefix(req.URL.Path, "/deals")
	targetURL := baseUrl + path + apiTokenParam

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

	updateMetrics(req.Method, latency)
	log.Printf("%s %s Request successfully completed with status %d \n", req.Method, req.URL.Path, proxyResp.StatusCode)

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
	log.Printf("Handling %s request for %s", req.Method, req.URL.Path)
	mu.Lock()
	defer mu.Unlock()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(metrics)

	log.Printf("Successfully returned metrics")
}

func invalidPathHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "Invalid path", http.StatusNotFound)
	log.Printf("Invalid path accessed: %s\n", req.URL.Path)
}
