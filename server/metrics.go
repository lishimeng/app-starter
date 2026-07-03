package server

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const DefaultMetricsPath = "/metrics"

var (
	metricsOnce sync.Once
	metricsReg  *prometheus.Registry

	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
)

func initMetrics() {
	metricsOnce.Do(func() {
		metricsReg = prometheus.NewRegistry()
		httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed.",
		}, []string{"method", "route", "status"})
		httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "route", "status"})

		metricsReg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			httpRequestsTotal,
			httpRequestDuration,
		)
	})
}

// MetricsRegistry returns the shared Prometheus registry for custom business metrics.
func MetricsRegistry() *prometheus.Registry {
	initMetrics()
	return metricsReg
}

// MetricsMiddleware records HTTP RED metrics into the shared registry (main web server).
func MetricsMiddleware() gin.HandlerFunc {
	initMetrics()
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())
		labels := []string{c.Request.Method, route, status}
		httpRequestsTotal.WithLabelValues(labels...).Inc()
		httpRequestDuration.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
	}
}

// MetricsHandler serves Prometheus metrics from the shared registry (admin listener).
func MetricsHandler() http.Handler {
	initMetrics()
	return promhttp.HandlerFor(metricsReg, promhttp.HandlerOpts{})
}
