package comp

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

func checkComparer[T any](t *testing.T, cmp Comparer[T], x, y T, exp int) {
	if actual := cmp(x, y); exp != actual {
		t.Errorf("\n"+
			"Unexpected value from Comparer:\n"+
			"\tKey Type:    %T\n"+
			"\tLeft Value:  %v\n"+
			"\tRight Value: %v\n"+
			"\tActual:      %d\n"+
			"\tExpected:    %d\n", x, x, y, actual, exp)
	}
}

func Test_Comp_OrderedComparer(t *testing.T) {
	c := Ordered[string]()
	checkComparer(t, c, `banana`, `cat`, -1)
	checkComparer(t, c, `cat`, `banana`, 1)
	checkComparer(t, c, `banana`, `banana`, 0)
	checkComparer(t, c, `cat`, `cat`, 0)
}

type pseudoComparable struct {
	name string
}

func (c *pseudoComparable) CompareTo(other *pseudoComparable) int {
	if c == nil {
		if other == nil {
			return 0
		}
		return -1
	}
	if other == nil {
		return 1
	}
	return strings.Compare(c.name, other.name)
}

func Test_Comp_Comparable(t *testing.T) {
	c := ComparableComparer[*pseudoComparable]()
	pc0 := (*pseudoComparable)(nil)
	pc1 := &pseudoComparable{name: `banana`}
	pc2 := &pseudoComparable{name: `cat`}

	checkComparer(t, c, pc0, pc0, 0)
	checkComparer(t, c, pc0, pc1, -1)
	checkComparer(t, c, pc0, pc2, -1)

	checkComparer(t, c, pc1, pc0, 1)
	checkComparer(t, c, pc1, pc1, 0)
	checkComparer(t, c, pc1, pc2, -1)

	checkComparer(t, c, pc2, pc0, 1)
	checkComparer(t, c, pc2, pc1, 1)
	checkComparer(t, c, pc2, pc2, 0)
}

func Test_Comp_FromLess(t *testing.T) {
	cmp := FromLess(func(x, y string) bool {
		return len(x) < len(y)
	})

	values := []string{`cat`, `dogs`, `doggo`, `apple`, `ox`}
	exp := `ox, cat, dogs, doggo, apple`
	c := slices.Clone(values)
	slices.SortFunc(c, cmp)
	if result := strings.Join(c, `, `); result != exp {
		t.Errorf("\n"+
			"Unexpected value from FromLess sort:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", result, exp)
	}
}

func Test_Comp_Epsilon(t *testing.T) {
	cmp := Epsilon(0.01)
	checkComparer(t, cmp, 0.0, 0.0, 0)
	checkComparer(t, cmp, 1.0, 1.0, 0)
	checkComparer(t, cmp, -1.0, -1.0, 0)

	checkComparer(t, cmp, 0.0, 0.005, 0)
	checkComparer(t, cmp, 0.0, 0.01, 0)
	checkComparer(t, cmp, 0.0, 0.02, -1)
	checkComparer(t, cmp, 0.0, 1.0, -1)
	checkComparer(t, cmp, 0.0, -0.005, 0)
	checkComparer(t, cmp, 0.0, -0.01, 0)
	checkComparer(t, cmp, 0.0, -0.02, 1)
	checkComparer(t, cmp, 0.0, -1.0, 1)

	checkComparer(t, cmp, 0.005, 0.0, 0)
	checkComparer(t, cmp, 0.01, 0.0, 0)
	checkComparer(t, cmp, 0.02, 0.0, 1)
	checkComparer(t, cmp, 1.0, 0.0, 1)
	checkComparer(t, cmp, -0.005, 0.0, 0)
	checkComparer(t, cmp, -0.01, 0.0, 0)
	checkComparer(t, cmp, -0.02, 0.0, -1)
	checkComparer(t, cmp, -1.0, 0.0, -1)

	cmp = Epsilon(-1.0) // defaults to ordered comparer, epsilon = 0
	checkComparer(t, cmp, 0.0, 0.0, 0)
	checkComparer(t, cmp, 1.0, 1.0, 0)
	checkComparer(t, cmp, -1.0, -1.0, 0)
}

