package collections

// Combiner is a function which merges the two given values into one value.
type Combiner[TFirst, TSecond, TOut any] func(first TFirst, second TSecond) TOut
