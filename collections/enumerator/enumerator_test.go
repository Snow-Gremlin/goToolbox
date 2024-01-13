package enumerator

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type pseudoComparable[T cmp.Ordered] struct {
	value T
}

func pc[T cmp.Ordered](value T) *pseudoComparable[T] {
	return &pseudoComparable[T]{value: value}
}

func (p *pseudoComparable[T]) CompareTo(other *pseudoComparable[T]) int {
	return utils.OrderedComparer[T]()(p.value, other.value)
}

func (p *pseudoComparable[T]) String() string {
	return utils.String(p.value)
}

func Test_Enumerator_Enumerate(t *testing.T) {
	e := Enumerate(`apple`, `dog`, `cat`, `bat`, `hat`)
	checkEqual(t, []string{`apple`, `dog`, `cat`, `bat`, `hat`}, e.ToSlice())
	checkLength(t, 5, e)
	checkEqual(t, []string{`cat`, `bat`, `hat`}, e.Where(predicate.Matches(`^.at$`)).ToSlice())

	e = Enumerate[string]()
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())
}

func Test_Enumerator_Range(t *testing.T) {
	e := Range(0, 5)
	checkEqual(t, []int{0, 1, 2, 3, 4}, e.ToSlice())
	checkLength(t, 5, e)

	e = Range(-3, 5)
	checkEqual(t, []int{-3, -2, -1, 0, 1}, e.ToSlice())
	checkLength(t, 5, e)

	e = Range(6, 0)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())

	e = Range(6, -5)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())
}

func Test_Enumerator_Stride(t *testing.T) {
	e := Stride(0, 3, 5)
	checkEqual(t, []int{0, 3, 6, 9, 12}, e.ToSlice())
	checkLength(t, 5, e)

	e = Stride(-3, 2, 5)
	checkEqual(t, []int{-3, -1, 1, 3, 5}, e.ToSlice())
	checkLength(t, 5, e)

	e = Stride(6, 2, 0)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())

	e = Stride(6, 2, -5)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())
}

func Test_Enumerator_Repeat(t *testing.T) {
	e := Repeat(`A`, 5)
	checkEqual(t, []string{`A`, `A`, `A`, `A`, `A`}, e.ToSlice())
	checkLength(t, 5, e)

	e = Repeat(`A`, 0)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())

	e = Repeat(`A`, -8)
	checkLength(t, 0, e)
	checkEqual(t, true, e.Empty())
}

func Test_Enumerator_Split(t *testing.T) {
	e := Split(`Cat dog hot cold mouse`, ` `)
	checkEqual(t, []string{`Cat`, `dog`, `hot`, `cold`, `mouse`}, e.ToSlice())
	checkLength(t, 5, e)

	checkEqual(t, []string{`Cat`, ` dog`, ` hot`, ` cold`, ``, ``, `mouse`},
		Split(`Cat, dog, hot, cold,,,mouse`, `,`).ToSlice())

	checkEqual(t, []string{`Cat dog hot cold mouse`},
		Split(`Cat dog hot cold mouse`, `#`).ToSlice())

	checkEqual(t, []string{},
		Split(``, ` `).ToSlice())

	checkEqual(t, []string{``},
		Split(` `, ` `).ToSlice())
}

func Test_Enumerator_SplitFunc(t *testing.T) {
	buf := &bytes.Buffer{}
	it := SplitFunc(`Cat dog mouse`, func(part string) (int, int) {
		_, _ = buf.WriteString(`[` + part + `]`)
		return strings.Index(part, ` `), 1
	}).Iterate()

	checkEqual(t, true, it.Next())
	checkEqual(t, `Cat`, it.Current())
	checkEqual(t, `[Cat dog mouse]`, buf.String())
	buf.Reset()

	checkEqual(t, true, it.Next())
	checkEqual(t, `dog`, it.Current())
	checkEqual(t, `[dog mouse]`, buf.String())
	buf.Reset()

	checkEqual(t, true, it.Next())
	checkEqual(t, `mouse`, it.Current())
	checkEqual(t, `[mouse]`, buf.String())
	buf.Reset()

	checkEqual(t, false, it.Next())
	checkEqual(t, ``, it.Current())
	checkEqual(t, ``, buf.String())

	checkPanic(t, `argument may not be nil {name: separator}`, func() {
		SplitFunc(`boom`, nil)
	})
}

