package metrics

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestErrors   *prometheus.CounterVec
	DatabaseDuration    *prometheus.HistogramVec
	DatabaseErrors      *prometheus.CounterVec
	Goroutines          prometheus.Gauge
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "Duration of HTTP requests in seconds.",
			},
			[]string{"method", "handler"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_errors_total",
				Help: "Total count of HTTP request errors.",
			},
			[]string{"method", "handler"},
		),
		DatabaseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "database_request_duration_seconds",
				Help: "Duration of database requests in seconds.",
			},
			[]string{"query"},
		),
		DatabaseErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_request_errors_total",
				Help: "Total count of database request errors.",
			},
			[]string{"query"},
		),
		Goroutines: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "goroutines",
				Help: "Number of active goroutines.",
			},
		),
	}
}

func (pa *Prometheus) RegisterHTTPDurationHistogram() *prometheus.HistogramVec {
	prometheus.MustRegister(pa.HTTPRequestDuration)
	return pa.HTTPRequestDuration
}

func (pa *Prometheus) RegisterHTTPErrorCounter() *prometheus.CounterVec {
	prometheus.MustRegister(pa.HTTPRequestErrors)
	return pa.HTTPRequestErrors
}

func (pa *Prometheus) RegisterDatabaseDurationHistogram() *prometheus.HistogramVec {
	prometheus.MustRegister(pa.DatabaseDuration)
	return pa.DatabaseDuration
}

func (pa *Prometheus) RegisterDatabaseErrorCounter() *prometheus.CounterVec {
	prometheus.MustRegister(pa.DatabaseErrors)
	return pa.DatabaseErrors
}

func (pa *Prometheus) RegisterGoroutineGauge() prometheus.Gauge {
	prometheus.MustRegister(pa.Goroutines)
	return pa.Goroutines
}

func (pa *Prometheus) UpdateGoroutineCount() {
	pa.Goroutines.Set(float64(runtime.NumGoroutine()))
}

//! main
// Initialize Prometheus adapter
// prometheus := metrics.Newprometheus()

// // Initialize HTTP handler with Prometheus metrics adapter
// handler := http.NewHandler(prometheus)

// // Initialize database repository with Prometheus metrics adapter
// repo := repository.NewRepository(prometheus)

// // Register routes
// http.HandleFunc("/hello", handler.HelloHandler)
// http.Handle("/metrics", prometheus.HTTPHandler())

// // Start HTTP server
// go func() {
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         panic(err)
//     }
// }()

// // Update goroutine count periodically
// go func() {
//     ticker := time.NewTicker(30 * time.Second)
//     defer ticker.Stop()

//     for range ticker.C {
//         prometheus.UpdateGoroutineCount()
//     }
// }()

// // Wait forever
// select {}
