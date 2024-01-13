package iterator

import (
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_Iterator_NilFetcher(t *testing.T) {
	it := New[int](nil)
	checkZero(t, it.Current())
	checkEqual(t, false, it.Next())
	checkZero(t, it.Current())
}

func Test_Iterator_New(t *testing.T) {
	count := 4
	it := New(func() (int, bool) {
		if count > 0 {
			count--
			return count, true
		}
		return count, false
	})

	checkZero(t, it.Current())

	checkEqual(t, true, it.Next())
	checkEqual(t, 3, it.Current())
	checkEqual(t, 3, count)

	checkEqual(t, true, it.Next())
	checkEqual(t, 2, it.Current())
	checkEqual(t, 2, count)

	checkEqual(t, true, it.Next())
	checkEqual(t, 1, it.Current())
	checkEqual(t, 1, count)

	checkEqual(t, true, it.Next())
	checkZero(t, it.Current())
	checkZero(t, count)

	checkEqual(t, false, it.Next())
	checkZero(t, it.Current())
	checkZero(t, count)

	checkEqual(t, false, it.Next())
	checkZero(t, it.Current())
	checkZero(t, count)
}

func Test_Iterator_Iterate(t *testing.T) {
	it := Iterate(1, 1, 2, 3, 5, 8)
	checkIt(t, it, 1, 1, 2, 3, 5, 8)

	it = Iterate[int]()
	checkIt(t, it)
}

func Test_Iterator_Range(t *testing.T) {
	it := Range(9, 4)
	checkIt(t, it, 9, 10, 11, 12)
}

func Test_Iterator_Stride(t *testing.T) {
	it := Stride(9.0, -1.5, 9)
	checkIt(t, it, 9.0, 7.5, 6.0, 4.5, 3.0, 1.5, 0.0, -1.5, -3.0)
}

func Test_Iterator_Repeat(t *testing.T) {
	it := Repeat(`cat`, 4)
	checkIt(t, it, `cat`, `cat`, `cat`, `cat`)
}

func Test_Iterator_Where(t *testing.T) {
	it := Where(Iterate(1, 9, 2, 3, 5, 1, 8), predicate.LessEq(3))
	checkIt(t, it, 1, 2, 3, 1)
}

func Test_Iterator_WhereNot(t *testing.T) {
	it := WhereNot(Iterate(1, 9, 2, 3, 5, 1, 8), predicate.LessEq(3))
	checkIt(t, it, 9, 5, 8)
}

func Test_Iterator_ToSlice(t *testing.T) {
	it := Iterate(1, 1, 2, 3, 5, 8)
	checkEqual(t, []int{1, 1, 2, 3, 5, 8}, ToSlice(it))

	it = Iterate[int]()
	checkEqual(t, []int{}, ToSlice(it))
}

func Test_Iterator_CopyToSlice(t *testing.T) {
	it := Iterate(1, 1, 2, 3, 5, 8)

	s := make([]int, 3)
	CopyToSlice(it, s)
	checkEqual(t, []int{1, 1, 2}, s)

	s = make([]int, 5)
	CopyToSlice(it, s)
	checkEqual(t, []int{3, 5, 8, 0, 0}, s)

	s = make([]int, 2)
	CopyToSlice(it, s)
	checkEqual(t, []int{0, 0}, s)

	CopyToSlice(it, nil) // Doesn't fail
}

func Test_Iterator_Foreach(t *testing.T) {
	product := 1
	f := func(i int) {
		product *= i
	}

	it := Iterate(1, 1, 2, 3, 5, 8)
	Foreach(it, f)
	checkEqual(t, 240, product)
}

func Test_Iterator_DoUntilError(t *testing.T) {
	results := []int{}
	add := func(i int) error {
		if i <= 0 {
			return terror.New(`Must be positive.`).With(`value`, i)
		}
		results = append(results, int(1000.0/float64(i)))
		return nil
	}

	it := Iterate(1, 2, 3, 4)
	err := DoUntilError(it, add)
	checkEqual(t, nil, err)
	checkEqual(t, []int{1000, 500, 333, 250}, results)

	it = Iterate(5, 6, -2, 7, 8)
	err = DoUntilError(it, add)
	checkEqual(t, `Must be positive. {value: -2}`, err.Error())
	checkEqual(t, []int{1000, 500, 333, 250, 200, 166}, results)
	checkIt(t, it, 7, 8) // Remainder
}

func Test_Iterator_DoUntilNotZero(t *testing.T) {
	bloop := func(i int) string {
		if i == 3 {
			return `Three`
		}
		return ``
	}

	it := Iterate(1, 2, 3, 4, 5, 6)
	val := DoUntilNotZero(it, bloop)
	checkEqual(t, `Three`, val)
	checkIt(t, it, 4, 5, 6) // Remainder

	it = Iterate(7, 8, 9)
	val = DoUntilNotZero(it, bloop)
	checkLength(t, 0, val)
	checkIt(t, it) // Remainder
}

func Test_Iterator_Any(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6)
	checkEqual(t, true, Any(it, predicate.LessEq(3)))
	checkIt(t, it, 7, 9, 6) // Remainder

	it = Iterate(7, 8, 6, 9, 7, 9, 6)
	checkEqual(t, false, Any(it, predicate.LessEq(3)))
	checkIt(t, it) // Remainder
}

