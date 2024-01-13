package dictionary

import (
	"fmt"
	"maps"
	"math"
	"math/rand"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func removeIf_StackStore[T comparable](test simpleSet.Set[T], p collections.Predicate[T]) bool {
	type node struct {
		key  T
		next *node
	}

	var remove *node
	for key := range test {
		if p(key) {
			remove = &node{
				key:  key,
				next: remove,
			}
		}
	}
	if remove == nil {
		return false
	}

	for remove != nil {
		delete(test, remove.key)
		remove = remove.next
	}
	return true
}

func removeIf_SliceStore[T comparable](test simpleSet.Set[T], p collections.Predicate[T]) bool {
	remove := []T{}
	for key := range test {
		if p(key) {
			remove = append(remove, key)
		}
	}
	if len(remove) <= 0 {
		return false
	}
	for _, key := range remove {
		delete(test, key)
	}
	return true
}

func removeIf_NoStore[T comparable](test simpleSet.Set[T], p collections.Predicate[T]) bool {
	count := len(test)
	maps.DeleteFunc(test, func(key T, _ struct{}) bool { return p(key) })
	return len(test) != count
}

func removeIf_NoStoreOrCount[T comparable](test simpleSet.Set[T], p collections.Predicate[T]) bool {
	return test.RemoveIf(p)
}

func removeIf_Comparison(b *testing.B, count int, removePercent float64) {
	maximum := math.MaxInt16
	src := simpleSet.Cap[int](count)
	for i := 0; i < count; i++ {
		key := int(rand.Int31n(int32(maximum)))
		src.Set(key)
	}

	threshold := int(removePercent * float64(maximum))
	p := predicate.LessThan(threshold)
	var test1, test2, test3, test4 simpleSet.Set[int]

	b.Run(`Stack Store`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			test1 = maps.Clone(src)
			removeIf_StackStore(test1, p)
		}
	})

	b.Run(`Slice Store`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			test2 = maps.Clone(src)
			removeIf_SliceStore(test2, p)
		}
	})

	b.Run(`No Store`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			test3 = maps.Clone(src)
			removeIf_NoStore(test3, p)
		}
	})

	b.Run(`No Store Or Count`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			test4 = maps.Clone(src)
			removeIf_NoStoreOrCount(test4, p)
		}
	})

	keys1 := utils.SortedKeys(test1)
	keys2 := utils.SortedKeys(test2)
	keys3 := utils.SortedKeys(test3)
	keys4 := utils.SortedKeys(test4)
	check.Equal(b, keys1).Assert(keys2)
	check.Equal(b, keys1).Assert(keys3)
	check.Equal(b, keys1).Assert(keys4)
	fmt.Printf("Final size = %d (%.02f%% is about %.02f)\n", len(keys1),
		float64(len(keys1))/float64(count)*100.0, (1.0-removePercent)*100.0)
}

func Benchmark_Dictionary_RemoveIf_Remove75(b *testing.B) {
	removeIf_Comparison(b, 1000, 0.75)
}

func Benchmark_Dictionary_RemoveIf_Remove25(b *testing.B) {
	removeIf_Comparison(b, 1000, 0.25)
}
