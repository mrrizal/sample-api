package main

import (
	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/handler"
	"github.com/mrrizal/sample-api/observer"
	"github.com/mrrizal/sample-api/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	observer.InitTracer()

	senders := []eventsender.EventSender{
		eventsender.NewAPISender(),
		eventsender.NewKinesisSender(),
		eventsender.NewSQSSender(),
	}
	eventService := service.NewEventService(senders)
	eventHandler := handler.NewEventHandler(eventService)

	app := fiber.New()

	app.Post("/v1/api/like", eventHandler.Like)
	app.Post("/v1/api/unlike", eventHandler.Unlike)

	err := app.Listen(":80")
	if err != nil {
		panic("Error starting the server: " + err.Error())
	}
}
