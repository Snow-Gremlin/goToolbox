package collections

import "github.com/Snow-Gremlin/goToolbox/collections/changeType"

// ChangeArgs is the value returned by an OnChange event.
type ChangeArgs interface {
	// Type gets the type of change that caused the OnChange event to invoke.
	Type() changeType.ChangeType
}
