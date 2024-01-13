package container

import (
	"github.com/Snow-Gremlin/goToolbox/differs"
	"github.com/Snow-Gremlin/goToolbox/differs/diff/internal"
)

// New creates a new container for full data.
func New(data differs.Data) internal.Container {
	return newSub(data,
		0, data.ACount(),
		0, data.BCount(),
		false)
}
