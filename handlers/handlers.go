package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"pdrive/pipedrive-test-api/metrics"
	"pdrive/pipedrive-test-api/proxy"
)

var mu sync.Mutex

// Handler takes in the request, forwards it to the proxy if it is an allowed request
func Handler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling %s request for %s\n", req.Method, req.URL.Path)
	allowed := false

	getPostDealsPath := regexp.MustCompile(`^/deals/?$`)
	putDealsPath := regexp.MustCompile(`^/deals/\d+/?$`)

	switch req.Method {
	case http.MethodGet, http.MethodPost:
		if getPostDealsPath.MatchString(req.URL.Path) {
			allowed = true
		}
	case http.MethodPut:
		if putDealsPath.MatchString(req.URL.Path) {
			allowed = true
		}

	}

	if allowed {
		proxy.Request(rw, req)
	} else {
		http.Error(rw, "Method "+req.Method+" not allowed for path "+req.URL.Path, http.StatusMethodNotAllowed)
		log.Printf("Method %s not allowed for %s\n", req.Method, req.URL.Path)
	}
}

// MetricsHandler returns the GET /metrics endpoint
func MetricsHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling %s request for %s", req.Method, req.URL.Path)
	mu.Lock()
	defer mu.Unlock()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(metrics.GetMetrics())

	log.Printf("Successfully returned metrics")
}

// InvalidPathHandler returns an error for all invalid requests
func InvalidPathHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "Invalid path", http.StatusNotFound)
	log.Printf("Invalid path accessed: %s\n", req.URL.Path)
}

// RequestMetricsMiddleware is an utility tool / wrapper for gathering metric data
func RequestMetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		duration := time.Since(start).Milliseconds()
		metrics.UpdateDuration(req.Method, duration)
	}
}
