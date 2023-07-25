package main

import (
	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/handler"
	"github.com/mrrizal/sample-api/middleware"
	"github.com/mrrizal/sample-api/observer"
	"github.com/mrrizal/sample-api/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	observer.InitTracer()
	observer.InitMetric()

	senders := []eventsender.EventSender{
		eventsender.NewAPISender(),
		eventsender.NewKinesisSender(),
		eventsender.NewSQSSender(),
	}
	eventService := service.NewEventService(senders)
	eventHandler := handler.NewEventHandler(eventService)

	app := fiber.New()

	metricMiddleware := middleware.HTTPHandlerWithMetrics
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Post("/v1/api/like", metricMiddleware(eventHandler.Like))
	app.Post("/v1/api/unlike", metricMiddleware(eventHandler.Unlike))

	err := app.Listen(":80")
	if err != nil {
		panic("Error starting the server: " + err.Error())
	}
}
