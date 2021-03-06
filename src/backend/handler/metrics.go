package handler

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "apiserver_request_count",
			Help: "Counter of apiserver requests broken out for each verb, API resource, client, and HTTP response contentType and code.",
		},
		[]string{"verb", "resource", "client", "contentType", "code"},
	)
	requestLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "apiserver_request_latencies",
			Help: "Response latency distribution in microseconds for each verb, resource and client.",
			// Use buckets ranging from 125 ms to 8 seconds.
			Buckets: prometheus.ExponentialBuckets(125000, 2.0, 7),
		},
		[]string{"verb", "resource"},
	)
	requestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "apiserver_request_latencies_summary",
			Help: "Response latency summary in microseconds for each verb and resource.",
			// Make the sliding window of 1h.
			MaxAge: time.Hour,
		},
		[]string{"verb", "resource"},
	)
)

// Initialize all metrics in prometheus
func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestLatencies)
	prometheus.MustRegister(requestLatenciesSummary)
}

// Track API call in prometheus
func monitor(verb, resource string, client, contentType string, httpCode int, reqStart time.Time) {
	elapsed := float64((time.Since(reqStart)) / time.Microsecond)
	requestCounter.WithLabelValues(verb, resource, client, contentType, strconv.Itoa(httpCode)).Inc()
	requestLatencies.WithLabelValues(verb, resource).Observe(elapsed)
	requestLatenciesSummary.WithLabelValues(verb, resource).Observe(elapsed)
}
