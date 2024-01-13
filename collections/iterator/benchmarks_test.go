package iterator

import (
	"math/rand"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

func toSlice_SimpleStack[T any](it collections.Iterator[T]) []T {
	type node struct {
		value T
		prev  *node
	}

	count := 0
	var prev *node
	for it.Next() {
		count++
		prev = &node{
			value: it.Current(),
			prev:  prev,
		}
	}

	result := make([]T, count)
	for i := count - 1; i >= 0; i-- {
		result[i] = prev.value
		prev = prev.prev
	}
	return result
}

func toSlice_ChunkyStack[T any](it collections.Iterator[T], chunkSize int) []T {
	type node struct {
		value []T
		prev  *node
	}

	count := 0
	chunkCount := 0
	prev := &node{
		value: make([]T, chunkSize),
		prev:  nil,
	}
	for it.Next() {
		count++
		prev.value[chunkCount] = it.Current()
		chunkCount++
		if chunkCount >= chunkSize {
			prev = &node{
				value: make([]T, chunkSize),
				prev:  prev,
			}
			chunkCount = 0
		}
	}

	result := make([]T, count)
	if chunkCount > 0 {
		count -= chunkCount
		copy(result[count:], prev.value)
	}
	prev = prev.prev
	for prev != nil {
		count -= chunkSize
		copy(result[count:], prev.value)
		prev = prev.prev
	}
	return result
}

func toSlice_ChunkySlice[T any](it collections.Iterator[T], chunkSize int) []T {
	count := 0
	chunkIndex := 0
	chunkCount := 0
	data := [][]T{
		make([]T, chunkSize),
	}
	for it.Next() {
		count++
		data[chunkIndex][chunkCount] = it.Current()
		chunkCount++
		if chunkCount >= chunkSize {
			chunkIndex++
			data = append(data, make([]T, chunkSize))
			chunkCount = 0
		}
	}

	result := make([]T, count)
	if chunkCount > 0 {
		copy(result[count-chunkCount:], data[chunkIndex])
	}
	for i, j := 0, 0; i < chunkIndex; i, j = i+1, j+chunkSize {
		copy(result[j:], data[i])
	}
	return result
}

func toSlice_SliceAppend[T any](it collections.Iterator[T]) []T {
	s := []T{}
	for it.Next() {
		s = append(s, it.Current())
	}
	return s
}

func toSlice_Comparison(b *testing.B, count int) {
	src := make([]int, count)
	for i := 0; i < count; i++ {
		src[i] = int(rand.Int31())
	}

	var result1, result2, result3, result4 []int
	b.Run(`Simple Stack`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result1 = toSlice_SimpleStack(Iterate(src...))
		}
	})

	b.Run(`Chunky Stack 256`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result2 = toSlice_ChunkyStack(Iterate(src...), 256)
		}
	})

	b.Run(`Chunky Slice 256`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result3 = toSlice_ChunkySlice(Iterate(src...), 256)
		}
	})

	b.Run(`Slice Append`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result4 = toSlice_SliceAppend(Iterate(src...))
		}
	})

	checkEqual(b, result1, result2)
	checkEqual(b, result1, result3)
	checkEqual(b, result1, result4)
}

func Benchmark_Iterator_ToSlice_10(b *testing.B) {
	toSlice_Comparison(b, 10)
}

func Benchmark_Iterator_ToSlice_100(b *testing.B) {
	toSlice_Comparison(b, 100)
}

func Benchmark_Iterator_ToSlice_1000(b *testing.B) {
	toSlice_Comparison(b, 1000)
}

func Benchmark_Iterator_ToSlice_10000(b *testing.B) {
	toSlice_Comparison(b, 10000)
}

func Benchmark_Iterator_ToSlice_100000(b *testing.B) {
	toSlice_Comparison(b, 100000)
}
