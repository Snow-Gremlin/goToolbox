package testers

// Tester is a subset of testing.TB needed by these testers.
type Tester interface {
	// Error logs the values concatenated together and fails the test.
	Error(args ...any)

	// FailNow marks the test as having failed and stops its execution.
	FailNow()
}
