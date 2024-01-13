package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/differs/data"
	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_Check_Require(t *testing.T) {
	pt := newTester(t)
	fVal := 12.0
	Equal(pt, 12.0).Require(fVal)
	pt.Check()

	Equal(pt, `World`).Require(`Hello`)
	pt.Check(`Must be equal:`,
		`\tActual Type:    string`,
		`\tActual Value:   "Hello"`,
		`\tExpected Type:  string`,
		`\tExpected Value: "World"`,
		`FAIL NOW`)

	Equal(pt, (error)(nil)).Require(nil)
	pt.Check()
}

func Test_Check_Required_Assert(t *testing.T) {
	pt := newTester(t)
	fVal := 12.0
	Equal(pt, 12.0).Required().Assert(fVal)
	pt.Check()

	Equal(pt, `World`).Required().Assert(`Hello`)
	pt.Check(`Must be equal:`,
		`\tActual Type:    string`,
		`\tActual Value:   "Hello"`,
		`\tExpected Type:  string`,
		`\tExpected Value: "World"`,
		`FAIL NOW`)

	Equal(pt, (error)(nil)).Required().Assert(nil)
	pt.Check()
}

func Test_Check_AssertAll(t *testing.T) {
	pt := newTester(t)
	GreaterThan(pt, 5).AssertAll([]int{6, 7, 8, 9})
	pt.Check()

	GreaterThan(pt, 5).AssertAll([]int{4, 6, 7, 3})
	pt.Check(`Should be greater than the expected value:`,
		`\tActual Collection: \[4 6 7 3\]`,
		`\tActual Type:       int`,
		`\tActual Value:      4`,
		`\tMinimum Type:      int`,
		`\tMinimum Value:     5`,
		`\tValue Index:       0`,
		``,
		`Should be greater than the expected value:`,
		`\tActual Collection: \[4 6 7 3\]`,
		`\tActual Type:       int`,
		`\tActual Value:      3`,
		`\tMinimum Type:      int`,
		`\tMinimum Value:     5`,
		`\tValue Index:       3`)

	Match(pt, `[]`).AssertAll([]string{`Hello`, `World`})
	pt.Check(`Must provide a valid regular expression pattern:`,
		`\tPattern: \[\]`,
		`FAIL NOW`)

	Match(pt, `^\w+$`).AssertAll([]string{`Hello`, `World`, `Pancakes!`})
	pt.Check(`Should match the given regular expression pattern:`,
		`\tActual Collection: \[Hello World Pancakes!\]`,
		`\tActual Type:       string`,
		`\tActual Value:      "Pancakes!"`,
		`\tPattern:           \^\\w\+\$`,
		`\tValue Index:       2`)
}

func Test_Check_RequireAll(t *testing.T) {
	pt := newTester(t)
	GreaterThan(pt, 5).RequireAll([]int{6, 7, 8, 9})
	pt.Check()

	GreaterThan(pt, 5).RequireAll([]int{4, 6, 7, 3})
	pt.Check(`Must be greater than the expected value:`,
		`\tActual Collection: \[4 6 7 3\]`,
		`\tActual Type:       int`,
		`\tActual Value:      4`,
		`\tMinimum Type:      int`,
		`\tMinimum Value:     5`,
		`\tValue Index:       0`,
		`FAIL NOW`,
		// If actually running the following wouldn't
		// be shown because of the FAIL NOW
		``,
		`Must be greater than the expected value:`,
		`\tActual Collection: \[4 6 7 3\]`,
		`\tActual Type:       int`,
		`\tActual Value:      3`,
		`\tMinimum Type:      int`,
		`\tMinimum Value:     5`,
		`\tValue Index:       3`,
		`FAIL NOW`)
}

