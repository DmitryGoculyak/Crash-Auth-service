package middleware

import (
	"Crash-Auth-service/pkg/metrics"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Metrics(m *metrics.MetricsHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		m.HttpRequestTotal.WithLabelValues(
			c.FullPath(),
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
		).Inc()
		m.HttpRequestDuration.WithLabelValues(
			c.FullPath(),
			c.Request.Method,
		).Observe(duration)
	}
}