func Test_Iterator_All(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6)
	checkEqual(t, false, All(it, predicate.GreaterThan(3)))
	checkIt(t, it, 7, 9, 6) // Remainder

	it = Iterate(7, 8, 6, 9, 7, 9, 6)
	checkEqual(t, true, All(it, predicate.GreaterThan(3)))
	checkIt(t, it) // Remainder
}

func Test_Iterator_StepsUntil(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6)
	checkEqual(t, 4, StepsUntil(it, predicate.LessEq(3)))
	checkIt(t, it, 7, 9, 6) // Remainder

	it = Iterate(7, 8, 6, 9, 7, 9, 6)
	checkEqual(t, -1, StepsUntil(it, predicate.LessEq(3)))
	checkIt(t, it) // Remainder
}

func Test_Iterator_StartsWith(t *testing.T) {
	it1 := Iterate(1, 2, 3, 4, 5)
	it2 := Iterate(1, 2, 3, 4, 5)
	checkEqual(t, true, StartsWith(it1, it2))
	checkIt(t, it1) // Remainder
	checkIt(t, it2) // Remainder

	it1 = Iterate(1, 2, 3, 4, 5)
	it2 = Iterate(1, 2, 3)
	checkEqual(t, true, StartsWith(it1, it2))
	checkIt(t, it1, 5) // Remainder
	checkIt(t, it2)    // Remainder

	it1 = Iterate(1, 2, 3)
	it2 = Iterate(1, 2, 3, 4, 5)
	checkEqual(t, false, StartsWith(it1, it2))
	checkIt(t, it1)    // Remainder
	checkIt(t, it2, 5) // Remainder

	it1 = Iterate(1, 8, 3, 4, 5)
	it2 = Iterate(1, 2, 3)
	checkEqual(t, false, StartsWith(it1, it2))
	checkIt(t, it1, 3, 4, 5) // Remainder
	checkIt(t, it2, 3)       // Remainder
}

func Test_Iterator_Equal(t *testing.T) {
	it1 := Iterate(1, 2, 3, 4, 5)
	it2 := Iterate(1, 2, 3, 4, 5)
	checkEqual(t, true, Equal(it1, it2))
	checkIt(t, it1) // Remainder
	checkIt(t, it2) // Remainder

	it1 = Iterate(1, 2, 3, 4, 5)
	it2 = Iterate(1, 2, 3)
	checkEqual(t, false, Equal(it1, it2))
	checkIt(t, it1, 5) // Remainder
	checkIt(t, it2)    // Remainder

	it1 = Iterate(1, 2, 3)
	it2 = Iterate(1, 2, 3, 4, 5)
	checkEqual(t, false, Equal(it1, it2))
	checkIt(t, it1)    // Remainder
	checkIt(t, it2, 5) // Remainder

	it1 = Iterate(1, 8, 3, 4, 5)
	it2 = Iterate(1, 2, 3)
	checkEqual(t, false, Equal(it1, it2))
	checkIt(t, it1, 3, 4, 5) // Remainder
	checkIt(t, it2, 3)       // Remainder
}

