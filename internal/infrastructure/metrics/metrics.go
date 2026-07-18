package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP (RED: Rate, Errors, Duration)

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by method, route template and status code.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "path"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response body size in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6), // 100B -> 10MB
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being served.",
		},
	)

	// Security / auth (attack visibility, low cardinality) 

	SecurityEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_events_total",
			Help: "Security-relevant events: auth_failure, login_failed, account_locked, rate_limit_exceeded, forbidden.",
		},
		[]string{"event", "reason"},
	)

	LoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logins_total",
			Help: "Login attempts by result (success|failure).",
		},
		[]string{"result"},
	)

	PanicsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "panics_recovered_total",
			Help: "Number of panics recovered by the HTTP recovery middleware.",
		},
	)

	// Cache 

	CacheOpsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Cache operations by cache name and result (hit|miss).",
		},
		[]string{"cache", "result"},
	)

	// Build info 

	BuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_build_info",
			Help: "Always 1; labels carry service/version/env for dashboard selection.",
		},
		[]string{"service", "version", "env"},
	)
)

func Register(db *sql.DB, service, version, env string) {
	if db != nil {
		prometheus.MustRegister(collectors.NewDBStatsCollector(db, "mahirlearning"))
	}
	BuildInfo.WithLabelValues(service, version, env).Set(1)
}

func RecordSecurityEvent(event, reason string) {
	SecurityEventsTotal.WithLabelValues(event, reason).Inc()
}

func RecordLogin(result string) {
	LoginsTotal.WithLabelValues(result).Inc()
}

func RecordPanic() {
	PanicsTotal.Inc()
}

func RecordCache(cache, result string) {
	CacheOpsTotal.WithLabelValues(cache, result).Inc()
}
