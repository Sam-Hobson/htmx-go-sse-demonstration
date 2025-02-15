package routes

import (
	"context"
	"fmt"
	"salad2/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

type EventData struct {
	Event string
	Data  []string
}

type EventDataProvider <-chan *EventData

func AddSSEHeaders(c echo.Context) {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
}

func WriteSSEEvent(c echo.Context, eventData *EventData) error {
	if _, err := fmt.Fprintf(c.Response(), "event: %s\n", eventData.Event); err != nil {
		return err
	}

	for _, line := range eventData.Data {
		if _, err := fmt.Fprintf(c.Response(), "data: %s\n", line); err != nil {
			return err
		}
	}

	if len(eventData.Data) == 0 {
		if _, err := fmt.Fprintf(c.Response(), "data:\n"); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(c.Response(), "\n"); err != nil {
		return err
	}

	c.Response().Flush()
	return nil
}

func HandleSSEConnection(c echo.Context, providers ...EventDataProvider) error {
	AddSSEHeaders(c)

	errCh := make(chan error, len(providers))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, provider := range providers {
		go func(p EventDataProvider) {
			for {
				select {
				case eventData, ok := <-p:
					// Provider is closed
					if !ok {
						errCh <- nil
						return
					}
					if err := WriteSSEEvent(c, eventData); err != nil {
						errCh <- err
						return
					}
				// When this function returns, either the connection is dead or
				// an error occurred in a provider. Either way, kill all providers.
				case <-ctx.Done():
					return
				}
			}
		}(provider)
	}

	select {
	case <-c.Request().Context().Done():
	case err := <-errCh:
		return err
	}

	return nil
}

func EventDataProviderFromTemplate[T any](ctx context.Context, source <-chan T, event string, templateName string) EventDataProvider {
	eventDataProvider := make(chan *EventData)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case obj := <-source:
				data, err := utils.RenderTemplateToString(templateName, obj)
				if err != nil {
					close(eventDataProvider)
					return
				}

				eventDataProvider <- &EventData{
					Event: event,
					Data:  strings.Split(data, "\n"),
				}
			}
		}
	}()

	return eventDataProvider
}

func EventDataProviderEventOnly[T any](ctx context.Context, source <-chan T, event string) EventDataProvider {
	eventDataProvider := make(chan *EventData)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-source:
				eventDataProvider <- &EventData{
					Event: event,
				}
			}
		}
	}()

	return eventDataProvider
}

func EventDataProviderCustomData[T any](ctx context.Context, source <-chan T, event string, custom func(T) []string) EventDataProvider {
	eventDataProvider := make(chan *EventData)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case obj := <-source:
				eventDataProvider <- &EventData{
					Event: event,
					Data:  custom(obj),
				}
			}
		}
	}()

	return eventDataProvider
}