func Test_Enumerator_Lines(t *testing.T) {
	e := Lines("Cat\ndog\n\nhot\ncold\nmouse")
	checkEqual(t, []string{`Cat`, `dog`, ``, `hot`, `cold`, `mouse`}, e.ToSlice())
	checkLength(t, 6, e)

	checkEqual(t, []string{`Cat`, `dog`, ``, `hot`, `cold`, `mouse`},
		Lines("Cat\n\rdog\r\rhot\r\ncold\u2029mouse").ToSlice())
}

func Test_Enumerator_Error(t *testing.T) {
	e := Errors(fmt.Errorf(`%w-%w-%w`, terror.New(`One`, errors.New(`Two`)), errors.New(`Three`), terror.New(`Four`)))
	checkEqual(t, []string{`One: Two-Three-Four`, `One: Two`, `Two`, `Three`, `Four`}, e.Strings().ToSlice())
	checkLength(t, 5, e)
}

func Test_Enumerator_Select(t *testing.T) {
	e1 := Select(Range(0, 10), func(i int) float64 {
		return math.Pi * float64(i) / 10.0
	})
	e2 := Select(e1, func(f float64) float64 {
		return (1.0 - math.Cos(f)) * 50.0
	})
	e3 := Select(e2, func(f float64) int {
		return int(math.Round(f))
	})
	checkEqual(t, []int{0, 2, 10, 21, 35, 50, 65, 79, 90, 98}, e3.ToSlice())
	checkLength(t, 10, e3)
}

func Test_Enumerator_Expand(t *testing.T) {
	e := Expand(Range(1, 4), func(i int) collections.Iterable[int] {
		return Repeat(i, i).Iterate
	})
	checkEqual(t, []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}, e.ToSlice())
	checkLength(t, 10, e)
}

func Test_Enumerator_Reduce(t *testing.T) {
	e1 := Enumerate(`horse`, `cat`, `mouse`, `wolf`)
	value := Reduce(e1, 100, func(s string, prior int) int {
		return prior + len(s)
	})
	checkEqual(t, 117, value)

	values := Reduce(e1, []int{}, func(s string, prior []int) []int {
		return append(prior, len(s))
	})
	checkEqual(t, []int{5, 3, 5, 4}, values)
}

func Test_Enumerator_Merge(t *testing.T) {
	e1 := Enumerate(`horse`, `cat`, `mouse`, `wolf`)
	value := e1.Merge(func(s, prior string) string {
		return fmt.Sprintf(`%s|%s`, s, prior)
	})
	checkEqual(t, `wolf|mouse|cat|horse`, value)
}

func Test_Enumerator_SlidingWindow(t *testing.T) {
	e1 := Select(Range(0, 20), func(i int) int {
		return int(math.Round((1.0 - math.Cos(math.Pi*float64(i)/20.0)) * 50.0))
	})
	checkEqual(t, []int{0, 1, 2, 5, 10, 15, 21, 27, 35, 42, 50, 58, 65, 73, 79, 85, 90, 95, 98, 99}, e1.ToSlice())

	e2 := SlidingWindow(e1, 2, 1, func(v []int) int {
		checkLength(t, 2, v)
		return v[1] - v[0]
	})
	checkEqual(t, []int{1, 1, 3, 5, 5, 6, 6, 8, 7, 8, 8, 7, 8, 6, 6, 5, 5, 3, 1}, e2.ToSlice())
	checkLength(t, 19, e2)

	e3 := SlidingWindow(e2, 4, 1, func(v []int) float64 {
		checkLength(t, 4, v)
		return float64(v[0]+v[1]+v[2]+v[3]) / 4.0
	})
	checkEqual(t, []float64{2.5, 3.5, 4.75, 5.5, 6.25, 6.75, 7.25, 7.75, 7.5, 7.75, 7.25, 6.75, 6.25, 5.5, 4.75, 3.5}, e3.ToSlice())
	checkLength(t, 16, e3)

	e4 := SlidingWindow(e3, 10, 1, func(v []float64) int {
		checkLength(t, 10, v)
		sum := 0.0
		for _, u := range v {
			sum += u
		}
		return int(sum)
	})
	checkEqual(t, []int{59, 64, 67, 69, 69, 67, 64}, e4.ToSlice())
	checkLength(t, 7, e4)

	e5 := SlidingWindow(e4, 40, 1, func(v []int) int {
		t.Fatal(`Shouldn't be called`)
		return 0
	})
	checkLength(t, 0, e5)

	checkPanic(t, `the given window size must be greater than zero {size: -6}`,
		func() {
			e6 := SlidingWindow(e4, -6, 1, func(v []int) int {
				t.Fatal(`Shouldn't be called`)
				return 0
			})
			checkLength(t, 0, e6) // not reached
		})
}