func Test_Check_Panic(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 12.0).Panic(func() {
		panic(12.0)
	})
	pt.Check()

	Equal(pt, 12.0).Panic(func() {
		panic(12.25)
	})
	pt.Check(`Should be equal:`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 12`,
		`\tPanicked Type:  float64`,
		`\tPanicked Value: 12\.25`)

	Equal(pt, 12.0).Panic(func() {
		panic(errors.New(`you smell something?`))
	})
	pt.Check(`Should be a panic of the expected type:`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 12`,
		`\tPanicked Type:  \*errors\.errorString`,
		`\tPanicked Value: you smell something\?`)

	Equal(pt, 12.0).Panic(func() {})
	pt.Check(`Should panic from given function:`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 12`)

	MatchError(pt, `({`).Panic(func() {}) // nil check
	pt.Check(`Must provide a valid regular expression pattern:`,
		`\tPattern: \(\{`,
		`FAIL NOW`)
}

func Test_Check_With(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 13).Name(`Pickle`).
		With(`Cakes`, "- Chocolate\n", "- Lemon\n", `- Vanilla`).
		Assert(12)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tCakes:          - Chocolate`,
		`\t                - Lemon`,
		`\t                - Vanilla`,
		`\tExpected Type:  int`,
		`\tExpected Value: 13`,
		`\tName:           Pickle`)
}

func Test_Check_Withf(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 13).Name(`Pickle`).
		Withf(`Cakes`, "- Chocolate\n- %s\n- Vanilla", `Lemon`).
		Assert(12)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tCakes:          - Chocolate`,
		`\t                - Lemon`,
		`\t                - Vanilla`,
		`\tExpected Type:  int`,
		`\tExpected Value: 13`,
		`\tName:           Pickle`)
}

func Test_Check_MaxLines(t *testing.T) {
	pt := newTester(t)
	lines := make([]string, 200)
	for i := range lines {
		lines[i] = strconv.Itoa(i)
	}
	Equal(pt, 13).
		With(`Long Output`, strings.Join(lines, "\n")).
		Assert(12)

	actual := pt.buf.String()
	expPattern := `^\n` +
		`Should be equal:\n` +
		`\tActual Type:    int\n` +
		`\tActual Value:   12\n` +
		`\tExpected Type:  int\n` +
		`\tExpected Value: 13\n` +
		`\tLong Output:    0\n` +
		`\t                1\n` +
		`\t                2\n` +
		`(?:\t                \d+\n)+` +
		`\t                93\n` +
		`\t                94\n` +
		`\t                95\n` +
		`\t                ...\n` +
		`\t                197\n` +
		`\t                198\n` +
		`\t                199\n$`
	matched, err := regexp.MatchString(expPattern, actual)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("failed to match expected regular expression:\n%q\n", actual)
	}
}

func Test_Check_TextHint(t *testing.T) {
	pt := newTester(t)

	Equal(pt, []rune(`Hello`)).Assert([]rune(`Help`))
	pt.Check(`Should be equal:`,
		`\tActual Type:    \[\]int32`,
		`\tActual Value:   \[72 101 108 112\]`,
		`\tExpected Type:  \[\]int32`,
		`\tExpected Value: \[72 101 108 108 111\]`)
	Equal(pt, []rune(`Hello`)).AsText().Assert([]rune(`Help`))
	pt.Check(`Should be equal:`,
		`\tActual Type:    \[\]int32`,
		`\tActual Value:   "Help"`,
		`\tExpected Type:  \[\]int32`,
		`\tExpected Value: "Hello"`)

	Equal(pt, []byte(`Hello`)).Assert([]byte(`Help`))
	pt.Check(`Should be equal:`,
		`\tActual Type:    \[\]uint8`,
		`\tActual Value:   \[72 101 108 112\]`,
		`\tExpected Type:  \[\]uint8`,
		`\tExpected Value: \[72 101 108 108 111\]`)
	Equal(pt, []byte(`Hello`)).AsText().Assert([]byte(`Help`))
	pt.Check(`Should be equal:`,
		`\tActual Type:    \[\]uint8`,
		`\tActual Value:   "Help"`,
		`\tExpected Type:  \[\]uint8`,
		`\tExpected Value: "Hello"`)

	Equal(pt, 'H').Assert('Q')
	pt.Check(`Should be equal:`,
		`\tActual Type:    int32`,
		`\tActual Value:   81`,
		`\tExpected Type:  int32`,
		`\tExpected Value: 72`)
	Equal(pt, 'H').AsText().Assert('Q')
	pt.Check(`Should be equal:`,
		`\tActual Type:    int32`,
		`\tActual Value:   'Q'`,
		`\tExpected Type:  int32`,
		`\tExpected Value: 'H'`)

	// Currently formatting doesn't dig down into slices or maps.
	Equal(pt, []any{13, `Hi`, 5.3}).AsBin().AsText().Assert([]any{42, `Bye`, 3.4})
	pt.Check(`Should be equal:`,
		`\tActual Type:    \[\]interface \{\}`,
		`\tActual Value:   \[42 Bye 3\.4\]`,
		`\tExpected Type:  \[\]interface \{\}`,
		`\tExpected Value: \[13 Hi 5\.3\]`)
}

func Test_Check_HexHint(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 0xDEAD).Assert(-0xF00D)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   -61453`,
		`\tExpected Type:  int`,
		`\tExpected Value: 57005`)
	Equal(pt, 0xDEAD).AsHex().Assert(-0xF00D)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   -0xF00D`,
		`\tExpected Type:  int`,
		`\tExpected Value: 0xDEAD`)
	Equal[uint64](pt, 0xFF66_55AA).AsHex().Assert(0xFEDCBA98_76543210)
	pt.Check(`Should be equal:`,
		`\tActual Type:    uint64`,
		`\tActual Value:   0xFEDC_BA98_7654_3210`,
		`\tExpected Type:  uint64`,
		`\tExpected Value: 0xFF66_55AA`)
	Equal(pt, `ABBA`).AsHex().Assert(`ACDC`) // no effect on strings
	pt.Check(`Should be equal:`,
		`\tActual Type:    string`,
		`\tActual Value:   "ACDC"`,
		`\tExpected Type:  string`,
		`\tExpected Value: "ABBA"`)
}

func Test_Check_OctHint(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 0o555).Assert(-0o234)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   -156`,
		`\tExpected Type:  int`,
		`\tExpected Value: 365`)
	Equal(pt, 0o555).AsOct().Assert(-0o234)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   -0o234`,
		`\tExpected Type:  int`,
		`\tExpected Value: 0o555`)
	Equal[uint64](pt, 0o76543210).AsOct().Assert(0o77_665544_33221100)
	pt.Check(`Should be equal:`,
		`\tActual Type:    uint64`,
		`\tActual Value:   0o7766_5544_3322_1100`,
		`\tExpected Type:  uint64`,
		`\tExpected Value: 0o7654_3210`)
}

func Test_Check_BinHint(t *testing.T) {
	pt := newTester(t)
	Equal(pt, 0xDEAD).AsBin().Assert(-0xF00D)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   -0b1111_0000_0000_1101`,
		`\tExpected Type:  int`,
		`\tExpected Value: 0b1101_1110_1010_1101`)
	Equal[uint64](pt, 0xFF66_55AA).AsBin().Assert(0xFEDCBA98_76543210)
	pt.Check(`Should be equal:`,
		`\tActual Type:    uint64`,
		`\tActual Value:   0b1111_1110_1101_1100_1011_1010_1001_1000_0111_`+
			`0110_0101_0100_0011_0010_0001_0000`,
		`\tExpected Type:  uint64`,
		`\tExpected Value: 0b1111_1111_0110_0110_0101_0101_1010_1010`)
}

func Test_Check_TimeHint(t *testing.T) {
	pt := newTester(t)
	time1 := time.Date(2023, time.August, 26, 9, 30, 0, 0, time.UTC)
	time2 := time.Date(2024, time.February, 12, 9, 30, 0, 0, time.UTC)

	Equal(pt, time1).Assert(time2)
	pt.Check(`Should be equal:`,
		`\tActual Type:    time\.Time`,
		`\tActual Value:   2024-02-12 09:30:00 \+0000 UTC`,
		`\tExpected Type:  time\.Time`,
		`\tExpected Value: 2023-08-26 09:30:00 \+0000 UTC`)

	Equal(pt, time1).TimeAs(time.DateOnly).Assert(time2)
	pt.Check(`Should be equal:`,
		`\tActual Type:    time\.Time`,
		`\tActual Value:   2024-02-12`,
		`\tExpected Type:  time\.Time`,
		`\tExpected Value: 2023-08-26`)

	// Using the wrong formatting may make error hard to understand
	Equal(pt, time1).TimeAs(time.TimeOnly).Assert(time2)
	pt.Check(`Should be equal:`,
		`\tActual Type:    time\.Time`,
		`\tActual Value:   09:30:00`,
		`\tExpected Type:  time\.Time`,
		`\tExpected Value: 09:30:00`)
}

func Test_Check_Helper(t *testing.T) {
	pt := &pseudoTesterWithHelper{
		pseudoTester: *newTester(t),
	}
	c := Equal(pt, 13)
	pt.Check(`Helper`)
	c.Assert(12)
	pt.Check(`Helper`,
		``,
		`Helper`,
		``,
		`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tExpected Type:  int`,
		`\tExpected Value: 13`)
}

func Test_Check_Nil(t *testing.T) {
	pt := newTester(t)
	i := 0
	Nil(pt).Assert(i)
	pt.Check(`Should be a nil-able type:`,
		`\tActual Type:  int`,
		`\tActual Value: 0`)

	pi := &i
	Nil(pt).Assert(pi)
	pt.Check(`Should be nil:`,
		`\tActual Type:  \*int`,
		`\tActual Value: 0x[0-9a-f]+`)

	pi = nil
	Nil(pt).Assert(pi)
	pt.Check()
}

func Test_Check_NotNil(t *testing.T) {
	pt := newTester(t)
	i := 0
	NotNil(pt).Assert(i)
	pt.Check(`Should be a nil-able type:`,
		`\tActual Type:  int`,
		`\tActual Value: 0`)

	pi := &i
	NotNil(pt).Assert(pi)
	pt.Check()

	pi = nil
	NotNil(pt).Assert(pi)
	pt.Check(`Should not be nil:`,
		`\tActual Type:  \*int`,
		`\tActual Value: <nil>`)
}

func Test_Check_Zero(t *testing.T) {
	pt := newTester(t)
	Zero(pt).Assert(12)
	pt.Check(`Should be a zero value:`,
		`\tActual Type:  int`,
		`\tActual Value: 12`)

	Zero(pt).Assert(0)
	Zero(pt).Assert(``)
	Zero(pt).Assert(0.0)
	Zero(pt).Assert(nil)
	Zero(pt).Assert(struct{ a int }{a: 0})
	pt.Check()
}

