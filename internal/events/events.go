package events

import (
	"context"
	"sync"
	"time"
)

type EventQueue[K any, T any] interface {
	Queue(data T)
	Subscribe(key K, chanSize int) <-chan *T
	Unsubscribe(key K)
	Close()
}

func addDataToAllChannels[T any](timeout time.Duration, m *sync.Map, data T) {
	m.Range(func(key, value any) bool {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			addDataToChannel(ctx, value.(chan T), data)
		}()
		return true
	})
}

func addDataToChannel[T any](ctx context.Context, ch chan T, data T) bool {
	for {
		select {
		case ch <- data:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

func closeAll[T any](m *sync.Map) {
	m.Range(func(key, value any) bool {
		close(value.(chan T))
		return true
	})
}
