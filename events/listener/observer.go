package listener

import (
	"slices"

	"github.com/Snow-Gremlin/goToolbox/events"
)

type listenerObserver[T any] struct {
	handle func(value T)
	events []events.Event[T]
}

func (lob *listenerObserver[T]) cancel() {
	if lob != nil {
		list := slices.Clone(lob.events)
		for _, e := range list {
			e.Remove(lob)
		}
		lob.events = []events.Event[T]{}
	}
}

func (lob *listenerObserver[T]) subscribe(event events.Event[T]) bool {
	return event.Add(lob)
}

func (lob *listenerObserver[T]) unsubscribe(event events.Event[T]) bool {
	return event.Remove(lob)
}

func (lob *listenerObserver[T]) Update(value T) {
	if lob != nil && lob.handle != nil {
		lob.handle(value)
	}
}

func (lob *listenerObserver[T]) Joined(event events.Event[T]) {
	if lob != nil && !slices.Contains(lob.events, event) {
		lob.events = append(lob.events, event)
	}
}

func (lob *listenerObserver[T]) Unjoined(event events.Event[T]) {
	if lob != nil {
		if index := slices.Index(lob.events, event); index >= 0 {
			maxIndex := len(lob.events) - 1
			copy(lob.events[index:], lob.events[index+1:])
			lob.events[maxIndex] = nil
			lob.events = lob.events[:maxIndex]
		}
	}
}