func Test_Check_NotZero(t *testing.T) {
	pt := newTester(t)
	NotZero(pt).Assert(0)
	pt.Check(`Should not be a zero value:`,
		`\tActual Type:  int`,
		`\tActual Value: 0`)

	v := 12
	NotZero(pt).Assert(32)
	NotZero(pt).Assert(`Zero`)
	NotZero(pt).Assert(0.001)
	NotZero(pt).Assert(&v)
	NotZero(pt).Assert(struct{ a int }{a: 56})
	pt.Check()
}

func Test_Check_True(t *testing.T) {
	pt := newTester(t)
	True(pt).Assert(false)
	pt.Check(`Should be true:`,
		`\tActual Type:  bool`,
		`\tActual Value: false`)

	True(pt).Assert(true)
	pt.Check()
}

func Test_Check_False(t *testing.T) {
	pt := newTester(t)
	False(pt).Assert(false)
	pt.Check()

	False(pt).Assert(true)
	pt.Check(`Should be false:`,
		`\tActual Type:  bool`,
		`\tActual Value: true`)
}

func Test_Check_Type(t *testing.T) {
	pt := newTester(t)
	v := 12
	Type[float64](pt).Assert(v)
	pt.Check(`Should be the expected type:`,
		`\tActual Type:   int`,
		`\tActual Value:  12`,
		`\tExpected Type: float64`)

	Type[int](pt).Assert(v)
	pt.Check()

	Type[*int](pt).Assert(&v)
	pt.Check()
}

func Test_Check_NotType(t *testing.T) {
	pt := newTester(t)
	v := 12
	NotType[float64](pt).Assert(v)
	pt.Check()

	NotType[int](pt).Assert(v)
	pt.Check(`Should not be the unexpected type:`,
		`\tActual Type:     int`,
		`\tActual Value:    12`,
		`\tUnexpected Type: int`)

	NotType[*int](pt).Assert(&v)
	pt.Check(`Should not be the unexpected type:`,
		`\tActual Type:     \*int`,
		`\tActual Value:    0x[0-9a-f]+`,
		`\tUnexpected Type: \*int`)
}

func Test_Check_Match(t *testing.T) {
	pt := newTester(t)
	Match(pt, ``).Assert(`Cat`)
	pt.Check(`Must provide a non-empty regular expression pattern`,
		`FAIL NOW`)

	Match(pt, `(()!`).Assert(`Cat`)
	pt.Check(`Must provide a valid regular expression pattern:`,
		`\tPattern: \(\(\)!`,
		`FAIL NOW`)

	Match(pt, `0x[0-9A-Z]{4}`).Assert(`Cat`)
	pt.Check(`Should match the given regular expression pattern:`,
		`\tActual Type:  string`,
		`\tActual Value: "Cat"`,
		`\tPattern:      0x\[0-9A-Z\]\{4\}`)

	Match(pt, `0x[0-9A-Z]{4}`).Assert(`0xF4B2`)
	pt.Check()

	reg := regexp.MustCompile(`^\w+$`)
	Match(pt, reg).Assert(`World`)
	pt.Check()

	Match(pt, reg).Assert(`Hello World`)
	pt.Check(`Should match the given regular expression pattern:`,
		`\tActual Type:  string`,
		`\tActual Value: "Hello World"`,
		`\tPattern:      \^\\w\+\$`)

	reg = nil
	Match(pt, reg).Assert(`Cat`)
	pt.Check(`Must provide a non-nil regular expression instance`,
		`FAIL NOW`)

	Match(pt, 4).Assert(`Cat`)
	pt.Check(`Must provide a string pattern or \*regexp\.Regexp:`,
		`\tGiven Type: int`,
		`FAIL NOW`)
}

func Test_Check_String(t *testing.T) {
	pt := newTester(t)
	ps := pseudoStringer{text: `Lumpy`}
	String(pt, `cat`).Assert(ps)
	pt.Check(`Should have string be equal:`,
		`\tActual Type:     check.pseudoStringer`,
		`\tActual Value:    "Lumpy"`,
		`\tExpected String: "cat"`)

	String(pt, `Lumpy`).Assert(ps)
	pt.Check()
}

func Test_Check_Equal(t *testing.T) {
	pt := newTester(t)
	iVal := 12
	Equal(pt, 12).Assert(iVal)
	pt.Check()

	Equal(pt, 13).Assert(iVal)
	pt.Check(`Should be equal:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tExpected Type:  int`,
		`\tExpected Value: 13`)

	u64Val := uint64(12)
	Equal[uint64](pt, 12).Assert(u64Val)
	pt.Check()
}

func Test_Check_NotEqual(t *testing.T) {
	pt := newTester(t)
	NotEqual(pt, 12).Assert(12)
	pt.Check(`Should not be equal:`,
		`\tActual Type:      int`,
		`\tActual Value:     12`,
		`\tUnexpected Type:  int`,
		`\tUnexpected Value: 12`)

	NotEqual(pt, 13).Assert(12)
	pt.Check()
}

func Test_Check_GreaterThan(t *testing.T) {
	pt := newTester(t)
	GreaterThan(pt, 4).Assert(6)
	GreaterThan(pt, 4).Assert(5)
	pt.Check()

	GreaterThan(pt, 4).Assert(4)
	pt.Check(`Should be greater than the expected value:`,
		`\tActual Type:   int`,
		`\tActual Value:  4`,
		`\tMinimum Type:  int`,
		`\tMinimum Value: 4`)
}

func Test_Check_GreaterEq(t *testing.T) {
	pt := newTester(t)
	GreaterEq(pt, 4).Assert(5)
	GreaterEq(pt, 4).Assert(4)
	pt.Check()

	GreaterEq(pt, 4).Assert(3)
	pt.Check(`Should be greater than or equal to the expected value:`,
		`\tActual Type:   int`,
		`\tActual Value:  3`,
		`\tMinimum Type:  int`,
		`\tMinimum Value: 4`)
}

func Test_Check_LessThan(t *testing.T) {
	pt := newTester(t)
	LessThan(pt, 4).Assert(2)
	LessThan(pt, 4).Assert(3)
	pt.Check()

	LessThan(pt, 4).Assert(4)
	pt.Check(`Should be less than the expected value:`,
		`\tActual Type:   int`,
		`\tActual Value:  4`,
		`\tMaximum Type:  int`,
		`\tMaximum Value: 4`)
}

func Test_Check_LessEq(t *testing.T) {
	pt := newTester(t)
	LessEq(pt, 4).Assert(3)
	LessEq(pt, 4).Assert(4)
	pt.Check()

	LessEq(pt, 4).Assert(5)
	pt.Check(`Should be less than or equal to the expected value:`,
		`\tActual Type:   int`,
		`\tActual Value:  5`,
		`\tMaximum Type:  int`,
		`\tMaximum Value: 4`)
}