func Test_Iterator_Empty(t *testing.T) {
	it := Iterate(1, 2, 3, 4)
	checkEqual(t, false, Empty(it))
	checkIt(t, it, 2, 3, 4) // Remainder

	it = Iterate[int]()
	checkEqual(t, true, Empty(it))
	checkIt(t, it) // Remainder

	it = Iterate(1, 2, 3, 4)
	checkEqual(t, false, Empty(it)) // Read 1
	checkEqual(t, false, Empty(it)) // Read 2
	checkEqual(t, false, Empty(it)) // Read 3
	checkEqual(t, false, Empty(it)) // Read 4
	checkEqual(t, true, Empty(it))
}

func Test_Iterator_Count(t *testing.T) {
	it := Iterate(1, 2, 3, 4)
	checkEqual(t, 4, Count(it))

	it = Iterate[int]()
	checkZero(t, Count(it))
}

func Test_Iterator_AtLeast(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5, 6, 7, 8)
	checkEqual(t, true, AtLeast(it, 3))
	checkIt(t, it, 4, 5, 6, 7, 8) // Remainder

	it = Iterate(1, 2, 3)
	checkEqual(t, true, AtLeast(it, 3))
	checkIt(t, it) // Remainder

	it = Iterate(1, 2)
	checkEqual(t, false, AtLeast(it, 3))
	checkIt(t, it) // Remainder
}

func Test_Iterator_AtMost(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5, 6, 7, 8)
	checkEqual(t, false, AtMost(it, 3))
	checkIt(t, it, 5, 6, 7, 8) // Remainder

	it = Iterate(1, 2, 3)
	checkEqual(t, true, AtMost(it, 3))
	checkIt(t, it) // Remainder

	it = Iterate(1, 2)
	checkEqual(t, true, AtMost(it, 3))
	checkIt(t, it) // Remainder
}

func Test_Iterator_First(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	value, has := First(it)
	checkEqual(t, true, has)
	checkEqual(t, 1, value)
	value, has = First(it)
	checkEqual(t, true, has)
	checkEqual(t, 2, value)
	checkIt(t, it, 3, 4, 5) // Remainder

	it = Iterate[int]()
	value, has = First(it)
	checkEqual(t, false, has)
	checkZero(t, value)
}

func Test_Iterator_Last(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	value, has := Last(it)
	checkEqual(t, true, has)
	checkEqual(t, 5, value)

	it = Iterate[int]()
	value, has = Last(it)
	checkEqual(t, false, has)
	checkZero(t, value)
}

func Test_Iterator_Single(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	value, has := Single(it)
	checkEqual(t, false, has)
	checkZero(t, value)

	it = Iterate(3)
	value, has = Single(it)
	checkEqual(t, true, has)
	checkEqual(t, 3, value)

	it = Iterate[int]()
	value, has = Single(it)
	checkEqual(t, false, has)
	checkZero(t, value)
}

func Test_Iterator_Skip(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	checkIt(t, Skip(it, 2), 3, 4, 5)

	it = Iterate(1, 2)
	checkIt(t, Skip(it, 2))

	it = Iterate[int]()
	checkIt(t, Skip(it, 2))

	count := 0
	it = Skip(watcher(&count, Iterate(11, 22, 33, 44, 55)), 2)
	checkZero(t, count) // Skip doesn't occur until first Next
	checkEqual(t, true, it.Next())
	checkEqual(t, 3, count)
	checkEqual(t, 33, it.Current())
}