func Test_Enumerator_Chunk(t *testing.T) {
	e1 := Enumerate(1, 2, 3, 4, 5, 6, 7, 8, 9)

	e2 := Chunk(e1, 3)
	checkEqual(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, e2.ToSlice())
	checkLength(t, 3, e2)

	e3 := Chunk(e1, 4)
	checkEqual(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9}}, e3.ToSlice())
	checkLength(t, 3, e3)
}

func Test_Enumerator_Max(t *testing.T) {
	checkEqual(t, `wolf`, Enumerate(`horse`, `wolf`, `cat`, `mouse`).Max())
	checkEqual(t, `wolf`, Enumerate(`Wolf`, `wolf`, `WOLF`).Max())
	checkEqual(t, 76, Enumerate(45, 76, 2, 56, 45, 5).Max())
	checkEqual(t, 5.6, Enumerate(0.45, 1.76, 2.0, 5.6, 0.45, 0.5).Max())

	e1 := Enumerate(pc(`horse`), pc(`wolf`), pc(`cat`), pc(`mouse`))
	checkEqual(t, pc(`wolf`), e1.Max())

	checkEqual(t, 76, Enumerate(45, 76, 2, 56, 45, 5).Max(utils.DefaultComparer[int]()))
	checkEqual(t, 2, Enumerate(45, 76, 2, 56, 45, 5).Max(utils.Descender(utils.DefaultComparer[int]())))

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, func() {
		Enumerate(45, 76, 2, 56, 45, 5).Max(utils.DefaultComparer[int](), utils.DefaultComparer[int]())
	})
}

func Test_Enumerator_Min(t *testing.T) {
	checkEqual(t, `cat`, Enumerate(`horse`, `wolf`, `cat`, `mouse`).Min())
	checkEqual(t, `WOLF`, Enumerate(`Wolf`, `wolf`, `WOLF`).Min())
	checkEqual(t, 2, Enumerate(45, 76, 2, 56, 45, 5).Min())
	checkEqual(t, 0.45, Enumerate(0.45, 1.76, 2.0, 5.6, 0.45, 0.5).Min())

	e1 := Enumerate(pc(`horse`), pc(`wolf`), pc(`cat`), pc(`mouse`))
	checkEqual(t, pc(`cat`), e1.Min())

	checkEqual(t, 2, Enumerate(45, 76, 2, 56, 45, 5).Min(utils.DefaultComparer[int]()))
	checkEqual(t, 76, Enumerate(45, 76, 2, 56, 45, 5).Min(utils.Descender(utils.DefaultComparer[int]())))

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, func() {
		Enumerate(45, 76, 2, 56, 45, 5).Min(utils.DefaultComparer[int](), utils.DefaultComparer[int]())
	})
}

func Test_Enumerator_Sum(t *testing.T) {
	sum1, count1 := Sum(Enumerate(45, 76, 2, 56, 45, 5))
	checkEqual(t, 229, sum1)
	checkEqual(t, 6, count1)

	sum2, count2 := Sum(Enumerate(0.25, 1.75, 2.0, 5.125, 0.25, 0.5))
	checkEqual(t, 9.875, sum2)
	checkEqual(t, 6, count2)
}

func Test_Enumerator_IsUnique(t *testing.T) {
	e := Enumerate(`cat`, `cat`, `wolf`, `cat`, `mouse`, `wolf`, `cat`)
	checkEqual(t, false, IsUnique(e))

	e = Enumerate(`cat`, `wolf`, `mouse`)
	checkEqual(t, true, IsUnique(e))
}

func Test_Enumerator_Unique(t *testing.T) {
	e := Unique(Enumerate(`cat`, `cat`, `wolf`, `cat`, `mouse`, `wolf`, `cat`))
	checkEqual(t, []string{`cat`, `wolf`, `mouse`}, e.ToSlice())
	checkLength(t, 3, e)
}

