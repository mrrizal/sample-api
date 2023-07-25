package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizal/sample-api/observer"
	"github.com/prometheus/client_golang/prometheus"
)

func HTTPHandlerWithMetrics(handler func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Increment the counter for the incoming request.
		observer.HTTPRequestsTotal.WithLabelValues(c.Method(), c.Path()).Inc()

		// Start recording the duration of the request processing.
		timer := prometheus.NewTimer(
			observer.HTTPDurationHistogram.WithLabelValues(c.Method(), c.Path(), ""), // Empty status_code for now.
		)

		// Execute the actual HTTP handler and wait for the response.
		err := handler(c)

		// Stop the timer after the actual handler has completed its work.
		timer.ObserveDuration()

		// Extract the status code from the response and set it for the histogram.
		statusCode := c.Response().StatusCode()
		observer.HTTPDurationHistogram.WithLabelValues(
			c.Method(),
			c.Path(),
			fmt.Sprintf("%d", statusCode)).
			Observe(timer.ObserveDuration().Seconds())

		return err
	}
}