func Test_Comp_DefaultComparer(t *testing.T) {
	checkComparer(t, Default[int](), 1, 3, -1)
	checkComparer(t, Default[int8](), 1, 3, -1)
	checkComparer(t, Default[int16](), 1, 3, -1)
	checkComparer(t, Default[int32](), 1, 3, -1)
	checkComparer(t, Default[int64](), 1, 3, -1)

	checkComparer(t, Default[uint](), 1, 3, -1)
	checkComparer(t, Default[uint8](), 1, 3, -1)
	checkComparer(t, Default[uint16](), 1, 3, -1)
	checkComparer(t, Default[uint32](), 1, 3, -1)
	checkComparer(t, Default[uint64](), 1, 3, -1)

	checkComparer(t, Default[float32](), 1.0, 3.0, -1)
	checkComparer(t, Default[float64](), 1.0, 3.0, -1)

	checkComparer(t, Default[uintptr](), 1, 3, -1)
	checkComparer(t, Default[string](), `apple`, `dog`, -1)
	checkComparer(t, Default[rune](), 'A', 'B', -1)
	checkComparer(t, Default[byte](), 1, 3, -1)

	cc := Default[*pseudoComparable]()
	pc0 := (*pseudoComparable)(nil)
	pc1 := &pseudoComparable{name: `apple`}
	pc2 := &pseudoComparable{name: `dog`}
	checkComparer(t, cc, pc0, pc0, 0)
	checkComparer(t, cc, pc0, pc1, -1)
	checkComparer(t, cc, pc0, pc2, -1)

	checkComparer(t, cc, pc1, pc0, 1)
	checkComparer(t, cc, pc1, pc1, 0)
	checkComparer(t, cc, pc1, pc2, -1)

	checkComparer(t, cc, pc2, pc0, 1)
	checkComparer(t, cc, pc2, pc1, 1)
	checkComparer(t, cc, pc2, pc2, 0)

	if d := Default[[]string](); !liteUtils.IsNil(d) {
		t.Errorf(`expected a nil value from Default with a slice but got %v`, d)
	}

	checkComparer(t, Default[time.Duration](), time.Second, time.Minute, -1)
	checkComparer(t, Default[time.Duration](), time.Hour, time.Minute, 1)
	checkComparer(t, Default[time.Duration](), time.Second, time.Second, 0)
	checkComparer(t, Default[time.Duration](), time.Hour, time.Hour, 0)

	time1, err := time.Parse(time.RFC822Z, `02 Jan 24 05:30 -0700`)
	if err != nil {
		panic(err)
	}
	time2, err := time.Parse(time.RFC822Z, `02 Jan 24 05:35 -0700`)
	if err != nil {
		panic(err)
	}
	checkComparer(t, Default[time.Time](), time1, time2, -1)
	checkComparer(t, Default[time.Time](), time1, time1, 0)
	checkComparer(t, Default[time.Time](), time2, time2, 0)
	checkComparer(t, Default[time.Time](), time2, time1, 1)
}

type pseudoEquatable struct{ success bool }

func (pe *pseudoEquatable) Equals(_ any) bool { return pe.success }

func checkEqual(t *testing.T, a, b any, exp bool) {
	t.Helper()
	if Equal(a, b) != exp {
		t.Errorf("\n"+
			"Unexpected value from Equal:\n"+
			"\tValue 1:  %v (%T)\n"+
			"\tValue 2:  %v (%T)\n"+
			"\tExpected: %t\n", a, a, b, b, exp)
	}
}

