package handler

import (
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizal/sample-api/observer"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricHandler struct {
	metrics *observer.Metrics
}

func NewMetricHandler(metrics *observer.Metrics) MetricHandler {
	return MetricHandler{metrics: metrics}
}

func (m *MetricHandler) HTTPHandlerWithMetrics(handler func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Increment the counter for the incoming request. this is for total request
		m.metrics.RequestsTotal.WithLabelValues(c.Route().Method, c.Route().Path).Inc()

		// Start recording the duration of the request processing. this is for duration/historgram
		timer := prometheus.NewTimer(
			m.metrics.DurationHistorgram.WithLabelValues(c.Route().Method, c.Route().Path, ""), // Empty status_code for now.
		)

		// this is for memory usage monitoring
		startMemoryStats := new(runtime.MemStats)
		runtime.ReadMemStats(startMemoryStats)

		// // Execute the actual HTTP handler and wait for the response.
		err := handler(c)

		endMemoryStats := new(runtime.MemStats)
		runtime.ReadMemStats(endMemoryStats)

		allocatedMemory := endMemoryStats.TotalAlloc - startMemoryStats.TotalAlloc

		// Stop the timer after the actual handler has completed its work.
		timer.ObserveDuration()

		// Extract the status code from the response and set it for the histogram.
		statusCode := c.Response().StatusCode()
		m.metrics.DurationHistorgram.WithLabelValues(
			c.Route().Method,
			c.Route().Path,
			fmt.Sprintf("%d", statusCode)).
			Observe(timer.ObserveDuration().Seconds())

		// send memory usage to metrics endpoint
		m.metrics.MemoryUsage.WithLabelValues(c.Route().Method, c.Route().Path).Set(
			float64(allocatedMemory),
		)
		return err
	}
}