func Test_Check_InRange(t *testing.T) {
	pt := newTester(t)
	InRange(pt, 0, 10).Assert(0)
	InRange(pt, 0, 10).Assert(1)
	InRange(pt, 0, 10).Assert(9)
	InRange(pt, 0, 10).Assert(10)
	pt.Check()

	InRange(pt, 0, 10).Assert(11)
	pt.Check(`Should be between or equal to the given maximum and minimum:`,
		`\tActual Type:   int`,
		`\tActual Value:  11`,
		`\tMaximum Value: 10`,
		`\tMinimum Value: 0`,
		`\tRange Type:    int`)

	InRange(pt, 0, 10).Assert(-1)
	pt.Check(`Should be between or equal to the given maximum and minimum:`,
		`\tActual Type:   int`,
		`\tActual Value:  -1`,
		`\tMaximum Value: 10`,
		`\tMinimum Value: 0`,
		`\tRange Type:    int`)

	time1 := time.Date(2024, time.February, 12, 9, 10, 0, 0, time.UTC)
	time2 := time.Date(2024, time.February, 12, 9, 30, 0, 0, time.UTC)
	time3 := time.Date(2024, time.February, 12, 9, 40, 0, 0, time.UTC)
	InRange(pt, time1, time2).Assert(time3)
	pt.Check(`Should be between or equal to the given maximum and minimum:`,
		`\tActual Type:   time\.Time`,
		`\tActual Value:  2024-02-12 09:40:00 \+0000 UTC`,
		`\tMaximum Value: 2024-02-12 09:30:00 \+0000 UTC`,
		`\tMinimum Value: 2024-02-12 09:10:00 \+0000 UTC`,
		`\tRange Type:    time\.Time`)

	InRange(pt, time1, time3).Assert(time2)
	pt.Check()

	InRange(pt, time2, time3).Assert(time1)
	pt.Check(`Should be between or equal to the given maximum and minimum:`,
		`\tActual Type:   time\.Time`,
		`\tActual Value:  2024-02-12 09:10:00 \+0000 UTC`,
		`\tMaximum Value: 2024-02-12 09:40:00 \+0000 UTC`,
		`\tMinimum Value: 2024-02-12 09:30:00 \+0000 UTC`,
		`\tRange Type:    time\.Time`)
}

func Test_Check_Epsilon(t *testing.T) {
	pt := newTester(t)
	Epsilon(pt, 4.0, -1.0).Assert(3.995)
	pt.Check(`Must provide an epsilon greater than zero:`,
		`\tEpsilon Type:  float64`,
		`\tEpsilon Value: -1`,
		`FAIL NOW`)

	Epsilon(pt, 4.0, 0.01).Assert(3.995)
	pt.Check()

	Epsilon(pt, 4.0, 0.01).Assert(4.00)
	pt.Check()

	Epsilon(pt, 4.0, 0.01).Assert(4.005)
	pt.Check()

	Epsilon(pt, 4.0, 0.01).Assert(4.02)
	pt.Check(`Should be within the epsilon of the expected value:`,
		`\tActual Type:    float64`,
		`\tActual Value:   4\.02`,
		`\tDelta:          -0\.019999\d*`,
		`\tEpsilon:        0\.01`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 4`)

	Epsilon(pt, 4.02, 0.01).Assert(4.0)
	pt.Check(`Should be within the epsilon of the expected value:`,
		`\tActual Type:    float64`,
		`\tActual Value:   4`,
		`\tDelta:          0\.019999\d*`,
		`\tEpsilon:        0\.01`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 4\.02`)
}

func Test_Check_NotEpsilon(t *testing.T) {
	pt := newTester(t)
	NotEpsilon(pt, 4.0, -1.0).Assert(3.995)
	pt.Check(`Must provide an epsilon greater than zero:`,
		`\tEpsilon Type:  float64`,
		`\tEpsilon Value: -1`,
		`FAIL NOW`)

	NotEpsilon(pt, 4.0, 0.01).Assert(3.995)
	pt.Check(`Should not be within the epsilon of the unexpected value:`,
		`\tActual Type:      float64`,
		`\tActual Value:     3\.995`,
		`\tDelta:            0\.0049999\d+`,
		`\tEpsilon:          0\.01`,
		`\tExpected Type:    float64`,
		`\tUnexpected Value: 4`)

	NotEpsilon(pt, 4.0, 0.01).Assert(4.00)
	pt.Check(`Should not be within the epsilon of the unexpected value:`,
		`\tActual Type:      float64`,
		`\tActual Value:     4`,
		`\tDelta:            0`,
		`\tEpsilon:          0\.01`,
		`\tExpected Type:    float64`,
		`\tUnexpected Value: 4`)

	NotEpsilon(pt, 4.0, 0.01).Assert(4.005)
	pt.Check(`Should not be within the epsilon of the unexpected value:`,
		`\tActual Type:      float64`,
		`\tActual Value:     4\.005`,
		`\tDelta:            -0\.0049999\d+`,
		`\tEpsilon:          0\.01`,
		`\tExpected Type:    float64`,
		`\tUnexpected Value: 4`)

	NotEpsilon(pt, 4.0, 0.01).Assert(4.02)
	pt.Check()

	NotEpsilon(pt, 4.02, 0.01).Assert(4.0)
	pt.Check()
}

func Test_Check_Positive(t *testing.T) {
	pt := newTester(t)
	Positive[float64](pt).Assert(42.0)
	pt.Check()

	Positive[float64](pt).Assert(0.0)
	pt.Check(`Should be a positive value:`,
		`\tActual Type:   float64`,
		`\tActual Value:  0`,
		`\tExpected Type: float64`)

	Positive[float64](pt).Assert(-42.0)
	pt.Check(`Should be a positive value:`,
		`\tActual Type:   float64`,
		`\tActual Value:  -42`,
		`\tExpected Type: float64`)
}

func Test_Check_Negative(t *testing.T) {
	pt := newTester(t)
	Negative[float64](pt).Assert(42.0)
	pt.Check(`Should be a negative value:`,
		`\tActual Type:   float64`,
		`\tActual Value:  42`,
		`\tExpected Type: float64`)

	Negative[float64](pt).Assert(0.0)
	pt.Check(`Should be a negative value:`,
		`\tActual Type:   float64`,
		`\tActual Value:  0`,
		`\tExpected Type: float64`)

	Negative[float64](pt).Assert(-42.0)
	pt.Check()
}

func Test_Check_Is(t *testing.T) {
	pt := newTester(t)
	Is(pt, leapYear).Name(`Is leap year?`).Assert(1900)
	Is(pt, leapYear).Name(`Is leap year?`).Assert(1996)
	Is(pt, leapYear).Name(`Is leap year?`).Assert(2000)
	Is(pt, leapYear).Name(`Is leap year?`).Assert(2024)
	Is(pt, leapYear).Name(`Is leap year?`).Assert(2025)
	pt.Check(`Should be accepted by the given predicate:`,
		`\tActual Type:  int`,
		`\tActual Value: 1900`,
		`\tName:         Is leap year\?`,
		``,
		`Should be accepted by the given predicate:`,
		`\tActual Type:  int`,
		`\tActual Value: 2025`,
		`\tName:         Is leap year\?`)
}

func Test_Check_IsNot(t *testing.T) {
	pt := newTester(t)
	IsNot(pt, leapYear).Name(`Is not leap year?`).Assert(1900)
	IsNot(pt, leapYear).Name(`Is not leap year?`).Assert(1996)
	IsNot(pt, leapYear).Name(`Is not leap year?`).Assert(2000)
	IsNot(pt, leapYear).Name(`Is not leap year?`).Assert(2024)
	IsNot(pt, leapYear).Name(`Is not leap year?`).Assert(2025)
	pt.Check(`Should not be accepted by the given predicate:`,
		`\tActual Type:  int`,
		`\tActual Value: 1996`,
		`\tName:         Is not leap year\?`,
		``,
		`Should not be accepted by the given predicate:`,
		`\tActual Type:  int`,
		`\tActual Value: 2000`,
		`\tName:         Is not leap year\?`,
		``,
		`Should not be accepted by the given predicate:`,
		`\tActual Type:  int`,
		`\tActual Value: 2024`,
		`\tName:         Is not leap year\?`)
}

