package readonlyVariantList

import (
	"bytes"
	"sort"
	"testing"

	"goToolbox/utils"
)

func Test_ReadonlyVariantList_Fallbacks(t *testing.T) {
	rv := (*impReadonlyVariantList[any])(nil)
	checkEqual(t, 0, rv.Count(), `Count() on nil struct`)
	checkEqual(t, nil, rv.liteGet(0), `liteGet(0) on nil struct`)

	rv = &impReadonlyVariantList[any]{
		countHandle: nil,
		getHandle:   func(i int) any { return 1 },
	}
	checkEqual(t, 0, rv.Count(), `Count() on nil "count"`)
	checkEqual(t, 1, rv.liteGet(0), `liteGet(0) on nil "count"`)

	rv = &impReadonlyVariantList[any]{
		countHandle: func() int { return 10 },
		getHandle:   nil,
	}
	checkEqual(t, 10, rv.Count(), `Count() on nil "get"`)
	checkEqual(t, nil, rv.Get(0), `Get(0) on nil "get"`)

	r2 := From[any](nil, nil)
	checkEqual(t, 0, r2.Count(), `Count() on nil struct via "From(nil, nil)"`)

	r2 = Wrap(nil)
	checkEqual(t, 0, r2.Count(), `Count() on nil struct via "Wrap(nil)"`)

	r2 = Wrap(12)
	checkEqual(t, 1, r2.Count(), `Count() on single value`)
	checkEqual(t, 12, r2.Get(0), `Get(0) on single value`)

	str := `Hello`
	pStr := &str
	r2 = Wrap(pStr)
	checkEqual(t, 1, r2.Count(), `Count() on single pointer to string`)
	checkEqual(t, pStr, r2.Get(0), `Get(0) on single pointer to string`)
}

func Test_ReadonlyVariantList_String(t *testing.T) {
	rv := Wrap(`Hello World`)
	checkEqual(t, 11, rv.Count(), `Count()`)
	checkEqual(t, false, rv.Empty(), `Empty()`)
	checkEqual(t, byte('H'), rv.Get(0), `Get(0)`)
	checkEqual(t, byte('d'), rv.Get(10), `Get(10)`)
	checkEqual(t, byte('H'), rv.First(), `First()`)
	checkEqual(t, byte('d'), rv.Last(), `Last()`)

	value, ok := rv.TryGet(1)
	checkEqual(t, byte('e'), value, `value, _ := TryGet(1)`)
	checkEqual(t, true, ok, `_, ok := TryGet(1)`)

	value, ok = rv.TryGet(-1)
	checkEqual(t, nil, value, `value, _ := TryGet(-1)`)
	checkEqual(t, false, ok, `_, ok := TryGet(-1)`)

	value, ok = rv.TryGet(11)
	checkEqual(t, nil, value, `value, _ := TryGet(11)`)
	checkEqual(t, false, ok, `_, ok := TryGet(11)`)

	checkPanic(t, `index out of bounds {count: 11, index: -1}`, `Get(-1)`, func() { rv.Get(-1) })
	checkPanic(t, `index out of bounds {count: 11, index: 11}`, `Get(11)`, func() { rv.Get(11) })

	checkEqual(t, true, rv.StartsWith(Wrap(`Hell`)), `StartsWith("Hell)`)
	checkEqual(t, false, rv.StartsWith(Wrap(`Help`)), `StartsWith("Help")`)
	checkEqual(t, true, rv.EndsWith(Wrap(`World`)), `EndsWith("World")`)
	checkEqual(t, false, rv.EndsWith(Wrap(`Cord`)), `EndsWith("Cord")`)

	checkEqual(t, []any{
		byte('H'), byte('e'), byte('l'), byte('l'), byte('o'), byte(' '),
		byte('W'), byte('o'), byte('r'), byte('l'), byte('d')},
		rv.ToSlice(), `Enumerate() and ToSlice()`)
	checkEqual(t, []any{
		byte('d'), byte('l'), byte('r'), byte('o'), byte('W'), byte(' '),
		byte('o'), byte('l'), byte('l'), byte('e'), byte('H')},
		rv.Backwards().ToSlice(), `Backwards()`)

	sc := make([]any, 3)
	rv.CopyToSlice(sc)
	checkEqual(t, []any{byte('H'), byte('e'), byte('l')}, sc, `CopyToSlice([]any{3})`)

	checkEqual(t, true, rv.Contains(byte('W')), `Contains(byte('W'))`)
	checkEqual(t, 6, rv.IndexOf(byte('W')), `Empty(byte('W'))`)
	checkEqual(t, false, rv.Contains(byte('X')), `Contains(byte('X'))`)
	checkEqual(t, -1, rv.IndexOf(byte('X')), `Empty(byte('X'))`)
	checkEqual(t, `72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100`, rv.String(), `String()`)

	r2 := Wrap(`Hello World`)
	checkEqual(t, true, rv.Equals(r2), `Equals new`)
	r2 = Wrap(rv.ToSlice())
	checkEqual(t, true, rv.Equals(r2), `Equals []any`)
	r2 = Wrap(`Hello Worlb`)
	checkEqual(t, false, rv.Equals(r2), `not equals wrong letter`)
	r2 = Wrap(`Hello World!`)
	checkEqual(t, false, rv.Equals(r2), `not equals extra letter`)

	r2 = Wrap(``)
	checkEqual(t, 0, r2.Count(), `Count()`)
	checkEqual(t, true, r2.Empty(), `Empty()`)
	checkPanic(t, `collection contains no values {action: First}`, `First()`, func() { r2.First() })
	checkPanic(t, `collection contains no values {action: Last}`, `Last()`, func() { r2.Last() })
}