func Test_Enumerator_DuplicateCounts(t *testing.T) {
	e := Enumerate(`cat`, `cat`, `wolf`, `cat`, `mouse`, `wolf`, `cat`)
	d := DuplicateCounts(e)
	keys := utils.SortedKeys(d)
	checkEqual(t, []string{`cat`, `mouse`, `wolf`}, keys)
	checkEqual(t, 4, d[`cat`])
	checkEqual(t, 2, d[`wolf`])
	checkEqual(t, 1, d[`mouse`])
}

func Test_Enumerator_Intersection_Union_Subtract(t *testing.T) {
	e1 := Enumerate(1, 3, 5, 7, 9)
	e2 := Enumerate(2, 4, 6, 8)
	checkEqual(t, []int{}, Intersection(e1, e2).ToSlice())
	checkEqual(t, []int{}, Intersection(e2, e1).ToSlice())
	checkEqual(t, []int{1, 3, 5, 7, 9, 2, 4, 6, 8}, Union(e1, e2).ToSlice())
	checkEqual(t, []int{2, 4, 6, 8, 1, 3, 5, 7, 9}, Union(e2, e1).ToSlice())
	checkEqual(t, []int{2, 4, 6, 8}, Subtract(e1, e2).ToSlice())
	checkEqual(t, []int{1, 3, 5, 7, 9}, Subtract(e2, e1).ToSlice())

	e2 = Enumerate(1, 2, 3, 4, 5)
	checkEqual(t, []int{1, 3, 5}, Intersection(e1, e2).ToSlice())
	checkEqual(t, []int{1, 3, 5}, Intersection(e2, e1).ToSlice())
	checkEqual(t, []int{1, 3, 5, 7, 9, 2, 4}, Union(e1, e2).ToSlice())
	checkEqual(t, []int{1, 2, 3, 4, 5, 7, 9}, Union(e2, e1).ToSlice())
	checkEqual(t, []int{2, 4}, Subtract(e1, e2).ToSlice())
	checkEqual(t, []int{7, 9}, Subtract(e2, e1).ToSlice())
}

func Test_Enumerator_Zip(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkLength(t, 6, e1)
	e2 := Enumerate(5, 4, 1, 76, 2, 678, 3, 12, 52, 99)
	checkLength(t, 10, e2)
	e3 := ZipToTuples(e1, e2)
	checkEqual(t, `[cat, 5], [bat, 4], [wolf, 1], [hat, 76], [mouse, 2], [dog, 678]`, e3.Join(`, `))
	checkLength(t, 6, e3)
}

func Test_Enumerator_Interweave(t *testing.T) {
	e1 := Enumerate(`I`, `II`, `III`, `IV`, `V`, `VI`, `VII`, `VIII`, `IX`, `X`)
	e2 := Enumerate(`one`, `two`, `three`, `four`, `five`)
	e3 := Enumerate(`1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`)
	e4 := Interweave(e1, e2, e3)
	checkEqual(t, []string{
		`I`, `one`, `1`, `II`, `two`, `2`, `III`, `three`, `3`, `IV`, `four`, `4`,
		`V`, `five`, `5`, `VI`, `6`, `VII`, `7`, `VIII`, `8`, `IX`, `X`,
	}, e4.ToSlice())
	checkLength(t, 23, e4)
}

func Test_Enumerator_SortInterweave(t *testing.T) {
	e1 := Enumerate(1, 3, 5, 7, 9)
	e2 := Enumerate(2, 4, 6, 8)
	checkEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, e1.SortInterweave(e2).ToSlice())
	checkEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, e2.SortInterweave(e1).ToSlice())

	e1 = Enumerate(1, 2, 3)
	e2 = Enumerate(4, 5, 6)
	checkEqual(t, []int{1, 2, 3, 4, 5, 6}, e1.SortInterweave(e2).ToSlice())
	checkEqual(t, []int{1, 2, 3, 4, 5, 6}, e2.SortInterweave(e1).ToSlice())

	// Not sorted enumerators merging
	e1 = Enumerate(5, 3, 1)
	e2 = Enumerate(2, 4, 6)
	checkEqual(t, []int{2, 4, 5, 3, 1, 6}, e1.SortInterweave(e2).ToSlice())
	checkEqual(t, []int{2, 4, 5, 3, 1, 6}, e2.SortInterweave(e1).ToSlice())

	e3 := Enumerate(pc(1), pc(3), pc(5), pc(7), pc(9))
	e4 := Enumerate(pc(2), pc(4), pc(6), pc(8))
	checkEqual(t, []*pseudoComparable[int]{
		pc(1), pc(2), pc(3), pc(4),
		pc(5), pc(6), pc(7), pc(8), pc(9),
	}, e3.SortInterweave(e4).ToSlice())

	checkEqual(t, []int{2, 4, 5, 3, 1, 6}, e1.SortInterweave(e2, utils.DefaultComparer[int]()).ToSlice())
	checkEqual(t, []int{5, 3, 2, 4, 6, 1}, e1.SortInterweave(e2, utils.Descender(utils.DefaultComparer[int]())).ToSlice())

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, func() {
		e1.SortInterweave(e2, utils.DefaultComparer[int](), utils.DefaultComparer[int]())
	})
}