func Test_Check_StartWith(t *testing.T) {
	pt := newTester(t)
	StartsWith(pt, 12).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should start with the given prefix:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Prefix: 12`,
		`\tExpected Type:   int`)

	StartsWith(pt, []int{1, 2, 3}).Assert(12)
	pt.Check(`Should start with the given prefix:`,
		`\tActual Type:     int`,
		`\tActual Value:    12`,
		`\tExpected Prefix: \[1 2 3\]`,
		`\tExpected Type:   \[\]int`)

	StartsWith(pt, []int{1, 4, 3}).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should start with the given prefix:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Prefix: \[1 4 3\]`,
		`\tExpected Type:   \[\]int`)

	StartsWith(pt, 1).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	StartsWith(pt, []int{1, 2, 3}).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	StartsWith(pt, `World`).Assert(`Hello World`)
	pt.Check(`Should start with the given prefix:`,
		`\tActual Type:     string`,
		`\tActual Value:    "Hello World"`,
		`\tExpected Prefix: "World"`,
		`\tExpected Type:   string`)

	StartsWith(pt, `He`).Assert(`Hello World`)
	pt.Check()

	StartsWith(pt, []int{}).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Must have at least one expected prefix value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)
}

func Test_Check_EndsWith(t *testing.T) {
	pt := newTester(t)
	EndsWith(pt, 12).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should end with the given suffix:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Suffix: 12`,
		`\tExpected Type:   int`)

	EndsWith(pt, []int{1, 2, 3}).Assert(12)
	pt.Check(`Should end with the given suffix:`,
		`\tActual Type:     int`,
		`\tActual Value:    12`,
		`\tExpected Suffix: \[1 2 3\]`,
		`\tExpected Type:   \[\]int`)

	EndsWith(pt, []int{3, 2, 5}).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should end with the given suffix:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Suffix: \[3 2 5\]`,
		`\tExpected Type:   \[\]int`)

	EndsWith(pt, []int{3, 4, 5}).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	EndsWith(pt, ``).Assert(`Hello World`)
	pt.Check(`Must have at least one expected suffix value:`,
		`\tExpected Type: string`,
		`FAIL NOW`)
}

func Test_Check_Empty(t *testing.T) {
	pt := newTester(t)
	Empty(pt).Assert(12)
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:  int`,
		`\tActual Value: 12`)

	s := []int{1, 2, 3, 4, 5}
	Empty(pt).Assert(s)
	pt.Check(`Should be empty:`,
		`\tActual Length: 5`,
		`\tActual Type:   \[\]int`,
		`\tActual Value:  \[1 2 3 4 5\]`)

	Empty(pt).Assert(s[:0])
	pt.Check()
}

func Test_Check_NotEmpty(t *testing.T) {
	pt := newTester(t)
	NotEmpty(pt).Assert(12)
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:  int`,
		`\tActual Value: 12`)

	s := []int{1, 2, 3, 4, 5}
	NotEmpty(pt).Assert(s)
	pt.Check()

	NotEmpty(pt).Assert(s[:0])
	pt.Check(`Should not be empty:`,
		`\tActual Length: 0`,
		`\tActual Type:   \[\]int`,
		`\tActual Value:  \[\]`)
}

func Test_Check_Single(t *testing.T) {
	pt := newTester(t)
	Single(pt).Assert(12)
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:  int`,
		`\tActual Value: 12`)

	s := []int{1, 2, 3, 4, 5}
	Single(pt).Assert(s)
	pt.Check(`Should have one and only one value:`,
		`\tActual Length: 5`,
		`\tActual Type:   \[\]int`,
		`\tActual Value:  \[1 2 3 4 5\]`)

	Single(pt).Assert(s[:0])
	pt.Check(`Should have one and only one value:`,
		`\tActual Length: 0`,
		`\tActual Type:   \[\]int`,
		`\tActual Value:  \[\]`)

	Single(pt).Assert([]int{42})
	pt.Check()

	Single(pt).Assert(`Hello`)
	pt.Check(`Should have one and only one value:`,
		`\tActual Length: 5`,
		`\tActual Type:   string`,
		`\tActual Value:  "Hello"`)

	Single(pt).Assert(`H`)
	pt.Check()
}

func Test_Check_Length(t *testing.T) {
	pt := newTester(t)
	Length(pt, 5).Assert(errors.New(`Oops`))
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:     \*errors\..+`,
		`\tActual Value:    Oops`,
		`\tExpected Length: 5`)

	m1 := map[string]int{`cat`: 5, `dog`: 9}
	Length(pt, 2).Assert(m1)
	pt.Check()

	m2 := map[string]int{`cat`: 5}
	Length(pt, 2).Assert(m2)
	pt.Check(`Should be the expected length:`,
		`\tActual Length:   1`,
		`\tActual Type:     map\[string\]int`,
		`\tActual Value:    map\[cat:5\]`,
		`\tExpected Length: 2`)

	obj1 := lenObj{len: 14}
	Length(pt, 14).Assert(obj1)
	pt.Check()

	obj2 := lengthObj{length: 27}
	Length(pt, 27).Assert(obj2)
	pt.Check()

	obj3 := countObj{count: 336}
	Length(pt, 336).Assert(obj3)
	pt.Check()

	obj4 := badCountObj{}
	Length(pt, 336).Assert(obj4)
	pt.Check(`Error: Vanjie`,
		`FAIL NOW`)
}

