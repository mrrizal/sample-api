package eventsender

import (
	"context"
	"log"
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
