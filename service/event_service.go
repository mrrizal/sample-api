package service

import (
	"context"

	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/model"
	"github.com/mrrizal/sample-api/utils"
)

type EventService struct {
	senders          []eventsender.EventSender
	senderProcessing eventsender.EventSenderProcessing
}

func NewEventService(senders []eventsender.EventSender,
	senderProcessing eventsender.EventSenderProcessing) *EventService {
	return &EventService{
		senders:          senders,
		senderProcessing: senderProcessing,
	}
}

func (e *EventService) Send(ctx context.Context, event model.Event) {
	ctx, span := utils.StartTracerSpan(ctx, "EventService/Send")
	defer span.End()

	e.senderProcessing.Send(ctx, e.senders, event)
}