func Test_Enumerator_Indexed(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	e2 := Indexed(e1)
	checkEqual(t, `[0, cat], [1, bat], [2, wolf], [3, hat], [4, mouse], [5, dog]`, e2.Join(`, `))
	checkLength(t, 6, e2)
}

func Test_Enumerator_Sorted(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, false, e1.Sorted())

	e1 = Enumerate(`bat`, `cat`, `dog`, `hat`, `mouse`, `wolf`)
	checkEqual(t, true, e1.Sorted())

	e1 = Enumerate(`cat`)
	checkEqual(t, true, e1.Sorted())

	e1 = Enumerate[string]()
	checkEqual(t, true, e1.Sorted())

	e2 := Enumerate(pc(`cat`), pc(`bat`), pc(`wolf`), pc(`hat`), pc(`mouse`), pc(`dog`))
	checkEqual(t, false, e2.Sorted())

	e2 = Enumerate(pc(`bat`), pc(`cat`), pc(`dog`), pc(`hat`), pc(`mouse`), pc(`wolf`))
	checkEqual(t, true, e2.Sorted())

	e1 = Enumerate(`bat`, `cat`, `dog`, `hat`, `mouse`, `wolf`)
	checkEqual(t, true, e1.Sorted(strings.Compare))
	checkEqual(t, false, e1.Sorted(utils.Descender(strings.Compare)))

	e1 = Enumerate(`wolf`, `mouse`, `hat`, `dog`, `cat`, `bat`)
	checkEqual(t, false, e1.Sorted(strings.Compare))
	checkEqual(t, true, e1.Sorted(utils.Descender(strings.Compare)))

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, func() {
		e1.Sorted(strings.Compare, strings.Compare)
	})
}

func Test_Enumerator_Sort(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`).Sort()
	checkEqual(t, []string{`bat`, `cat`, `dog`, `hat`, `mouse`, `wolf`}, e1.ToSlice())
	checkLength(t, 6, e1)

	e2 := Enumerate(pc(`Giz`), pc(`Mo`), pc(`Gizmo`)).Sort().Strings()
	checkEqual(t, []string{`Giz`, `Gizmo`, `Mo`}, e2.ToSlice())

	e3 := Enumerate(`cat`, `horse`, `wolf`, `elephant`, `mouse`, `dog`).Sort(func(x, y string) int {
		return len(x) - len(y)
	})
	checkEqual(t, []string{`cat`, `dog`, `wolf`, `horse`, `mouse`, `elephant`}, e3.ToSlice())
	checkLength(t, 6, e3)

	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, func() {
		Enumerate(`cat`, `bat`).Sort(strings.Compare, strings.Compare)
	})
}

func Test_Enumerator_Where(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`).
		Where(func(value string) bool { return len(value) == 3 })
	checkEqual(t, []string{`cat`, `bat`, `hat`, `dog`}, e1.ToSlice())
}

func Test_Enumerator_WhereNot(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`).
		WhereNot(func(value string) bool { return len(value) == 3 })
	checkEqual(t, []string{`wolf`, `mouse`}, e1.ToSlice())
}

func Test_Enumerator_NotNil(t *testing.T) {
	v1, v2, v3 := 23, 15, 55
	e1 := Enumerate[*int](nil, &v1, nil, &v2, nil, nil, &v3, nil).NotNil()
	checkEqual(t, []*int{&v1, &v2, &v3}, e1.ToSlice())
}

func Test_Enumerator_NotZero(t *testing.T) {
	e1 := Enumerate(``, `cat`, ``, `bat`, `wolf`, ``, ``).NotZero()
	checkEqual(t, []string{`cat`, `bat`, `wolf`}, e1.ToSlice())
}

func Test_Enumerator_ToSlice(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`}, e1.ToSlice())
}

