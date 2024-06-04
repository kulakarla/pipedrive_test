package metrics

import "sync"

//MethodMetrics defines the metric data shown for different request types in the response
type MethodMetrics struct {
	TotalRequests int64   `json:"total_requests"`
	MeanDuration  float64 `json:"mean_duration"`
	TotalDuration int64   `json:"total_duration"`
	MeanLatency   float64 `json:"mean_latency"`
	TotalLatency  int64   `json:"total_latency"`
}

//Metrics defines the overall response body for the metrics request
type Metrics struct {
	GET  MethodMetrics `json:"GET"`
	POST MethodMetrics `json:"POST"`
	PUT  MethodMetrics `json:"PUT"`
}

var (
	metrics Metrics
	mu      sync.Mutex
)

//GetMetrics returns the current metrics
func GetMetrics() Metrics {
	return metrics
}

//UpdateMetrics is an utility function for updating the latency metric
func UpdateMetrics(method string, latency int64) {
	mu.Lock()
	defer mu.Unlock()
	var methodMetrics *MethodMetrics
	switch method {
	case "GET":
		methodMetrics = &metrics.GET
	case "POST":
		methodMetrics = &metrics.POST
	case "PUT":
		methodMetrics = &metrics.PUT
	}

	if methodMetrics != nil {
		methodMetrics.TotalRequests++
		methodMetrics.TotalLatency += latency
		methodMetrics.MeanLatency = float64(methodMetrics.TotalLatency) / float64(methodMetrics.TotalRequests)
	}
}

//UpdateDuration is an utility function for updating the resoponse duration metric
func UpdateDuration(method string, duration int64) {
	mu.Lock()
	defer mu.Unlock()

	var methodMetrics *MethodMetrics
	switch method {
	case "GET":
		methodMetrics = &metrics.GET
	case "POST":
		methodMetrics = &metrics.POST
	case "PUT":
		methodMetrics = &metrics.PUT
	}

	if methodMetrics != nil {
		methodMetrics.TotalDuration += duration
		methodMetrics.MeanDuration = float64(methodMetrics.TotalDuration) / float64(methodMetrics.TotalRequests)
	}
}

//ResetMetrics is an utility function for testing the metrics handler (reset the metrics before running an unit test)
func ResetMetrics() {
	mu.Lock()
	defer mu.Unlock()

	metrics = Metrics{
		GET:  MethodMetrics{},
		POST: MethodMetrics{},
		PUT:  MethodMetrics{},
	}
}