func Test_Check_ShorterThan(t *testing.T) {
	pt := newTester(t)
	ShorterThan(pt, 15).Assert(12)
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tMaximum Length: 15`)

	s := `Hello World`
	ShorterThan(pt, 15).Assert(s)
	pt.Check()

	ShorterThan(pt, 8).Assert(s)
	pt.Check(`Should be shorter than the expected length:`,
		`\tActual Length:  11`,
		`\tActual Type:    string`,
		`\tActual Value:   "Hello World"`,
		`\tMaximum Length: 8`)
}

func Test_Check_LongerThan(t *testing.T) {
	pt := newTester(t)
	LongerThan(pt, 15).Assert(12)
	pt.Check(`Should be a type that has length:`,
		`\tActual Type:    int`,
		`\tActual Value:   12`,
		`\tMinimum Length: 15`)

	s := `Hello World`
	LongerThan(pt, 15).Assert(s)
	pt.Check(`Should be longer than the expected length:`,
		`\tActual Length:  11`,
		`\tActual Type:    string`,
		`\tActual Value:   "Hello World"`,
		`\tMinimum Length: 15`)

	LongerThan(pt, 8).Assert(s)
	pt.Check()
}

func Test_Check_NoError(t *testing.T) {
	pt := newTester(t)
	NoError(pt).Assert(errors.New(`Science!`))
	pt.Check(`Should be no error:`,
		`\tActual Type:  \*errors.+`,
		`\tActual Value: Science!`)

	NoError(pt).Assert(nil)
	pt.Check()
}

func Test_Check_MatchError(t *testing.T) {
	pt := newTester(t)

	// Check that failure returns a Check that doesn't cause a nil dereference.
	MatchError(pt, ``).With(`A`, `B`).Required().Assert(nil)
	pt.Check(`Must provide a non-empty regular expression pattern`,
		`FAIL NOW`)

	MatchError(pt, `[[))`).With(`A`, `B`).Required().Assert(nil)
	pt.Check(`Must provide a valid regular expression pattern:`,
		`\tPattern: \[\[\)\)`,
		`FAIL NOW`)

	MatchError(pt, `[[))`).Withf(`A`, `B`).Required().Assert(nil)
	pt.Check(`Must provide a valid regular expression pattern:`,
		`\tPattern: \[\[\)\)`,
		`FAIL NOW`)

	MatchError(pt, `.*`).Assert(nil)
	pt.Check(`Should not be a nil error:`,
		`\tActual Type:  \<nil\>`,
		`\tActual Value: \<nil\>`)

	MatchError(pt, `^[a-z]+$`).Assert(errors.New(`Maths!`))
	pt.Check(`Should have error sting match the given regular expression pattern:`,
		`\tActual Type:  \*errors.+`,
		`\tActual Value: Maths!`,
		`\tPattern:      \^\[a-z\]\+\$`)

	MatchError(pt, `^M[a-z]+s!$`).Assert(errors.New(`Maths!`))
	pt.Check()
}

func Test_Check_ErrorHas(t *testing.T) {
	pt := newTester(t)
	_, err := strconv.Atoi(`apple`)
	ErrorHas[*json.UnsupportedTypeError](pt).Assert(err)
	pt.Check(`Should have an error of the target type in the error tree:`,
		`\tActual Type:  \*strconv\.NumError`,
		`\tActual Value: strconv\.Atoi: parsing "apple": invalid syntax`,
		`\tTarget Type:  \*json\.UnsupportedTypeError`)

	ErrorHas[*strconv.NumError](pt).Assert(err)
	ErrorHas[*strconv.NumError](pt).Assert(fmt.Errorf(`=>%w`, err))
	pt.Check()
}

func Test_Check_Implements(t *testing.T) {
	pt := newTester(t)
	v := 3.14
	Implements[float64](pt).Assert(v)
	pt.Check(`Must provide an interface type:`,
		`\tType: float64`,
		`FAIL NOW`)

	Implements[error](pt).Assert(v)
	pt.Check(`Should implement the target type:`,
		`\tActual Type:  float64`,
		`\tActual Value: 3.14`,
		`\tTarget Type:  error`)

	_, err := strconv.Atoi(`apple`)
	Implements[error](pt).Assert(err)
	pt.Check()

	Implements[utils.Stringer](pt).Assert(err)
	pt.Check(`Should implement the target type:`,
		`\tActual Type:  \*strconv\.NumError`,
		`\tActual Value: strconv\.Atoi: parsing "apple": invalid syntax`,
		`\tTarget Type:  utils\.Stringer`)
}

func Test_Check_ConvertibleTo(t *testing.T) {
	pt := newTester(t)
	var v1 float64 = 3.14
	ConvertibleTo[float64](pt).Assert(v1)
	pt.Check()

	var v2 int = 13
	ConvertibleTo[float64](pt).Assert(v2)
	pt.Check()

	var v3 string = `Hello`
	ConvertibleTo[float64](pt).Assert(v3)
	pt.Check(`Should be convertible to the target type:`,
		`\tActual Type:  string`,
		`\tActual Value: "Hello"`,
		`\tTarget Type:  float64`)

	type tt string
	var v4 tt = `Hello`
	ConvertibleTo[string](pt).Assert(v4)
	pt.Check()
}

func Test_Check_SameType(t *testing.T) {
	pt := newTester(t)
	var v1 float64 = 3.14
	var v2 float64 = 4.5
	SameType(pt, v1).Assert(v2)
	pt.Check()

	SameType(pt, 42).Assert(v1)
	pt.Check(`Should be the same expected type:`,
		`\tActual Type:    float64`,
		`\tActual Value:   3\.14`,
		`\tExpected Type:  int`,
		`\tExpected Value: 42`)

	SameType(pt, `Hello`).Assert(v1)
	pt.Check(`Should be the same expected type:`,
		`\tActual Type:    float64`,
		`\tActual Value:   3\.14`,
		`\tExpected Type:  string`,
		`\tExpected Value: "Hello"`)

	SameType(pt, &v1).Assert(&v2)
	pt.Check()

	SameType(pt, 'a').Assert(int32(24601))
	pt.Check()

	SameType(pt, []int{1, 3, 5}).Assert([]int(nil))
	pt.Check()
}

func Test_Check_NotSameType(t *testing.T) {
	pt := newTester(t)
	var v1 float64 = 3.14
	var v2 float64 = 4.5
	NotSameType(pt, v1).Assert(v2)
	pt.Check(`Should not be the unexpected type:`,
		`\tActual Type:      float64`,
		`\tActual Value:     4\.5`,
		`\tUnexpected Type:  float64`,
		`\tUnexpected Value: 3\.14`)

	NotSameType(pt, 42).Assert(v1)
	pt.Check()

	NotSameType(pt, `Hello`).Assert(v1)
	pt.Check()

	NotSameType(pt, 'a').Assert(int32(24601))
	pt.Check(`Should not be the unexpected type:`,
		`\tActual Type:      int32`,
		`\tActual Value:     24601`,
		`\tUnexpected Type:  int32`,
		`\tUnexpected Value: 97`)
}

func Test_Check_Same(t *testing.T) {
	pt := newTester(t)
	v1 := 3.14
	p1 := &v1
	Same(pt, p1).Assert(&v1)
	pt.Check()

	v2 := 3.14
	p2 := &v2
	Same(pt, p1).Assert(p2)
	pt.Check(`Should be the same:`,
		`\tActual Type:    \*float64`,
		`\tActual Value:   0x[0-9a-f]+`,
		`\tExpected Type:  \*float64`,
		`\tExpected Value: 0x[0-9a-f]+`)

	Same(pt, v2).Assert(v1)
	pt.Check()

	Same(pt, any(3.0)).Assert(3)
	pt.Check(`Should be the same:`,
		`\tActual Type:    int`,
		`\tActual Value:   3`,
		`\tExpected Type:  float64`,
		`\tExpected Value: 3`)
}

func Test_Check_NotSame(t *testing.T) {
	pt := newTester(t)
	NotSame(pt, 3).Assert(3)
	pt.Check(`Should not be the same:`,
		`\tActual Type:    int`,
		`\tActual Value:   3`,
		`\tExpected Type:  int`,
		`\tExpected Value: 3`)

	NotSame(pt, any(4)).Assert(4.0)
	NotSame(pt, 4).Assert(5)
}

func Test_Check_Includes(t *testing.T) {
	pt := newTester(t)
	s := []int{1, 2, 3, 4, 5}
	Includes(pt, []int{}).Assert(s)
	pt.Check(`Must provide at least one expected value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)

	Includes(pt, 4).Assert(s)
	pt.Check()

	Includes(pt, 7).Assert(s)
	pt.Check(`Should have the expected values:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Type:   int`,
		`\tExpected Values: 7`,
		`\tMissing Values:  \[7\]`)

	Includes(pt, []int{3, 4, 2}).Assert(s)
	pt.Check()

	Includes(pt, []int{3, 7, 1, 9}).Assert(s)
	pt.Check(`Should have the expected values:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Type:   \[\]int`,
		`\tExpected Values: \[3 7 1 9\]`,
		`\tMissing Values:  \[7 9\]`)
}

