package metrics

import "sync"

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
	metrics Metrics
	mu      sync.Mutex
)

func GetMetrics() Metrics {
	return metrics
}

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

func ResetMetrics() {
	mu.Lock()
	defer mu.Unlock()

	metrics = Metrics{
		GET:  MethodMetrics{},
		POST: MethodMetrics{},
		PUT:  MethodMetrics{},
	}
}
