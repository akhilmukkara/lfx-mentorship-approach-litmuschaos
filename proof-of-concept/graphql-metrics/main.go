package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	graphqlRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "litmus_graphql_requests_total",
			Help: "Total GraphQL requests by operation and status",
		},
		[]string{"operation", "status"},
	)

	graphqlDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "litmus_graphql_request_duration_seconds",
			Help:    "GraphQL request duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"operation"},
	)

	graphqlErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "litmus_graphql_errors_total",
			Help: "Total GraphQL errors by operation",
		},
		[]string{"operation", "error_type"},
	)

	activeSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "litmus_active_sessions",
			Help: "Number of active user sessions",
		},
	)
)

// GraphQL operations
var operations = []string{
	"listChaosEngines",
	"getChaosEngine",
	"createChaosEngine",
	"listWorkflows",
	"getWorkflow",
}

// GraphQL handler with metrics
func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	operation := r.URL.Query().Get("operation")
	if operation == "" {
		operation = "unknown"
	}

	start := time.Now()

	// Simulate processing time
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	// Simulate 10% error rate
	if rand.Float32() < 0.1 {
		duration := time.Since(start).Seconds()
		graphqlDuration.WithLabelValues(operation).Observe(duration)
		graphqlRequests.WithLabelValues(operation, "error").Inc()
		graphqlErrors.WithLabelValues(operation, "internal_error").Inc()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// Success
	duration := time.Since(start).Seconds()
	graphqlDuration.WithLabelValues(operation).Observe(duration)
	graphqlRequests.WithLabelValues(operation, "success").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]string{
			"operation": operation,
			"status":    "success",
		},
	})
}

// Home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>LitmusChaos GraphQL Metrics PoC</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #5B21B6; }
        .endpoint { background: #f0f0f0; padding: 15px; margin: 15px 0; border-radius: 5px; border-left: 4px solid #5B21B6; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; font-family: monospace; }
        a { color: #5B21B6; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .metrics-badge { background: #10B981; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>LitmusChaos GraphQL Metrics - Proof of Concept</h1>
        <p>This PoC demonstrates how Prometheus metrics would be instrumented in the LitmusChaos GraphQL server.</p>
        
        <h2>Metrics Exposed</h2>
        <div class="endpoint">
            <strong>litmus_graphql_requests_total</strong> <span class="metrics-badge">Counter</span><br>
            Total GraphQL requests by operation and status (success/error)
        </div>
        <div class="endpoint">
            <strong>litmus_graphql_request_duration_seconds</strong> <span class="metrics-badge">Histogram</span><br>
            Request duration in seconds by operation
        </div>
        <div class="endpoint">
            <strong>litmus_graphql_errors_total</strong> <span class="metrics-badge">Counter</span><br>
            Total errors by operation and error type
        </div>
        <div class="endpoint">
            <strong>litmus_active_sessions</strong> <span class="metrics-badge">Gauge</span><br>
            Number of active user sessions
        </div>
        
        <h2>Try it:</h2>
        <ul>
            <li><a href="/graphql?operation=listChaosEngines">List Chaos Engines</a></li>
            <li><a href="/graphql?operation=getChaosEngine">Get Chaos Engine</a></li>
            <li><a href="/graphql?operation=createChaosEngine">Create Chaos Engine</a></li>
            <li><a href="/graphql?operation=listWorkflows">List Workflows</a></li>
            <li><a href="/metrics" target="_blank"><strong>View Prometheus Metrics</strong></a></li>
        </ul>
        
        <p style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #666;">
            <strong>Note:</strong> This is a standalone PoC. The actual implementation would be integrated into the LitmusChaos GraphQL server handlers.
        </p>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

// Background traffic generator
func trafficGenerator() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		operation := operations[rand.Intn(len(operations))]
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/graphql?operation=%s", operation))
		if err == nil {
			resp.Body.Close()
		}
		activeSessions.Set(float64(10 + rand.Intn(40)))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go trafficGenerator()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/graphql", graphqlHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("LitmusChaos GraphQL Metrics PoC")
	fmt.Println("Server: http://localhost:8080")
	fmt.Println("Metrics: http://localhost:8080/metrics")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
