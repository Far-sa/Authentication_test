package ports

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics interface {
	RegisterHTTPDurationHistogram() *prometheus.HistogramVec
	RegisterHTTPErrorCounter() *prometheus.CounterVec
}

// DatabaseMetricsAdapter defines the interface for Prometheus metrics related to database operations.
type DatabaseMetrics interface {
	RegisterDatabaseDurationHistogram() *prometheus.HistogramVec
	RegisterDatabaseErrorCounter() *prometheus.CounterVec
}

// PrometheusAdapter defines the combined interface for Prometheus metrics.
type GoroutineMetrics interface {
	RegisterGoroutineGauge() prometheus.Gauge
	UpdateGoroutineCount()
}