func Test_ReadonlyVariantList_Slice(t *testing.T) {
	rv := Wrap([]byte{255, 12, 42, 39})
	checkEqual(t, 4, rv.Count(), `Count()`)
	checkEqual(t, byte(255), rv.Get(0), `Get(0)`)

	rv = Wrap([]int{1355, 112, 42, 399})
	checkEqual(t, 4, rv.Count(), `Count()`)
	checkEqual(t, 1355, rv.Get(0), `Get(0)`)

	rv = Wrap([]string{`cat`, `dog`, `Gizmo`})
	checkEqual(t, 3, rv.Count(), `Count()`)
	checkEqual(t, `cat`, rv.Get(0), `Get(0)`)

	rv = Wrap([]rune{'A', 'B', 'C'})
	checkEqual(t, 3, rv.Count(), `Count()`)
	checkEqual(t, 'A', rv.Get(0), `Get(0)`)

	rv = Wrap([]any{'A', 134, `dog`, byte(12)})
	checkEqual(t, 4, rv.Count(), `Count()`)
	checkEqual(t, 'A', rv.Get(0), `Get(0)`)

	buf := &bytes.Buffer{}
	_, _ = buf.WriteString(`Giggle`)
	rv = Wrap(buf)
	checkEqual(t, 6, rv.Count(), `Count()`)
	checkEqual(t, byte('G'), rv.Get(0), `Get(0)`)

	rv = Wrap(&buf) // double pointer
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, &buf, rv.Get(0), `Get(0)`)

	rv = Wrap([]float64{3.14, 2.2, 8.87, 14.6, 56.7})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 3.14, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoSliceableP{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoSliceableR{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoSliceableP{})
	// doesn't have access to pointer method so used as a single value
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, pseudoSliceableP{}, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoSliceableR{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)
}