func Test_Iterator_SkipWhile(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	checkIt(t, SkipWhile(it, predicate.LessThan(3)), 3, 4, 5)

	it = Iterate(1, 2)
	checkIt(t, SkipWhile(it, predicate.LessThan(8)))

	it = Iterate(1, 2)
	checkIt(t, SkipWhile(it, nil), 1, 2)
}

func Test_Iterator_Take(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5, 6, 7, 8)
	checkIt(t, Take(it, 3), 1, 2, 3)
	checkIt(t, Take(it, 3), 4, 5, 6)
	checkIt(t, Take(it, 3), 7, 8)
	checkIt(t, Take(it, 3))

	it = Iterate(1, 2, 3, 4, 5)
	checkIt(t, Take(it, 4), 1, 2, 3, 4)
}

func Test_Iterator_TakeWhile(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5, 6, 7, 8)
	checkIt(t, TakeWhile(it, predicate.LessThan(4)), 1, 2, 3)
	checkIt(t, it, 5, 6, 7, 8) // Remainder

	it = Iterate(1, 2, 3, 4, 5)
	checkIt(t, TakeWhile(it, predicate.GreaterThan(0)), 1, 2, 3, 4, 5)
}

func Test_Iterator_Replace(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	checkIt(t, Replace(it, func(v int) int { return v % 4 }),
		1, 2, 3, 0, 1, 2, 3, 0, 1, 2)
}

func Test_Iterator_Reverse(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	checkIt(t, Reverse(it), 5, 4, 3, 2, 1)

	it = Iterate[int]()
	checkIt(t, Reverse(it))

	count := 0
	it = Reverse(watcher(&count, Iterate(11, 22, 33, 44, 55)))
	checkZero(t, count) // Reverse doesn't start reverse until first Next
	checkEqual(t, true, it.Next())
	checkEqual(t, 6, count) // After first Next all values have been read.
	checkEqual(t, 55, it.Current())
}

func Test_Iterator_Append(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	checkIt(t, Append(it, []int{6, 7, 8}), 1, 2, 3, 4, 5, 6, 7, 8)

	it = Iterate[int]()
	checkIt(t, Append(it, []int{6, 7, 8}), 6, 7, 8)

	it = Iterate(1, 2, 3, 4, 5)
	checkIt(t, Append(it, []int{}), 1, 2, 3, 4, 5)
}

