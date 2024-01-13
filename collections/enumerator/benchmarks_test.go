package enumerator

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type buffered_Handle[T any] func(e collections.Enumerator[T]) collections.Enumerator[T]

func buffered_SimpleStack[T any](e collections.Enumerator[T]) collections.Enumerator[T] {
	type node struct {
		value T
		next  *node
	}
	var first, last *node
	var source collections.Iterator[T]
	return New(func() collections.Iterator[T] {
		if source == nil {
			source = e.Iterate()
		}
		n := first
		return iterator.New(func() (T, bool) {
			if n != nil {
				value := n.value
				n = n.next
				return value, true
			}
			if source.Next() {
				value := source.Current()
				tail := &node{
					value: value,
					next:  nil,
				}
				if last != nil {
					last.next = tail
				} else {
					first = tail
				}
				last = tail
				return value, true
			}
			return utils.Zero[T](), false
		})
	})
}

func buffered_GrowingSlice[T any](e collections.Enumerator[T]) collections.Enumerator[T] {
	return e.Buffered()
}

func buffered_FullSlice[T any](e collections.Enumerator[T]) collections.Enumerator[T] {
	first := true
	var source []T
	var count int
	return New(func() collections.Iterator[T] {
		index := 0
		return iterator.New(func() (T, bool) {
			if first {
				first = false
				source = e.ToSlice()
				count = len(source)
			}
			if index < count {
				value := source[index]
				index++
				return value, true
			}
			return utils.Zero[T](), false
		})
	})
}

func buffered_Trial(b *testing.B, count int, frac float64, src []int, buffered buffered_Handle[int]) {
	partSize := int(float64(count) * frac)
	e0 := Enumerate(src...)
	halfSrc := e0.Take(partSize).ToSlice()
	fracStr := fmt.Sprintf(`%d%%`, int(frac*100.0))

	data := make([]int, partSize)
	b.Run(`Load `+fracStr, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffered(e0).CopyToSlice(data)
		}
	})
	checkEqual(b, halfSrc, data)

	e1 := buffered(e0)
	e1.CopyToSlice(data) // preload
	b.Run(`Read `+fracStr, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e1.CopyToSlice(data)
		}
	})
	checkEqual(b, halfSrc, data)
}

func buffered_Comparison(b *testing.B, count int) {
	src := make([]int, count)
	for i := 0; i < count; i++ {
		src[i] = int(rand.Int31())
	}

	b.Run(`Simple Stack`, func(b *testing.B) {
		buffered_Trial(b, count, 0.1, src, buffered_SimpleStack)
		buffered_Trial(b, count, 1.0, src, buffered_SimpleStack)
	})

	b.Run(`Growing Slice Stack`, func(b *testing.B) {
		buffered_Trial(b, count, 0.1, src, buffered_GrowingSlice)
		buffered_Trial(b, count, 1.0, src, buffered_GrowingSlice)
	})

	b.Run(`Full Slice`, func(b *testing.B) {
		buffered_Trial(b, count, 0.1, src, buffered_FullSlice)
		buffered_Trial(b, count, 1.0, src, buffered_FullSlice)
	})
}

func Benchmark_Enumerator_Buffered_100000(b *testing.B) {
	buffered_Comparison(b, 100000)
}
