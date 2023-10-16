package diff

import (
	"sync"

	"goToolbox/differs"
	"goToolbox/differs/data"
	"goToolbox/differs/diff/internal"
	"goToolbox/differs/diff/internal/collector"
	"goToolbox/differs/diff/internal/container"
)

func wrap(alg internal.Algorithm) differs.Diff {
	return &differImp{
		alg:  alg,
		lock: &sync.Mutex{},
	}
}

type differImp struct {
	alg  internal.Algorithm
	lock *sync.Mutex
}

func (i *differImp) Diff(data differs.Data) differs.Result {
	col := collector.New(data.ACount(), data.BCount())
	cont := container.New(data)
	cont, before, after := cont.Reduce()
	col.InsertEqual(after)
	if !cont.EndCase(col) {
		i.lock.Lock()
		defer i.lock.Unlock()
		i.alg.Diff(cont, col)
	}
	col.InsertEqual(before)
	return col.Finish()
}

func (i *differImp) PlusMinus(a, b []string) []string {
	return PlusMinus(i.Diff(data.Strings(a, b)), a, b)
}

func (i *differImp) Merge(a, b []string) []string {
	return Merge(i.Diff(data.Strings(a, b)), a, b)
}
