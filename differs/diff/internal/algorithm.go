package internal

// Algorithm is the interface for a diff algorithm.
type Algorithm interface {

	// NoResizeNeeded determines if the diff algorithm can handle a container with
	// the amount of data inside of the given container. If this returns false a
	// larger matrix, cache, vector, or whatever would be created to perform the diff.
	NoResizeNeeded(cont Container) bool

	// Diff performs the algorithm on the given container
	// and writes the results to the collector.
	Diff(cont Container, col Collector)
}
