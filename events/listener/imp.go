package listener

import "github.com/Snow-Gremlin/goToolbox/events"

type listenerImp[T any] struct {
	obv *listenerObserver[T]
}

func (lis listenerImp[T]) Cancel() {
	lis.obv.cancel()
}

func (lis listenerImp[T]) Subscribe(event events.Event[T]) bool {
	return lis.obv.subscribe(event)
}

func (lis listenerImp[T]) Unsubscribe(event events.Event[T]) bool {
	return lis.obv.unsubscribe(event)
}
