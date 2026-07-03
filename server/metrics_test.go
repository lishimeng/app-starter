package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsRegistryRegister(t *testing.T) {
	reg := MetricsRegistry()
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_business_events_total",
		Help: "test counter",
	})
	if err := reg.Register(counter); err != nil {
		t.Fatal(err)
	}
	counter.Inc()

	req := httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	w := httptest.NewRecorder()
	MetricsHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "test_business_events_total") {
		t.Fatalf("missing custom metric in output")
	}
}