func Test_Utils_Equal(t *testing.T) {
	checkEqual(t, true, true, true)
	checkEqual(t, false, true, false)
	checkEqual(t, true, false, false)
	checkEqual(t, false, false, true)

	checkEqual(t, 1, 1, true)
	checkEqual(t, 1, 2, false)
	checkEqual(t, 2, 1, false)
	checkEqual(t, 2, 2, true)

	e0 := (error)(nil)
	e1 := errors.New(`fred`)
	checkEqual(t, nil, nil, true)
	checkEqual(t, e0, e0, true)
	checkEqual(t, e1, e0, false)
	checkEqual(t, e0, e1, false)
	checkEqual(t, e1, e1, true)

	var v1 int = 0
	var v2 float64 = 0.0
	checkEqual(t, nil, v1, false)
	checkEqual(t, v2, v1, false)
	checkEqual(t, v1, v2, false)

	e2 := &pseudoEquatable{success: false}
	checkEqual(t, e2, nil, false)
	checkEqual(t, nil, e2, false)
	checkEqual(t, e2, e2, false)
	checkEqual(t, 4, e2, false)
	checkEqual(t, e2, 4, false)

	e3 := &pseudoEquatable{success: true}
	checkEqual(t, e3, nil, false)
	checkEqual(t, nil, e3, false)
	checkEqual(t, e3, e3, true)
	checkEqual(t, 4, e3, true)
	checkEqual(t, e3, 4, true)

	e4 := (*pseudoEquatable)(nil)
	checkEqual(t, e4, nil, false)
	checkEqual(t, nil, e4, false)
	checkEqual(t, e4, e0, false)
	checkEqual(t, e0, e4, false)
	checkEqual(t, e4, e4, true)

	e5 := ([]int)(nil)
	e6 := []int{}
	e7 := []int{1, 2, 3}
	e8 := []int{1, 4, 3}
	checkEqual(t, e5, nil, false)
	checkEqual(t, nil, e5, false)
	checkEqual(t, e5, e5, true)
	checkEqual(t, e5, e6, false)
	checkEqual(t, e5, e7, false)
	checkEqual(t, e5, e8, false)
	checkEqual(t, e6, e6, true)
	checkEqual(t, e6, e7, false)
	checkEqual(t, e6, e8, false)
	checkEqual(t, e7, e7, true)
	checkEqual(t, e7, e8, false)
	checkEqual(t, e8, e8, true)

	e9 := func() { print(`boom`) }
	var e10 func()
	checkEqual(t, e9, e9, false)
	checkEqual(t, e9, e10, false)
	checkEqual(t, e10, e9, false)
	checkEqual(t, e10, e10, true)
	checkEqual(t, e10, nil, false)
	checkEqual(t, nil, e10, false)
}

func checkSort[T any, S ~[]T](t *testing.T, cmp Comparer[T], values S, exp string) {
	t.Helper()
	slices.SortFunc(values, cmp)
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = fmt.Sprint(v)
	}
	if result := strings.Join(strs, `, `); exp != result {
		t.Errorf("\n"+
			"Unexpected value from SortedKeys:\n"+
			"\tKey Type: %T\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", liteUtils.Zero[T](), result, exp)
	}
}

func Test_Comp_Sorting(t *testing.T) {
	checkSort(t, Default[int](),
		[]int{3, 1, 6, 5},
		`1, 3, 5, 6`)

	checkSort(t, Default[float64](),
		[]float64{4.3, 2.16, 333.333, 12.34},
		`2.16, 4.3, 12.34, 333.333`)

	checkSort(t, Default[rune](),
		[]rune{'k', 'a', 'q', 'H'},
		`72, 97, 107, 113`)

	checkSort(t, Default[string](),
		[]string{`cat`, `pig`, `Dog`, `apple`},
		`Dog, apple, cat, pig`)

	checkSort(t, Descender(Ordered[int]()),
		[]int{3, 1, 6, 5},
		`6, 5, 3, 1`)

	checkSort(t, Descender(Ordered[float64]()),
		[]float64{4.3, 2.16, 333.333, 12.34},
		`333.333, 12.34, 4.3, 2.16`)

	checkSort(t, Descender(Ordered[rune]()),
		[]rune{'k', 'a', 'q', 'H'},
		`113, 107, 97, 72`)

	checkSort(t, Descender(Ordered[string]()),
		[]string{`cat`, `pig`, `Dog`, `apple`},
		`pig, cat, apple, Dog`)

	checkSort(t, Default[bool](),
		[]bool{false, true, false, true, true, false, false},
		`false, false, false, false, true, true, true`)
}