func Test_Check_OneOf(t *testing.T) {
	pt := newTester(t)
	OneOf(pt, []int{}).Assert(4)
	pt.Check(`Must provide at least one expected value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)

	OneOf(pt, []int{3, 4, 2}).Assert(4)
	pt.Check()

	OneOf(pt, []int{3, 4, 2}).Assert(5)
	pt.Check(`Should be one of the expected values:`,
		`\tActual Type:     int`,
		`\tActual Value:    5`,
		`\tExpected Type:   \[\]int`,
		`\tExpected Values: \[3 4 2\]`)

	OneOf(pt, []string{`Cat`, `Dog`, `Pickle`}).Assert(3.14)
	pt.Check(`Should be one of the expected values:`,
		`\tActual Type:     float64`,
		`\tActual Value:    3\.14`,
		`\tExpected Type:   \[\]string`,
		`\tExpected Values: \[Cat Dog Pickle\]`)

	m := map[string]int{`One`: 1, `Two`: 2, `Three`: 3}
	OneOf(pt, m).Assert(`Two`)
	pt.Check(`Should be one of the expected values:`,
		`\tActual Type:     string`,
		`\tActual Value:    "Two"`,
		`\tExpected Type:   map\[string\]int`,
		`\tExpected Values: map\[(?:One:1 ?|Two:2 ?|Three:3 ?){3}\]`)

	OneOf(pt, m).Assert(tuple2.New(`Two`, 2))
	pt.Check()
}

func Test_Check_Excludes(t *testing.T) {
	pt := newTester(t)
	s := []int{1, 2, 3, 4, 5}
	Excludes(pt, []int{}).Assert(s)
	pt.Check(`Must provide at least one unexpected value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)

	Excludes(pt, 4).Assert(s)
	pt.Check(`Should not have the any of the unexpected values:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 3 4 5\]`,
		`\tContained Values:  \[4\]`,
		`\tExpected Type:     int`,
		`\tUnexpected Values: 4`)

	Excludes(pt, 7).Assert(s)
	pt.Check()

	Excludes(pt, []int{3, 7, 1, 9}).Assert(s)
	pt.Check(`Should not have the any of the unexpected values:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 3 4 5\]`,
		`\tContained Values:  \[3 1\]`,
		`\tExpected Type:     \[\]int`,
		`\tUnexpected Values: \[3 7 1 9\]`)

	Excludes(pt, []int{7, 8, 9}).Assert(s)
	pt.Check()
}

func Test_Check_Intersects(t *testing.T) {
	pt := newTester(t)
	s := []int{1, 2, 3, 4, 5}
	Intersects(pt, []int{}).Assert(s)
	pt.Check(`Must provide at least one expected value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)

	Intersects(pt, []int{5, 6, 7, 8, 9}).Assert(s)
	pt.Check()

	Intersects(pt, []int{6, 7, 8, 9}).Assert(s)
	pt.Check(`Should have at least one of the expected values:`,
		`\tActual Type:     \[\]int`,
		`\tActual Value:    \[1 2 3 4 5\]`,
		`\tExpected Type:   \[\]int`,
		`\tExpected Values: \[6 7 8 9\]`)

	Intersects(pt, []int{5, 2, 1, 9}).Assert(s)
	pt.Check()
}

func Test_Check_Sorted(t *testing.T) {
	pt := newTester(t)
	Sorted[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	Sorted[int](pt).Assert([]int{5, 2, 1, 4, 3})
	pt.Check(`Should be in sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[5 2 1 4 3\]`)

	Sorted[int](pt).Assert([]int{5, 4, 3, 2, 1})
	pt.Check(`Should be in sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[5 4 3 2 1\]`)
}

func Test_Check_NotSorted(t *testing.T) {
	pt := newTester(t)
	NotSorted[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should not be in sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[1 2 3 4 5\]`)

	NotSorted[int](pt).Assert([]int{5, 2, 1, 4, 3})
	pt.Check()

	NotSorted[int](pt).Assert([]int{5, 4, 3, 2, 1})
	pt.Check()
}

