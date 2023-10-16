package collections

// Reducer is a function which merges the given value with the prior
// value to create a new value to use in next reduce call.
type Reducer[TIn, TOut any] func(value TIn, prior TOut) TOut
