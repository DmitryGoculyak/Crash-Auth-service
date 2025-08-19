package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Path string `mapstructure:"path"`
}

type MetricsHelper struct {
	HttpRequestTotal    *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
	path                string
}

func NewMetrics(cfg *Config) *MetricsHelper {
	metric := &MetricsHelper{
		HttpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total number of HTTP requests processed",
			},
			[]string{"path", "method", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency distribution",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		path: cfg.Path,
	}

	prometheus.MustRegister(metric.HttpRequestTotal, metric.HttpRequestDuration)
	return metric
}

func (m *MetricsHelper) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

func (m *MetricsHelper) Path() string {
	return m.path
}
