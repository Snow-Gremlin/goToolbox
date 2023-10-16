package collections

// Selector is a function to convert one value type into another.
type Selector[TIn, TOut any] func(value TIn) TOut