func Test_ReadonlyVariantList_List(t *testing.T) {
	rv := Wrap(&pseudoListPP{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoListPR{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoListRP{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoListRR{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoListPP{})
	// doesn't have access to pointer methods so used as a single value
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, pseudoListPP{}, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoListPR{})
	// doesn't have access to the pointer count method so used as a single value
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, pseudoListPR{}, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoListRP{})
	// doesn't have access to the pointer get method so used as a single value
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, pseudoListRP{}, rv.Get(0), `Get(0)`)

	rv = Wrap(pseudoListRR{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	p1 := &pseudoOnlyCount{}
	rv = Wrap(p1)
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, p1, rv.Get(0), `Get(0)`)

	p2 := &pseudoOnlyGet{}
	rv = Wrap(p2)
	checkEqual(t, 1, rv.Count(), `Count()`)
	checkEqual(t, p2, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoListLen{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)

	rv = Wrap(&pseudoListLength{})
	checkEqual(t, 5, rv.Count(), `Count()`)
	checkEqual(t, 15, rv.Get(0), `Get(0)`)
}

func Test_ReadonlyVariantList_Map(t *testing.T) {
	rv := Wrap(map[string]int{`one`: 1, `two`: 2, `three`: 3})
	checkEqual(t, 3, rv.Count(), `Count()`)
	results := utils.Strings([]any{rv.Get(0), rv.Get(1), rv.Get(2)})
	sort.Strings(results)

	checkEqual(t, `[one, 1]`, results[0], `results[0]`)
	checkEqual(t, `[three, 3]`, results[1], `results[1]`)
	checkEqual(t, `[two, 2]`, results[2], `results[2]`)
}

func Test_ReadonlyVariantList_Cast(t *testing.T) {
	rv := Cast[int](Wrap([]int{44, 22, 55}))
	checkEqual(t, 3, rv.Count(), `Count()`)
	checkEqual(t, 44, rv.Get(0), `Get(0)`)

	rv = Cast[int](Wrap([]float64{4.4, 2.2, 5.5}), func(v any) int { return int(v.(float64)) })
	checkEqual(t, 3, rv.Count(), `Count()`)
	checkEqual(t, 4, rv.Get(0), `Get(0)`)

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: selector}`,
		`too many selectors`, func() {
			rv = Cast[int](Wrap([]int{44, 22, 55}), func(v any) int { return 0 }, func(v any) int { return 1 })
		})
}

func Test_ReadonlyVariantList_UnstableIteration(t *testing.T) {
	s := &pseudoSliceWrapper{}
	s.slice = []int{1, 2, 3, 4, 5}
	rv := Wrap(s)

	it1 := rv.Enumerate().Iterate()
	checkEqual(t, true, it1.Next(), `1st forward next`)
	checkEqual(t, 1, it1.Current(), `after 1st forward next`)

	it2 := rv.Backwards().Iterate()
	checkEqual(t, true, it2.Next(), `1st backward next`)
	checkEqual(t, 5, it2.Current(), `after 1st backward next`)

	s.slice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // acts like insertion

	checkEqual(t, true, it1.Next(), `2nd forward next`)
	checkEqual(t, 2, it1.Current(), `after 2nd forward next`)

	checkEqual(t, true, it2.Next(), `2nd backward next`)
	checkEqual(t, 4, it2.Current(), `after 2nd backward next`)

	s.slice = []int{7, 8, 9} // acts like removal

	checkEqual(t, true, it1.Next(), `3rd forward next`)
	checkEqual(t, 9, it1.Current(), `after 3rd forward next`)
	checkEqual(t, false, it1.Next(), `4th forward next`)

	checkEqual(t, true, it2.Next(), `3rd backward next`)
	checkEqual(t, 9, it2.Current(), `after 3rd backward next`)
	checkEqual(t, true, it2.Next(), `4th backward next`)
	checkEqual(t, 8, it2.Current(), `after 4th backward next`)
	checkEqual(t, true, it2.Next(), `5th backward next`)
	checkEqual(t, 7, it2.Current(), `after 5th backward next`)
	checkEqual(t, false, it2.Next(), `6th backward next`)
}

func checkEqual(t *testing.T, exp, actual any, name string) {
	t.Helper()
	if !utils.Equal(actual, exp) {
		t.Errorf("\nUnexpected result from %q:\n"+
			"\tActual:   %v (%T)\n"+
			"\tExpected: %v (%T)", name, actual, actual, exp, exp)
	}
}

func checkPanic(t *testing.T, exp, name string, handle func()) {
	t.Helper()
	r := func() (r any) {
		defer func() { r = recover() }()
		handle()
		return `no panic`
	}()
	actual := utils.String(r)
	checkEqual(t, exp, actual, name)
}

var data = []int{15, 51, 24, 42, 33}

type pseudoSliceableP struct{}

func (s *pseudoSliceableP) ToSlice() []int { return data }

type pseudoSliceableR struct{}

func (s pseudoSliceableR) ToSlice() []int { return data }

type pseudoListPP struct{}

func (pl *pseudoListPP) Count() int    { return len(data) }
func (pl *pseudoListPP) Get(i int) int { return data[i] }

type pseudoListPR struct{}

func (pl *pseudoListPR) Count() int   { return len(data) }
func (pl pseudoListPR) Get(i int) int { return data[i] }

type pseudoListRP struct{}

func (pl pseudoListRP) Count() int     { return len(data) }
func (pl *pseudoListRP) Get(i int) int { return data[i] }

type pseudoListRR struct{}

func (pl pseudoListRR) Count() int    { return len(data) }
func (pl pseudoListRR) Get(i int) int { return data[i] }

type pseudoOnlyCount struct{}

func (pl *pseudoOnlyCount) Count() int { return len(data) }

type pseudoOnlyGet struct{}

func (pl *pseudoOnlyGet) Get(i int) int { return data[i] }

type pseudoListLen struct{}

func (pl pseudoListLen) Len() int      { return len(data) }
func (pl pseudoListLen) Get(i int) int { return data[i] }

type pseudoListLength struct{}

func (pl pseudoListLength) Length() int   { return len(data) }
func (pl pseudoListLength) Get(i int) int { return data[i] }

type pseudoSliceWrapper struct{ slice []int }

func (psw *pseudoSliceWrapper) Count() int    { return len(psw.slice) }
func (psw *pseudoSliceWrapper) Get(i int) int { return psw.slice[i] }
