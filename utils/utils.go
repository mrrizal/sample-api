package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/mrrizal/sample-api/model"
)

func GetMessageTemplate(senderName string, event *model.Event) string {
	eventTime := event.Time.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s sender: data has been sent, %s by %s", eventTime, senderName, event.Name, event.Actor)
}

func RandomDuration(min, max time.Duration) time.Duration {
	n := rand.Int63n(int64(max - min))
	return min + time.Duration(n)
}

func GenerateRandomEvent() model.Event {
	var event model.Event
	faker.FakeData(&event)
	event.Time = time.Now().Add(-time.Duration(rand.Intn(1440-1)+1) * time.Minute)
	return event
}