func Test_Enumerator_CopyToSlice(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)

	s := make([]string, 4)
	e1.CopyToSlice(s)
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`}, s)

	s = make([]string, 8)
	e1.CopyToSlice(s)
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`, ``, ``}, s)
}

func Test_Enumerator_Foreach(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	total := 0
	e1.Foreach(func(value string) {
		total += len(value)
	})
	checkEqual(t, 21, total)
}

func Test_Enumerator_DoUntilError(t *testing.T) {
	results := []int{}
	parseHex := func(value string) error {
		v, err := strconv.ParseInt(value, 16, 32)
		if err != nil {
			return err
		}
		results = append(results, int(v))
		return nil
	}

	e1 := Enumerate(`12`, `F`, `FF`)
	err := e1.DoUntilError(parseHex)
	checkEqual(t, nil, err)
	checkEqual(t, []int{18, 15, 255}, results)

	e2 := Enumerate(`1A`, `Cat`, `20`, `A2`)
	err = e2.DoUntilError(parseHex)
	checkEqual(t, `strconv.ParseInt: parsing "Cat": invalid syntax`, err.Error())
	checkEqual(t, []int{18, 15, 255, 26}, results)
}

func Test_Enumerator_DoUntilNotZero(t *testing.T) {
	parseBool := func(value string) bool {
		v, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		return v
	}

	e1 := Enumerate(`0`, `F`, `false`, `true`)
	v := DoUntilNotZero(e1, parseBool)
	checkEqual(t, true, v)
}

func Test_Enumerator_Any(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, true, e1.Any(predicate.Matches(`^.o..$`)))
	checkEqual(t, true, e1.Any(predicate.Eq(`hat`)))
	checkEqual(t, false, e1.Any(predicate.Eq(`cheese`)))
}

func Test_Enumerator_All(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, false, e1.All(predicate.Matches(`^.at$`)))

	e2 := Enumerate(`cat`, `bat`, `hat`)
	checkEqual(t, true, e2.All(predicate.Matches(`^.at$`)))
}

func Test_Enumerator_StepsUntil(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, 2, e1.StepsUntil(predicate.Matches(`^.o..$`)))
	checkEqual(t, 3, e1.StepsUntil(predicate.Eq(`hat`)))
	checkEqual(t, -1, e1.StepsUntil(predicate.Eq(`cheese`)))
}

func Test_Enumerator_Empty(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkEqual(t, false, e1.Empty())

	e2 := Enumerate[string]()
	checkEqual(t, true, e2.Empty())
}

func Test_Enumerator_Count(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkLength(t, 6, e1)

	e2 := Enumerate(`cat`, `bat`, `wolf`)
	checkLength(t, 3, e2)

	e3 := Enumerate[string]()
	checkLength(t, 0, e3)
}

func Test_Enumerator_AtLeast(t *testing.T) {
	e := Enumerate(`cat`, `bat`, `wolf`)
	checkEqual(t, false, e.AtLeast(4))
	checkEqual(t, true, e.AtLeast(3))
	checkEqual(t, true, e.AtLeast(2))
}

func Test_Enumerator_AtMost(t *testing.T) {
	e := Enumerate(`cat`, `bat`, `wolf`)
	checkEqual(t, true, e.AtMost(4))
	checkEqual(t, true, e.AtMost(3))
	checkEqual(t, false, e.AtMost(2))
}

func Test_Enumerator_First(t *testing.T) {
	result, ok := Enumerate(`cat`, `bat`, `wolf`).First()
	checkEqual(t, `cat`, result)
	checkEqual(t, true, ok)

	result, ok = Enumerate[string]().First()
	checkEqual(t, ``, result)
	checkEqual(t, false, ok)
}

func Test_Enumerator_Last(t *testing.T) {
	result, ok := Enumerate(`cat`, `bat`, `wolf`).Last()
	checkEqual(t, `wolf`, result)
	checkEqual(t, true, ok)

	result, ok = Enumerate[string]().Last()
	checkEqual(t, ``, result)
	checkEqual(t, false, ok)
}

func Test_Enumerator_Single(t *testing.T) {
	result, ok := Enumerate(`cat`, `bat`).Single()
	checkEqual(t, ``, result)
	checkEqual(t, false, ok)

	result, ok = Enumerate(`cat`).Single()
	checkEqual(t, `cat`, result)
	checkEqual(t, true, ok)

	result, ok = Enumerate[string]().Single()
	checkEqual(t, ``, result)
	checkEqual(t, false, ok)
}