func Test_Iterator_Concat(t *testing.T) {
	it1 := Iterate(1, 2, 3)
	it2 := Iterate(4, 5, 6)
	it3 := Iterate[int]()
	it4 := Iterate(7, 8, 9)
	it5 := Concat([]collections.Iterator[int]{it1, it2, it3, it4})
	checkIt(t, it5, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	it5 = Concat([]collections.Iterator[int]{})
	checkIt(t, it5)
}

func Test_Iterator_Select(t *testing.T) {
	it := Select(Range(1, 5), func(v int) float64 { return 10.0 - float64(v)*2.0 })
	checkIt(t, it, 8.0, 6.0, 4.0, 2.0, 0.0)
}

func Test_Iterator_OfType(t *testing.T) {
	it1 := OfType[any, int](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it1, 5, 42)

	it2 := OfType[any, float64](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it2, 1.0, 3.1)

	it3 := OfType[any, string](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it3, `cat`, `dog`)
}

func Test_Iterator_Cast(t *testing.T) {
	it1 := Cast[any, int](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it1, 0, 1, 5, 0, 42, 3)

	it2 := Cast[any, float64](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it2, 0.0, 1.0, 5.0, 0.0, 42.0, 3.1)

	it3 := Cast[any, string](Iterate[any](`cat`, 1.0, 5, `dog`, 42, 3.1))
	checkIt(t, it3, `cat`, ``, "\x05", `dog`, "\x2A", ``)
}

func Test_Iterator_Expand(t *testing.T) {
	it1 := Iterate(`a`, ``, `cat`, `dog`, `apple`)
	it2 := Expand(it1, func(s string) collections.Iterable[rune] {
		return func() collections.Iterator[rune] {
			return Iterate[rune]([]rune(s)...)
		}
	})
	checkIt[rune](t, it2, 'a', 'c', 'a', 't', 'd', 'o', 'g', 'a', 'p', 'p', 'l', 'e')

	it1 = Iterate(`one`)
	it2 = Expand(it1, func(s string) collections.Iterable[rune] {
		return nil
	})
	checkIt(t, it2)
}

func Test_Iterator_Reduce(t *testing.T) {
	it := Iterate(`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`, `dog`, ``, `is`, `cute`, `and`, `dumb`)
	sizes := Reduce(it, []int{}, func(s string, sizes []int) []int {
		w := len(s)
		for len(sizes) <= w {
			sizes = append(sizes, 0)
		}
		sizes[w]++
		return sizes
	})
	checkEqual(t, []int{1, 0, 2, 7, 3, 0, 1}, sizes)
}

func Test_Iterator_Merge(t *testing.T) {
	it := Iterate(`dog`, `cat`, `horse`, `taco`)
	result := Merge(it, func(value, prior string) string {
		return value + `|` + prior
	})
	checkEqual(t, `taco|horse|cat|dog`, result)
}

func Test_Iterator_SlidingWindow(t *testing.T) {
	en := func() collections.Iterator[string] {
		return Iterate(`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`, `dog`, ``, `is`, `cute`, `and`, `dumb`)
	}

	windowHandle := func(size int) collections.Window[string, string] {
		return func(window []string) string {
			checkLength(t, size, window)
			return strings.Join(window, `|`)
		}
	}

	doNotWindowHandle := func() collections.Window[string, string] {
		return func(window []string) string {
			t.Error(`window method should not be called if the window is size zero`)
			return strings.Join(window, `|`)
		}
	}

	checkIt(t, SlidingWindow(en(), 4, 1, windowHandle(4)),
		`the|cat|is|sneaky`, `cat|is|sneaky|and`, `is|sneaky|and|evil`,
		`sneaky|and|evil|but`, `and|evil|but|the`, `evil|but|the|dog`, `but|the|dog|`,
		`the|dog||is`, `dog||is|cute`, `|is|cute|and`, `is|cute|and|dumb`)

	checkIt(t, SlidingWindow(en(), 4, 2, windowHandle(4)),
		`the|cat|is|sneaky`, `is|sneaky|and|evil`, `and|evil|but|the`,
		`but|the|dog|`, `dog||is|cute`, `is|cute|and|dumb`)

	checkIt(t, SlidingWindow(en(), 4, 3, windowHandle(4)),
		`the|cat|is|sneaky`, `sneaky|and|evil|but`, `but|the|dog|`, `|is|cute|and`)

	checkIt(t, SlidingWindow(en(), 4, 4, windowHandle(4)),
		`the|cat|is|sneaky`, `and|evil|but|the`, `dog||is|cute`)

	checkIt(t, SlidingWindow(Iterate(`the`, `cat`), 4, 1, doNotWindowHandle()))

	checkPanic(t, `the given window size must be greater than zero {size: 0}`,
		func() { SlidingWindow(en(), 0, 1, doNotWindowHandle()) })
	checkPanic(t, `the given window stride must be greater than zero and less than or equal to the size {size: 4, stride: 0}`,
		func() { SlidingWindow(en(), 4, 0, doNotWindowHandle()) })
	checkPanic(t, `the given window stride must be greater than zero and less than or equal to the size {size: 4, stride: 5}`,
		func() { SlidingWindow(en(), 4, 5, doNotWindowHandle()) })
}

func Test_Iterator_Chunk(t *testing.T) {
	en := func() collections.Iterator[string] {
		return Iterate(`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`, `dog`, ``, `is`, `cute`, `and`, `dumb`)
	}

	checkIt(t, Chunk(en(), 1),
		[]string{`the`}, []string{`cat`}, []string{`is`}, []string{`sneaky`}, []string{`and`}, []string{`evil`}, []string{`but`},
		[]string{`the`}, []string{`dog`}, []string{``}, []string{`is`}, []string{`cute`}, []string{`and`}, []string{`dumb`})

	checkIt(t, Chunk(en(), 2),
		[]string{`the`, `cat`}, []string{`is`, `sneaky`}, []string{`and`, `evil`}, []string{`but`, `the`},
		[]string{`dog`, ``}, []string{`is`, `cute`}, []string{`and`, `dumb`})

	checkIt(t, Chunk(en(), 3),
		[]string{`the`, `cat`, `is`}, []string{`sneaky`, `and`, `evil`}, []string{`but`, `the`, `dog`},
		[]string{``, `is`, `cute`}, []string{`and`, `dumb`})

	checkIt(t, Chunk(en(), 8),
		[]string{`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`},
		[]string{`dog`, ``, `is`, `cute`, `and`, `dumb`})

	checkIt(t, Chunk(en(), 13),
		[]string{`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`, `dog`, ``, `is`, `cute`, `and`}, []string{`dumb`})

	checkIt(t, Chunk(en(), 20),
		[]string{`the`, `cat`, `is`, `sneaky`, `and`, `evil`, `but`, `the`, `dog`, ``, `is`, `cute`, `and`, `dumb`})

	checkPanic(t, `the given chunk size must be greater than zero {size: 0}`,
		func() { Chunk(en(), 0) })
}

func Test_Iterator_Sum(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6)
	result, count := Sum(it)
	checkEqual(t, 53, result)
	checkEqual(t, 8, count)
}

