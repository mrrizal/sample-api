package service

import (
	"context"
	"log"
	"sync"

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

func (e *EventService) Send(ctx context.Context, event model.Event) {
	ctx, span := utils.StartTracerSpan(ctx, "EventService/Send")
	defer span.End()

	var wg sync.WaitGroup
	for _, sender := range e.senders {
		wg.Add(1)
		go func(ctx context.Context, sender eventsender.EventSender, wg *sync.WaitGroup) {
			defer wg.Done()
			err := sender.Send(ctx, event)
			if err != nil {
				log.Fatal(err.Error())
			}
		}(ctx, sender, &wg)
	}
	wg.Wait()
}