func Test_Comp_Slice(t *testing.T) {
	cmp := Slice[[]int](Default[int]())
	check := func(exp int, a, b []int) {
		t.Helper()
		checkEqual(t, cmp(a, b), exp, true)
	}

	check(0, []int{}, []int{})
	check(0, []int{1, 2, 3}, []int{1, 2, 3})
	check(-1, []int{1, 2, 3}, []int{1, 2, 3, 4})
	check(1, []int{1, 2, 3, 4}, []int{1, 2, 3})
	check(1, []int{4}, []int{1, 2, 3})
	check(-1, []int{1, 2, 3}, []int{4})
}

func Test_Comp_Or(t *testing.T) {
	type person struct {
		first string
		last  string
		age   int
	}

	people := []person{
		{first: `Tim`, last: `Smith`, age: 13},
		{first: `Jen`, last: `Smith`, age: 18},
		{first: `Bob`, last: `Smith`, age: 54},
		{first: `Bob`, last: `Hicks`, age: 23},
		{first: `Kim`, last: `Hicks`, age: 25},
		{first: `Bob`, last: `Smith`, age: 23},
		{first: `Jen`, last: `Hicks`, age: 21},
		{first: `Kim`, last: `Hicks`, age: 25},
		{first: `Sal`, last: `Atoms`, age: 13},
	}

	check := func(c Comparer[person], exp ...string) {
		t.Helper()
		ps := slices.Clone(people)
		slices.SortFunc(ps, c)
		str := make([]string, len(ps))
		for i, p := range ps {
			str[i] = fmt.Sprintf(`%s %s %d`, p.first, p.last, p.age)
		}
		if !Equal(exp, str) {
			t.Error("unexpected result from Or comparison sort:\n" +
				fmt.Sprintf("\tgot: [%s]\n", strings.Join(str, ",\n\t\t")) +
				fmt.Sprintf("\texp: [%s]\n", strings.Join(exp, ",\n\t\t")))
		}
	}

	check(func(a, b person) int {
		return Or(
			func() int { return Default[string]()(a.first, b.first) },
			func() int { return Default[string]()(a.last, b.last) },
			func() int { return Default[int]()(a.age, b.age) },
		)
	}, `Bob Hicks 23`,
		`Bob Smith 23`,
		`Bob Smith 54`,
		`Jen Hicks 21`,
		`Jen Smith 18`,
		`Kim Hicks 25`,
		`Kim Hicks 25`,
		`Sal Atoms 13`,
		`Tim Smith 13`)

	check(func(a, b person) int {
		return Or(
			func() int { return Default[string]()(a.last, b.last) },
			func() int { return Default[string]()(a.first, b.first) },
			func() int { return Default[int]()(a.age, b.age) },
		)
	}, `Sal Atoms 13`,
		`Bob Hicks 23`,
		`Jen Hicks 21`,
		`Kim Hicks 25`,
		`Kim Hicks 25`,
		`Bob Smith 23`,
		`Bob Smith 54`,
		`Jen Smith 18`,
		`Tim Smith 13`)

	check(func(a, b person) int {
		return Or(
			func() int { return Default[int]()(a.age, b.age) },
			func() int { return Default[string]()(a.first, b.first) },
			func() int { return Default[string]()(a.last, b.last) },
		)
	}, `Sal Atoms 13`,
		`Tim Smith 13`,
		`Jen Smith 18`,
		`Jen Hicks 21`,
		`Bob Hicks 23`,
		`Bob Smith 23`,
		`Kim Hicks 25`,
		`Kim Hicks 25`,
		`Bob Smith 54`)
}
