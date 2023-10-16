package linkedList

import (
	"math/rand"
	"testing"

	"goToolbox/collections/enumerator"
	"goToolbox/testers/check"
)

func new_ViaIteration[T any](src ...T) *linkedListImp[T] {
	return impFrom(enumerator.Enumerate(src...))
}

func new_BatchBuildNodes[T any](src ...T) *linkedListImp[T] {
	return newImp(src...)
}

func new_Comparison(b *testing.B, count int) {
	src := make([]int, count)
	for i := 0; i < count; i++ {
		src[i] = int(rand.Int31())
	}

	var list1, list2 *linkedListImp[int]
	b.Run(`Via Iteration`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			list2 = new_ViaIteration(src...)
		}
	})

	b.Run(`Batch Build Nodes`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			list2 = new_BatchBuildNodes(src...)
		}
	})

	check.Equal(b, list1).Assert(list2)
}

func Benchmark_LinkedList_NewFromSlice(b *testing.B) {
	new_Comparison(b, 10)
}

func Benchmark_LinkedList_NewFromSlice_10000(b *testing.B) {
	new_Comparison(b, 10000)
}
