package pkg

// middleware/metrics.go

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method", "status"})
)

func init() {
	prometheus.MustRegister(requestDuration)
}

// MetricsMiddleware is a middleware for capturing request duration metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&ww, r)

		duration := time.Since(start).Seconds()

		requestDuration.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.status)).Observe(duration)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
