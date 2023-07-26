package eventsender

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mrrizal/sample-api/model"
	"github.com/mrrizal/sample-api/utils"
)

type EventSender interface {
	Send(ctx context.Context, event model.Event) error
}

type KinesisSender struct{}

func NewKinesisSender() *KinesisSender {
	return &KinesisSender{}
}

func (k *KinesisSender) Send(ctx context.Context, event model.Event) error {
	_, span := utils.StartTracerSpan(ctx, "KinesisSender/Send")
	defer span.End()

	if event.Name == "like" {
		time.Sleep(utils.RandomDuration(300, 700) * time.Millisecond)
	} else {
		time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	}

	log.Println(utils.GetMessageTemplate("kinesis", &event))
	return nil
}

type SQSSender struct{}

func NewSQSSender() *SQSSender {
	return &SQSSender{}
}

func (s *SQSSender) Send(ctx context.Context, event model.Event) error {
	_, span := utils.StartTracerSpan(ctx, "SQSSender/Send")
	defer span.End()

	if event.Name == "like" {
		time.Sleep(utils.RandomDuration(300, 700) * time.Millisecond)
	} else {
		time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	}

	log.Println(utils.GetMessageTemplate("sqs", &event))
	return nil
}

type APISender struct{}

func NewAPISender() *APISender {
	return &APISender{}
}

func (a *APISender) Send(ctx context.Context, event model.Event) error {
	_, span := utils.StartTracerSpan(ctx, "APISender/Send")
	defer span.End()

	if event.Name == "like" {
		time.Sleep(utils.RandomDuration(300, 700) * time.Millisecond)
	} else {
		time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	}

	log.Println(utils.GetMessageTemplate("api", &event))
	return nil
}

type EventData struct {
	Ctx     context.Context
	Senders []EventSender
	Event   model.Event
}

type EventSenderProcessing struct {
	queue chan EventData
	wg    sync.WaitGroup
}

func NewEventSenderProcessing() EventSenderProcessing {
	return EventSenderProcessing{
		queue: make(chan EventData, 50),
	}
}

func (e *EventSenderProcessing) StartProcessing() {
	nWorker := 437
	for i := 0; i < nWorker; i++ {
		e.wg.Add(1)
		go e.Worker()
	}
}

func (e *EventSenderProcessing) StopProcessing() {
	close(e.queue)
	e.wg.Wait()
}

func (e *EventSenderProcessing) Worker() {
	defer e.wg.Done()
	for data := range e.queue {
		var wg sync.WaitGroup
		senders := data.Senders
		for _, sender := range senders {
			wg.Add(1)
			go func(wg *sync.WaitGroup, sender EventSender) {
				defer wg.Done()
				err := sender.Send(data.Ctx, data.Event)
				if err != nil {
					log.Printf("Error: %s\n", err.Error())
				}
			}(&wg, sender)
		}
		wg.Wait()
	}
}

func (e *EventSenderProcessing) Send(ctx context.Context, senders []EventSender, event model.Event) {
	e.queue <- EventData{
		Ctx:     ctx,
		Senders: senders,
		Event:   event,
	}
}
