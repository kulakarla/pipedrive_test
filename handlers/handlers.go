package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"pdrive/pipedrive-test-api/metrics"
	"pdrive/pipedrive-test-api/proxy"
)

var mu sync.Mutex

func Handler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling %s request for %s\n", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxy.ProxyRequest(rw, req)
	default:
		http.Error(rw, "Method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method %s not allowed for %s\n", req.Method, req.URL.Path)
	}
}

func MetricsHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling %s request for %s", req.Method, req.URL.Path)
	mu.Lock()
	defer mu.Unlock()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(metrics.GetMetrics())

	log.Printf("Successfully returned metrics")
}

func InvalidPathHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "Invalid path", http.StatusNotFound)
	log.Printf("Invalid path accessed: %s\n", req.URL.Path)
}

func RequestMetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		duration := time.Since(start).Milliseconds()

		metrics.UpdateDuration(req.Method, duration)
	}
}
