package handler

import (
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
