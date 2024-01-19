package collections

import "github.com/Snow-Gremlin/goToolbox/events"

// OnChanger is an object which can emit a change event.
type OnChanger interface {
	// OnChange gets the event that is invoked on change.
	OnChange() events.Event[ChangeArgs]
}
