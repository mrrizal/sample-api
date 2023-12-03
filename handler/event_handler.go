package handler

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizal/sample-api/service"
	"github.com/mrrizal/sample-api/utils"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}

func (e *EventHandler) Like(c *fiber.Ctx) error {
	ctx, span := utils.StartTracerSpan(c.Context(), "/like")
	defer span.End()

	event := utils.GenerateRandomEvent()
	event.Name = "like"

	e.svc.Send(ctx, event)
	return c.SendStatus(fiber.StatusNoContent)
}

func (e *EventHandler) Unlike(c *fiber.Ctx) error {
	ctx, span := utils.StartTracerSpan(c.Context(), "/unlike")
	defer span.End()

	event := utils.GenerateRandomEvent()
	event.Name = "unlike"

	e.svc.Send(ctx, event)
	return c.SendStatus(fiber.StatusNoContent)
}

func (e *EventHandler) DummyEndpoint(c *fiber.Ctx) error {
	_, span := utils.StartTracerSpan(c.Context(), "/unlike")
	defer span.End()

	// Generate a random size between 1MB and 10MB
	minSize := 1 * 1024 * 1024  // 1MB
	maxSize := 10 * 1024 * 1024 // 10MB
	size := rand.Intn(maxSize-minSize+1) + minSize

	data := make([]byte, size)

	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (e *EventHandler) DummyEndpoint1(c *fiber.Ctx) error {
	_, span := utils.StartTracerSpan(c.Context(), "/unlike")
	defer span.End()

	// Generate a random size between 1MB and 10MB
	minSize := 1 * 1024 * 1024 // 1MB
	maxSize := 5 * 1024 * 1024 // 5MB
	size := rand.Intn(maxSize-minSize+1) + minSize

	data := make([]byte, size)

	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