func Test_Iterator_IsUnique(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6, 1, 1, 8, 9)
	checkEqual(t, false, IsUnique(it))

	it = Iterate(7, 8, 6, 9, 1)
	checkEqual(t, true, IsUnique(it))
}

func Test_Iterator_Unique(t *testing.T) {
	it := Iterate(7, 8, 6, 9, 1, 7, 9, 6, 1, 1, 8, 9)
	checkIt(t, Unique(it), 7, 8, 6, 9, 1)
}

func Test_Iterator_Intersection(t *testing.T) {
	count1, count2 := 0, 0
	it1 := watcher(&count1, Iterate(7, 8, 6, 1, 1, 8, 9, 4))
	it2 := watcher(&count2, Iterate(7, 5, 6, 9, 3, 1, 7, 9, 5, 9))
	it3 := Intersection(it1, it2)

	checkZero(t, count1)
	checkZero(t, count2)
	checkEqual(t, true, it3.Next())
	checkEqual(t, 7, it3.Current())
	checkEqual(t, 1, count1) // iterator 1 only read as much as needed
	checkEqual(t, 1, count2)

	checkIt(t, it3, 6, 9, 1, 7, 9, 9) // order and duplicates from it2
}

func Test_Iterator_Subtract(t *testing.T) {
	count1, count2 := 0, 0
	it1 := watcher(&count1, Iterate(7, 8, 6, 1, 1, 8, 9, 4))
	it2 := watcher(&count2, Iterate(7, 5, 6, 9, 3, 1, 7, 9, 5, 9))
	it3 := Subtract(it1, it2)

	checkZero(t, count1)
	checkZero(t, count2)
	checkEqual(t, true, it3.Next())
	checkEqual(t, 5, it3.Current())
	checkEqual(t, 9, count1) // iterator 1 reads whole thing trying to find another 5
	checkEqual(t, 2, count2)

	checkIt(t, it3, 3, 5) // order and duplicates from it2
}

func Test_Iterator_Zip(t *testing.T) {
	it1 := Iterate(`a`, `b`, `c`, `d`, `ex`, `cat `)
	it2 := Iterate(1, 3, 1, 3, 1, 2)
	it3 := Zip(it1, it2, strings.Repeat)
	checkIt(t, it3, `a`, `bbb`, `c`, `ddd`, `ex`, `cat cat `)
}

