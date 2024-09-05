package event

import (
	"slices"

	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

type eventImp[T any] struct {
	obs []events.Observer[T]
}

func (e *eventImp[T]) Add(observer events.Observer[T]) bool {
	if e == nil || liteUtils.IsNil(observer) || slices.Contains(e.obs, observer) {
		return false
	}
	e.obs = append(e.obs, observer)
	if jobs, ok := observer.(events.Joinable[T]); ok {
		jobs.Joined(e)
	}
	return true
}

func (e *eventImp[T]) Remove(observer events.Observer[T]) bool {
	if e == nil || liteUtils.IsNil(observer) {
		return false
	}
	index := slices.Index(e.obs, observer)
	if index < 0 {
		return false
	}

	maxIndex := len(e.obs) - 1
	copy(e.obs[index:], e.obs[index+1:])
	e.obs[maxIndex] = nil
	e.obs = e.obs[:maxIndex]
	if jobs, ok := observer.(events.Unjoinable[T]); ok {
		jobs.Unjoined(e)
	}
	return true
}

func (e *eventImp[T]) Clear() {
	if e != nil {
		for _, ob := range e.obs {
			if jobs, ok := ob.(events.Unjoinable[T]); ok {
				jobs.Unjoined(e)
			}
		}
		e.obs = nil
	}
}

func (e *eventImp[T]) Invoke(value T) {
	if e != nil {
		for _, ob := range e.obs {
			ob.Update(value)
		}
	}
}
