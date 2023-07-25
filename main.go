package main

import (
	"context"
	"log"
	"os"
	"sync"

	eventsender "github.com/mrrizal/sample-api/event_sender"
	"github.com/mrrizal/sample-api/model"
	"github.com/mrrizal/sample-api/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/gofiber/fiber/v2"
)

type EventService struct {
	senders []eventsender.EventSender
}

func NewEventService(senders []eventsender.EventSender) *EventService {
	return &EventService{senders: senders}
}

func (e *EventService) Send(ctx context.Context, event model.Event) {
	ctx, span := utils.StartTracerSpan(ctx, "eventService/Send")
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

type EventHandler struct {
	svc *EventService
}

func NewEventHandler(svc *EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}

func (e *EventHandler) Like(c *fiber.Ctx) error {
	ctx, span := utils.StartTracerSpan(c.Context(), "EventHandler/Like")
	defer span.End()

	event := utils.GenerateRandomEvent()
	event.Name = "like"

	e.svc.Send(ctx, event)
	return c.SendStatus(fiber.StatusNoContent)
}

func (e *EventHandler) Unlike(c *fiber.Ctx) error {
	ctx, span := utils.StartTracerSpan(c.Context(), "EventHandler/Unlike")
	defer span.End()

	event := utils.GenerateRandomEvent()
	event.Name = "unlike"

	e.svc.Send(ctx, event)
	return c.SendStatus(fiber.StatusNoContent)
}

func iniTrace() {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	serviceName := os.Getenv("SERVICE_NAME")

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerEndpoint),
		),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(jaegerExporter),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
		),
	)
}

func main() {
	iniTrace()

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

	err := app.Listen(":80")
	if err != nil {
		panic("Error starting the server: " + err.Error())
	}
}
