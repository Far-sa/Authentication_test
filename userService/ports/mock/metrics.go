package mocks

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

type MockGoroutineMetrics struct {
	mock.Mock
}

func NewMockGoroutineMetrics() *MockGoroutineMetrics {
	return &MockGoroutineMetrics{}
}

func (m *MockGoroutineMetrics) RegisterGoroutineGauge() prometheus.Gauge {
	args := m.Called()
	return args.Get(0).(prometheus.Gauge)
}

func (m *MockGoroutineMetrics) UpdateGoroutineCount() {
	m.Called()
}

// ! --->
type MockHTTPMetrics struct {
	mock.Mock
}

func NewMockHTTPMetrics() *MockHTTPMetrics {
	return &MockHTTPMetrics{}
}

func (m *MockHTTPMetrics) RegisterHTTPDurationHistogram() *prometheus.HistogramVec {
	args := m.Called()
	return args.Get(0).(*prometheus.HistogramVec)
}

func (m *MockHTTPMetrics) RegisterHTTPErrorCounter() *prometheus.CounterVec {
	args := m.Called()
	return args.Get(0).(*prometheus.CounterVec)
}

// !---->
type MockDatabaseMetrics struct {
	mock.Mock
}

func NewMockDatabaseMetrics() *MockDatabaseMetrics {
	return &MockDatabaseMetrics{}
}

func (m *MockDatabaseMetrics) RegisterDatabaseDurationHistogram() *prometheus.HistogramVec {
	args := m.Called()
	return args.Get(0).(*prometheus.HistogramVec)
}

func (m *MockDatabaseMetrics) RegisterDatabaseErrorCounter() *prometheus.CounterVec {
	args := m.Called()
	return args.Get(0).(*prometheus.CounterVec)
}

// MockDatabaseMetricsAdapter is a mock implementation of DatabaseMetricsAdapter for testing
// type MockDatabaseMetricsAdapter struct {
// 	RegisterDatabaseDurationHistogramFn func() *prometheus.HistogramVec
// 	RegisterDatabaseErrorCounterFn      func() *prometheus.CounterVec
// }

// func (m *MockDatabaseMetricsAdapter) RegisterDatabaseDurationHistogram() *prometheus.HistogramVec {
// 	if m.RegisterDatabaseDurationHistogramFn != nil {
// 		return m.RegisterDatabaseDurationHistogramFn()
// 	}
// 	return nil
// }

// func (m *MockDatabaseMetricsAdapter) RegisterDatabaseErrorCounter() *prometheus.CounterVec {
// 	if m.RegisterDatabaseErrorCounterFn != nil {
// 		return m.RegisterDatabaseErrorCounterFn()
// 	}
// 	return nil
// }

// // MockHTTPMetricsAdapter is a mock implementation of HTTPMetricsAdapter for testing
// type MockHTTPMetricsAdapter struct {
// 	RegisterHTTPDurationHistogramFn func() *prometheus.HistogramVec
// 	RegisterHTTPErrorCounterFn      func() *prometheus.CounterVec
// }

// func (m *MockHTTPMetricsAdapter) RegisterHTTPDurationHistogram() *prometheus.HistogramVec {
// 	if m.RegisterHTTPDurationHistogramFn != nil {
// 		return m.RegisterHTTPDurationHistogramFn()
// 	}
// 	return nil
// }

// func (m *MockHTTPMetricsAdapter) RegisterHTTPErrorCounter() *prometheus.CounterVec {
// 	if m.RegisterHTTPErrorCounterFn != nil {
// 		return m.RegisterHTTPErrorCounterFn()
// 	}
// 	return nil
// }

// // MockMetrics is a mock implementation of Metrics for testing
// type MockMetrics struct {
// 	HTTPMetricsAdapter MockHTTPMetricsAdapter
// 	DatabaseMetricsAdapter MockDatabaseMetricsAdapter
// 	RegisterGoroutineGaugeFn func() prometheus.Gauge
// 	UpdateGoroutineCountFn   func()
// }

// func (m *MockMetrics) RegisterGoroutineGauge() prometheus.Gauge {
// 	if m.RegisterGoroutineGaugeFn != nil {
// 		return m.RegisterGoroutineGaugeFn()
// 	}
// 	return nil
// }

// func (m *MockMetrics) UpdateGoroutineCount() {
// 	if m.UpdateGoroutineCountFn != nil {
// 		m.UpdateGoroutineCountFn()
// 	}
// }