func Test_Iterator_Interweave(t *testing.T) {
	it1 := Repeat(1, 9)
	it2 := Repeat(2, 7)
	it3 := Repeat(3, 5)
	it4 := Repeat(4, 3)
	it5 := Interweave([]collections.Iterator[int]{it1, it2, it3, it4})
	checkIt(t, it5,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3,
		1, 2, 3,
		1, 2,
		1, 2,
		1,
		1)
}

func Test_Iterator_SortInterweave(t *testing.T) {
	it1 := Iterate(1, 6, 2, 5, 3, 4)
	it2 := Iterate(4, 3, 5, 2, 6, 1)
	it3 := SortInterweave(it1, it2, utils.OrderedComparer[int]())
	checkIt(t, it3, 1, 4, 3, 5, 2, 6, 2, 5, 3, 4, 6, 1)
}

func Test_Iterator_Indexed(t *testing.T) {
	it := Indexed(Repeat(`a`, 5))
	checkIt(t, it,
		tuple2.New(0, `a`),
		tuple2.New(1, `a`),
		tuple2.New(2, `a`),
		tuple2.New(3, `a`),
		tuple2.New(4, `a`))
}

func Test_Iterator_Sort(t *testing.T) {
	it1 := Iterate(`cat`, `is`, `sneaky`, `evil`, ``, `apple`)
	it2 := Sort(it1, func(x, y string) int {
		return len(x) - len(y)
	})
	checkIt(t, it2, ``, `is`, `cat`, `evil`, `apple`, `sneaky`)

	it1 = Iterate(`cat`, `is`, `sneaky`, `evil`, ``, `apple`)
	it2 = Sort(it1)
	checkIt(t, it2, ``, `apple`, `cat`, `evil`, `is`, `sneaky`)

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`,
		func() { Sort(it1, utils.OrderedComparer[string](), utils.OrderedComparer[string]()) })
}

func Test_Iterator_Sorted(t *testing.T) {
	it := Iterate(1, 2, 3, 4, 5)
	checkEqual(t, true, Sorted(it, utils.OrderedComparer[int]()))
	checkIt(t, it) // Remainder

	it = Iterate(1, 2, 5, 4, 3)
	checkEqual(t, false, Sorted(it, utils.OrderedComparer[int]()))
	checkIt(t, it, 3) // Remainder
}

func watcher[T any](count *int, it collections.Iterator[T]) collections.Iterator[T] {
	return New(func() (T, bool) {
		*count++
		hasNext := it.Next()
		return it.Current(), hasNext
	})
}

func checkIt[T any](t *testing.T, it collections.Iterator[T], exp ...T) {
	var parts []T
	for it.Next() {
		parts = append(parts, it.Current())
	}
	checkEqual(t, exp, parts)
	checkZero(t, it.Current())
}

func checkEqual(t testing.TB, exp, actual any) {
	t.Helper()
	if !utils.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value didn't match the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkZero(t testing.TB, actual any) {
	t.Helper()
	if !utils.IsZero(actual) {
		t.Errorf("\n"+
			"Expected value to be zero:\n"+
			"Actual:   %v (%T)", actual, actual)
	}
}

func checkLength(t testing.TB, exp int, value any) {
	t.Helper()
	actual, ok := utils.Length(value)
	if !ok {
		t.Errorf("\n"+
			"The given value does not have a length\n"+
			"Value:   %v (%T)", value, value)
	}
	if exp != actual {
		t.Errorf("\n"+
			"Expected value didn't match the actual value:\n"+
			"Value:    %v (%T)\n"+
			"Actual:   %d\n"+
			"Expected: %v (%T)", value, value, actual, exp, exp)
	}
}

func checkPanic(t testing.TB, exp string, handle func()) {
	t.Helper()
	actual := func() (r string) {
		defer func() { r = utils.String(recover()) }()
		handle()
		t.Error(`expected a panic`)
		return ``
	}()
	checkEqual(t, exp, actual)
}