func Test_Check_DescendingSorted(t *testing.T) {
	pt := newTester(t)
	DescendingSorted[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should be in descending sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[1 2 3 4 5\]`)

	DescendingSorted[int](pt).Assert([]int{5, 2, 1, 4, 3})
	pt.Check(`Should be in descending sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[5 2 1 4 3\]`)

	DescendingSorted[int](pt).Assert([]int{5, 4, 3, 2, 1})
	pt.Check()
}

func Test_Check_NotDescendingSorted(t *testing.T) {
	pt := newTester(t)
	NotDescendingSorted[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	NotDescendingSorted[int](pt).Assert([]int{5, 2, 1, 4, 3})
	pt.Check()

	NotDescendingSorted[int](pt).Assert([]int{5, 4, 3, 2, 1})
	pt.Check(`Should not be in descending sorted order:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[5 4 3 2 1\]`)
}

func Test_Check_Unique(t *testing.T) {
	pt := newTester(t)
	Unique[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check()

	Unique[int](pt).Assert([]int{1, 2, 3, 4, 2})
	pt.Check(`Should have unique values:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[1 2 3 4 2\]`)

	Unique[byte](pt).Assert(`World`)
	pt.Check()

	Unique[byte](pt).Assert(`Hello`)
	pt.Check(`Should have unique values:`,
		`\tActual Type:  string`,
		`\tActual Value: "Hello"`)
}

func Test_Check_NotUnique(t *testing.T) {
	pt := newTester(t)
	NotUnique[int](pt).Assert([]int{1, 2, 3, 4, 5})
	pt.Check(`Should not have all unique values:`,
		`\tActual Type:  \[\]int`,
		`\tActual Value: \[1 2 3 4 5\]`)

	NotUnique[int](pt).Assert([]int{1, 2, 3, 4, 2})
	pt.Check()

	NotUnique[byte](pt).Assert(`World`)
	pt.Check(`Should not have all unique values:`,
		`\tActual Type:  string`,
		`\tActual Value: "World"`)

	NotUnique[byte](pt).Assert(`Hello`)
	pt.Check()
}

func Test_Check_HasKeys(t *testing.T) {
	pt := newTester(t)
	m := map[string]int{`One`: 1, `Two`: 2, `Three`: 3, `Four`: 4, `Five`: 5}
	HasKeys[map[string]int](pt).Assert(m)
	pt.Check(`Must provide at least one expected key:`,
		`\tExpected Type: \[\]string`,
		`FAIL NOW`)

	HasKeys[map[string]int](pt, `Four`).Assert(m)
	pt.Check()

	HasKeys[map[string]int](pt, `Cat`).Assert(m)
	pt.Check(`Should have the expected keys:`,
		`\tActual Type:   map\[string\]int`,
		`\tActual Value:  map\[Five:5 Four:4 One:1 Three:3 Two:2\]`,
		`\tExpected Keys: \[Cat\]`,
		`\tExpected Type: \[\]string`,
		`\tMissing Keys:  \[Cat\]`)

	HasKeys[map[string]int](pt, `Three`, `Four`, `Two`).Assert(m)
	pt.Check()

	HasKeys[map[string]int](pt, `Three`, `Cat`, `One`, `Apple`).Assert(m)
	pt.Check(`Should have the expected keys:`,
		`\tActual Type:   map\[string\]int`,
		`\tActual Value:  map\[Five:5 Four:4 One:1 Three:3 Two:2\]`,
		`\tExpected Keys: \[Three Cat One Apple\]`,
		`\tExpected Type: \[\]string`,
		`\tMissing Keys:  \[Cat Apple\]`)
}

func Test_Check_HasValues(t *testing.T) {
	pt := newTester(t)
	m := map[string]int{`One`: 1, `Two`: 2, `Three`: 3, `Four`: 4, `Five`: 5}
	HasValues[map[string]int](pt).Assert(m)
	pt.Check(`Must provide at least one expected value:`,
		`\tExpected Type: \[\]int`,
		`FAIL NOW`)

	HasValues[map[string]int](pt, 4).Assert(m)
	pt.Check()

	HasValues[map[string]int](pt, 8).Assert(m)
	pt.Check(`Should have the expected values:`,
		`\tActual Type:     map\[string\]int`,
		`\tActual Value:    map\[Five:5 Four:4 One:1 Three:3 Two:2\]`,
		`\tExpected Type:   \[\]int`,
		`\tExpected Values: \[8\]`,
		`\tMissing Values:  \[8\]`)

	HasValues[map[string]int](pt, 3, 4, 2).Assert(m)
	pt.Check()

	HasValues[map[string]int](pt, 3, 8, 1, 12).Assert(m)
	pt.Check(`Should have the expected values:`,
		`\tActual Type:     map\[string\]int`,
		`\tActual Value:    map\[Five:5 Four:4 One:1 Three:3 Two:2\]`,
		`\tExpected Type:   \[\]int`,
		`\tExpected Values: \[3 8 1 12\]`,
		`\tMissing Values:  \[8 12\]`)
}

func Test_Check_EqualElems(t *testing.T) {
	pt := newTester(t)
	s := []int{1, 2, 3, 4, 5}
	EqualElems(pt, []int{}).Assert(s)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 3 4 5\]`,
		`\tExpected Elements: \[\]`,
		`\tExpected Type:     \[\]int`,
		`\tExtra Elements:    \[1 2 3 4 5\]`)

	EqualElems(pt, 4).Assert(s)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 3 4 5\]`,
		`\tExpected Elements: 4`,
		`\tExpected Type:     int`,
		`\tExtra Elements:    \[1 2 3 5\]`)

	EqualElems(pt, []int{2, 3, 4, 5, 6}).Assert(s)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 3 4 5\]`,
		`\tExpected Elements: \[2 3 4 5 6\]`,
		`\tExpected Type:     \[\]int`,
		`\tExtra Elements:    \[1\]`,
		`\tMissing Elements:  \[6\]`)

	EqualElems(pt, []int{5, 3, 1, 4, 2}).Assert(s)
	pt.Check()

	EqualElems(pt, `abcdefghijklmnopqrstuvwxyz`).
		Assert(`the quick brown fox jumps over the lazy dog`)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       string`,
		`\tActual Value:      "the quick brown fox jumps over the lazy dog"`,
		`\tExpected Elements: "abcdefghijklmnopqrstuvwxyz"`,
		`\tExpected Type:     string`,
		`\tExtra Elements:    \[' '\]`)
}

func Test_Check_SameElems(t *testing.T) {
	pt := newTester(t)
	s := []int{1, 2, 2, 3, 4, 5, 5, 5}
	SameElems(pt, []int{}).Assert(s)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 2 3 4 5 5 5\]`,
		`\tExpected Elements: \[\]`,
		`\tExpected Type:     \[\]int`,
		`\tExtra Elements:    \[1 2\(x2\) 3 4 5\(x3\)\]`)

	SameElems(pt, []int{1, 2, 2, 2, 3, 3, 4, 5}).Assert(s)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       \[\]int`,
		`\tActual Value:      \[1 2 2 3 4 5 5 5\]`,
		`\tExpected Elements: \[1 2 2 2 3 3 4 5\]`,
		`\tExpected Type:     \[\]int`,
		`\tExtra Elements:    \[5\(x2\)\]`,
		`\tMissing Elements:  \[2 3\]`)

	SameElems(pt, []int{1, 2, 2, 3, 4, 5, 5, 5}).Assert(s)
	pt.Check()

	SameElems(pt, `abcdefghijklmnopqrstuvwxyz`).
		Assert(`the quick brown fox jumps over the lazy dog`)
	pt.Check(`Should have the expected elements:`,
		`\tActual Type:       string`,
		`\tActual Value:      "the quick brown fox jumps over the lazy dog"`,
		`\tExpected Elements: "abcdefghijklmnopqrstuvwxyz"`,
		`\tExpected Type:     string`,
		`\tExtra Elements:    \[' '\(x8\) 'e'\(x2\) 'h' 'o'\(x3\) 'r' 't' 'u'\]`)
}

type pseudoTester struct {
	t   *testing.T
	buf *bytes.Buffer
}

func newTester(t *testing.T) *pseudoTester {
	return &pseudoTester{
		t:   t,
		buf: &bytes.Buffer{},
	}
}

func (p pseudoTester) Error(msg ...any) {
	if _, err := p.buf.WriteString(fmt.Sprint(msg...)); err != nil {
		panic(err)
	}
}

func (p pseudoTester) FailNow() {
	if _, err := p.buf.WriteString("FAIL NOW\n"); err != nil {
		panic(err)
	}
}

func (p pseudoTester) Check(patterns ...string) {
	actual := p.buf.String()
	if len(actual) <= 0 && len(patterns) <= 0 {
		return
	}

	p.t.Helper()

	if len(actual) > 0 {
		var ok bool
		actual, ok = strings.CutPrefix(actual, "\n")
		if !ok {
			p.t.Errorf("\nExpected a new line prefix but didn't have one.\n\t%q", actual)
		}
		actual, ok = strings.CutSuffix(actual, "\n")
		if !ok {
			p.t.Errorf("\nExpected a new line suffix but didn't have one.\n\t%q", actual)
		}
	}
	lines := strings.Split(actual, "\n")

	for i, pattern := range patterns {
		if !strings.HasPrefix(pattern, `^`) {
			pattern = `^` + pattern
		}
		if !strings.HasSuffix(pattern, `$`) {
			pattern += `$`
		}
		patterns[i] = pattern
	}

	results := diff.Default().Diff(data.Regex(patterns, lines))
	if results.HasDiff() {
		delta := strings.Join(diff.PlusMinus(results, patterns, lines), "\n\t")
		p.t.Helper()
		p.t.Errorf("\nUnexpected Results:\n\t%s\nGotten:\n\t%s\n",
			delta, strings.Join(lines, "\n\t"))
	}
	p.buf.Reset()
}

type pseudoTesterWithHelper struct {
	pseudoTester
}

func (p pseudoTesterWithHelper) Helper() {
	if _, err := p.buf.WriteString("\nHelper\n"); err != nil {
		panic(err)
	}
}

func leapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

type (
	pseudoStringer struct{ text string }
	lenObj         struct{ len int }
	lengthObj      struct{ length int }
	countObj       struct{ count int }
	badCountObj    struct{}
)

func (ps pseudoStringer) String() string { return ps.text }
func (n lenObj) Len() int                { return n.len }
func (n lengthObj) Length() int          { return n.length }
func (n countObj) Count() int            { return n.count }
func (n badCountObj) Count() int         { panic(terror.New(`Vanjie`)) }
