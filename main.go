package main

import (
	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/handler"
	"github.com/mrrizal/sample-api/observer"
	"github.com/mrrizal/sample-api/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	// tracker
	observer.InitTracer()

	// Meter
	reg := prometheus.NewRegistry()
	metrics := observer.NewMetrics(reg)

	senders := []eventsender.EventSender{
		eventsender.NewAPISender(),
		eventsender.NewKinesisSender(),
		eventsender.NewSQSSender(),
	}
	eventService := service.NewEventService(senders)
	eventHandler := handler.NewEventHandler(eventService)
	metricHandler := handler.NewMetricHandler(metrics)

	app := fiber.New()

	app.Post("/v1/api/like", metricHandler.HTTPHandlerWithMetrics(eventHandler.Like))
	app.Post("/v1/api/unlike", metricHandler.HTTPHandlerWithMetrics(eventHandler.Unlike))

	app.Get("/metrics", adaptor.HTTPHandler(
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})))

	err := app.Listen(":80")
	if err != nil {
		panic("Error starting the server: " + err.Error())
	}
}
