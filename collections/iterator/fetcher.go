package iterator

// Fetcher is the source of values which can be used in iterators.
//
// Each time it is called it can return a different value.
// Returns true if a new value was fetched and false if there are no new values.
// Once false is returned it is expected to always return false from then on.
type Fetcher[T any] func() (T, bool)
