package collections

// Window is a function which handles a sliding window of values
// and returns a result from that window.
type Window[TIn, TOut any] func(values []TIn) TOut