func Test_Enumerator_Skip(t *testing.T) {
	e := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkLength(t, 6, e.Skip(-1))
	checkLength(t, 6, e.Skip(0))
	checkLength(t, 5, e.Skip(1))
	checkEqual(t, []string{`bat`, `wolf`, `hat`, `mouse`, `dog`}, e.Skip(1).ToSlice())
	checkLength(t, 3, e.Skip(3))
	checkEqual(t, []string{`hat`, `mouse`, `dog`}, e.Skip(3).ToSlice())
	checkLength(t, 0, e.Skip(6).ToSlice())
	checkLength(t, 0, e.Skip(8).ToSlice())
}

func Test_Enumerator_SkipWhile(t *testing.T) {
	e := Enumerate(1, 2, 3, 4, 5, 6, 7, 8, 9)
	checkLength(t, 9, e.SkipWhile(predicate.LessThan(-1)))
	checkLength(t, 9, e.SkipWhile(predicate.GreaterThan(10)))
	checkLength(t, 0, e.SkipWhile(predicate.LessThan(10)))
	checkLength(t, 8, e.SkipWhile(predicate.LessThan(2)))
	checkEqual(t, []int{4, 5, 6, 7, 8, 9}, e.SkipWhile(predicate.LessThan(4)).ToSlice())
	checkLength(t, 6, e.SkipWhile(predicate.LessThan(4)))
}

func Test_Enumerator_Take(t *testing.T) {
	e := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`)
	checkLength(t, 0, e.Take(-1))
	checkLength(t, 0, e.Take(0))
	checkLength(t, 1, e.Take(1))
	checkEqual(t, []string{`cat`}, e.Take(1).ToSlice())
	checkLength(t, 3, e.Take(3))
	checkEqual(t, []string{`cat`, `bat`, `wolf`}, e.Take(3).ToSlice())
	checkLength(t, 6, e.Take(6))
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`}, e.Take(6).ToSlice())
	checkLength(t, 6, e.Take(8))
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`}, e.Take(8).ToSlice())
}

func Test_Enumerator_TakeWhile(t *testing.T) {
	e := Enumerate(1, 2, 3, 4, 5, 6, 7, 8, 9)
	checkLength(t, 0, e.TakeWhile(predicate.GreaterThan(10)))
	checkLength(t, 0, e.TakeWhile(predicate.LessEq(-1)))
	checkLength(t, 1, e.TakeWhile(predicate.LessEq(1)))
	checkEqual(t, []int{1, 2, 3}, e.TakeWhile(predicate.LessEq(3)).ToSlice())
	checkLength(t, 3, e.TakeWhile(predicate.LessEq(3)))
}

func Test_Enumerator_Replace(t *testing.T) {
	e := Enumerate(1, 5, 7, 2, 6, 4, 3, 7, 8, 2, 1, 5, 6, 8)
	checkEqual(t, []int{1, 1, 3, 2, 2, 4, 3, 3, 4, 2, 1, 1, 2, 4},
		e.Replace(func(value int) int {
			if value < 5 {
				return value
			}
			return value - 4
		}).ToSlice())
}

func Test_Enumerator_Reverse(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`).Reverse()
	checkEqual(t, []string{`dog`, `mouse`, `hat`, `wolf`, `bat`, `cat`}, e1.ToSlice())

	e2 := Enumerate[string]().Reverse()
	checkLength(t, 0, e2.ToSlice())
}

type pseudoStringer struct{ text string }

func (p *pseudoStringer) String() string { return p.text }

func Test_Enumerator_ToStrings(t *testing.T) {
	s1 := &pseudoStringer{text: `plop`}
	s2 := fmt.Errorf(`bar`)
	e1 := Enumerate[any](1, `bat`, 3.25, s1, nil, s2, false).Strings()
	checkEqual(t, []string{`1`, `bat`, `3.25`, `plop`, `<nil>`, `bar`, `false`}, e1.ToSlice())
}

func Test_Enumerator_Trim(t *testing.T) {
	s1 := &pseudoStringer{text: "  plop\t"}
	e1 := Enumerate[any](`  cat`, `dog  `, s1, ` w w `, 6.8).Trim()
	checkEqual(t, []string{`cat`, `dog`, `plop`, `w w`, `6.8`}, e1.ToSlice())
}

