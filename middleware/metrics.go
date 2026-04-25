package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/bolean304/e-commerce-cart/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		endpoint := c.FullPath()

		metrics.HttpRequests.WithLabelValues(
			c.Request.Method,
			endpoint,
			status,
		).Inc()

		metrics.HttpDuration.WithLabelValues(endpoint).Observe(duration)
	}
}