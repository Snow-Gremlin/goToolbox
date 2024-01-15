package changeArgs

import "github.com/Snow-Gremlin/goToolbox/collections/changeType"

type (
	changeArgsAddedImp    struct{}
	changeArgsRemovedImp  struct{}
	changeArgsReplacedImp struct{}
)

func (c changeArgsAddedImp) Type() changeType.ChangeType {
	return changeType.Added
}

func (c changeArgsRemovedImp) Type() changeType.ChangeType {
	return changeType.Removed
}

func (c changeArgsReplacedImp) Type() changeType.ChangeType {
	return changeType.Replaced
}
