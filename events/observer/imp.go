package observer

type observerImp[T any] struct {
	handle func(T)
}

func (e *observerImp[T]) Update(value T) {
	if e.handle != nil {
		e.handle(value)
	}
}
