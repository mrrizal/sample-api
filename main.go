package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/model"
	"github.com/mrrizal/sample-api/utils"
)

type EventService struct {
	senders []eventsender.EventSender
}

func NewEventService(senders []eventsender.EventSender) *EventService {
	return &EventService{senders: senders}
}

func (e *EventService) Send(event model.Event) {
	var wg sync.WaitGroup
	for _, sender := range e.senders {
		wg.Add(1)
		go func(sender eventsender.EventSender, wg *sync.WaitGroup) {
			defer wg.Done()
			err := sender.Send(event)
			if err != nil {
				log.Fatal(err.Error())
			}
		}(sender, &wg)
	}
	wg.Wait()
}

type EventHandler struct {
	svc *EventService
}

func NewEventHandler(svc *EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}

func (e *EventHandler) Like(c *fiber.Ctx) error {
	event := utils.GenerateRandomEvent()
	event.Name = "like"

	e.svc.Send(event)
	return c.SendStatus(fiber.StatusNoContent)
}

func (e *EventHandler) Unlike(c *fiber.Ctx) error {
	event := utils.GenerateRandomEvent()
	event.Name = "unlike"

	e.svc.Send(event)
	return c.SendStatus(fiber.StatusNoContent)
}

func main() {
	senders := []eventsender.EventSender{
		eventsender.NewAPISender(),
		eventsender.NewKinesisSender(),
		eventsender.NewSQSSender(),
	}
	eventService := NewEventService(senders)
	eventHandler := NewEventHandler(eventService)

	app := fiber.New()

	app.Post("/v1/api/like", eventHandler.Like)
	app.Post("/v1/api/unlike", eventHandler.Unlike)

	err := app.Listen(":3000")
	if err != nil {
		panic("Error starting the server: " + err.Error())
	}
}
