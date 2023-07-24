package eventsender

import (
	"log"
	"time"

	"github.com/mrrizal/sample-api/model"
	"github.com/mrrizal/sample-api/utils"
)

type EventSender interface {
	Send(event model.Event) error
}

type KinesisSender struct{}

func NewKinesisSender() *KinesisSender {
	return &KinesisSender{}
}

func (k *KinesisSender) Send(event model.Event) error {
	time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	log.Println(utils.GetMessageTemplate("kinesis", &event))
	return nil
}

type SQSSender struct{}

func NewSQSSender() *SQSSender {
	return &SQSSender{}
}

func (s *SQSSender) Send(event model.Event) error {
	time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	log.Println(utils.GetMessageTemplate("sqs", &event))
	return nil
}

type APISender struct{}

func NewAPISender() *APISender {
	return &APISender{}
}

func (a *APISender) Send(event model.Event) error {
	time.Sleep(utils.RandomDuration(50, 500) * time.Millisecond)
	log.Println(utils.GetMessageTemplate("api", &event))
	return nil
}
