package changeArgs

import "github.com/Snow-Gremlin/goToolbox/collections"

// NewAdded creates a new change argument for an "Added" change.
func NewAdded() collections.ChangeArgs {
	return &changeArgsAddedImp{}
}

// NewRemoved creates a new change argument for a "Removed" change.
func NewRemoved() collections.ChangeArgs {
	return &changeArgsRemovedImp{}
}

// NewReplaced creates a new change argument for a "Replaced" change.
func NewReplaced() collections.ChangeArgs {
	return &changeArgsReplacedImp{}
}