func Test_Enumerator_Join(t *testing.T) {
	s1 := &pseudoStringer{text: `plop`}
	e1 := Enumerate[any](`cat`, `dog`, s1, ` w w `, 6.8)
	checkEqual(t, `cat|dog|plop| w w |6.8`, e1.Join(`|`))
}

func Test_Enumerator_Append(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`).Append(`wolf`, `hat`, `mouse`).Append(`dog`).Append()
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`}, e1.ToSlice())
}

func Test_Enumerator_Concat(t *testing.T) {
	e1 := Enumerate(`cat`, `bat`)
	e2 := Enumerate(`wolf`, `hat`, `mouse`)
	e3 := Enumerate(`dog`)
	e4 := e1.Concat(e2, e3)
	checkEqual(t, []string{`cat`, `bat`, `wolf`, `hat`, `mouse`, `dog`}, e4.ToSlice())
}

func Test_Enumerator_OfType(t *testing.T) {
	e := Enumerate[any](1, 2.0, int64(3), uint(4), `five`, fmt.Errorf(`six`), true)
	checkEqual(t, []bool{true}, OfType[bool](e).ToSlice())
	checkEqual(t, []int{1}, OfType[int](e).ToSlice())
	checkEqual(t, []string{`five`}, OfType[string](e).ToSlice())
	checkEqual(t, []float64{2.0}, OfType[float64](e).ToSlice())
	checkEqual(t, []any{1, 2.0, int64(3), uint(4), `five`, fmt.Errorf(`six`), true}, OfType[any](e).ToSlice())
}

func Test_Enumerator_Cast(t *testing.T) {
	e := Enumerate[any](1, 2.0, int64(3), uint(4), `five`, fmt.Errorf(`six`), true)
	checkEqual(t, []bool{false, false, false, false, false, false, true}, Cast[bool](e).ToSlice())
	checkEqual(t, []int{1, 2, 3, 4, 0, 0, 0}, Cast[int](e).ToSlice())
	checkEqual(t, []string{`"\x01"`, `""`, `"\x03"`, `"\x04"`, `"five"`, `""`, `""`}, Cast[string](e).Quotes().ToSlice())
	checkEqual(t, []float64{1.0, 2.0, 3.0, 4.0, 0.0, 0.0, 0.0}, Cast[float64](e).ToSlice())
}

func Test_Enumerator_Buffered(t *testing.T) {
	iterCreated := 0
	e := New(func() collections.Iterator[int] {
		iterCreated++
		return iterator.Range(1, 5)
	}).Buffered()

	checkEqual(t, 0, iterCreated)
	checkEqual(t, []int{1, 2, 3, 4, 5}, e.ToSlice())
	checkEqual(t, 1, iterCreated)
	checkEqual(t, []int{1, 2, 3, 4, 5}, e.ToSlice())
	checkEqual(t, 1, iterCreated)
	checkEqual(t, []int{1, 2, 3, 4, 5}, e.ToSlice())
	checkEqual(t, 1, iterCreated)
}

func Test_Enumerator_StartsWith(t *testing.T) {
	e1 := Enumerate(1, 2, 3, 4, 5)
	e2 := Enumerate(1, 2, 3, 4, 5)
	checkEqual(t, true, e1.StartsWith(e2))

	e1 = Enumerate(1, 2, 3, 4, 5)
	e2 = Enumerate(1, 2, 3)
	checkEqual(t, true, e1.StartsWith(e2))

	e1 = Enumerate(1, 2, 3)
	e2 = Enumerate(1, 2, 3, 4, 5)
	checkEqual(t, false, e1.StartsWith(e2))

	e1 = Enumerate(1, 8, 3, 4, 5)
	e2 = Enumerate(1, 2, 3)
	checkEqual(t, false, e1.StartsWith(e2))
}

func Test_Enumerator_Equal(t *testing.T) {
	e1 := Enumerate(1, 2, 3)
	e2 := Enumerate(1, 2, 3)
	checkEqual(t, e1, e2)

	e2 = Enumerate(1, 2, 3, 4)
	checkNotEqual(t, e1, e2)

	e2 = Enumerate(1, 2)
	checkNotEqual(t, e1, e2)

	e2 = Enumerate(1, 2, 4)
	checkNotEqual(t, e1, e2)

	e2 = Enumerate(1, 2, 4)
	checkNotEqual(t, any(5), e2)
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

func checkNotEqual(t testing.TB, exp, actual any) {
	t.Helper()
	if utils.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value shouldn't have matched the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
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
