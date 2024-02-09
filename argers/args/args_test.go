package args

import (
	"errors"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Args_Empty(t *testing.T) {
	err := New().Process([]string{})
	check.NoError(t).Assert(err)

	err = New().Process([]string{`cat`})
	check.MatchError(t, `^too many optional arguments \{arguments: cat, gotten: 1, maximum: 0\}$`).
		Assert(err)

	err = New().Process([]string{`-c`})
	check.MatchError(t, `^unknown short name found in arguments `+
		`\{argument: -c, argument index: 0, name: c, name index: 0\}$`).
		Assert(err)

	err = New().Process([]string{`--cat`})
	check.MatchError(t, `^unknown long name found in arguments \{argument: --cat, argument index: 0\}$`).
		Assert(err)
}

func Test_Args_Positional(t *testing.T) {
	var s1, s2 string
	r := New().PosStr(&s1).PosStr(&s2)
	err := r.Process([]string{`cat`, `dog`})
	check.NoError(t).Assert(err)
	check.Equal(t, `cat`).Assert(s1)
	check.Equal(t, `dog`).Assert(s2)

	err = r.Process([]string{`cat`})
	check.MatchError(t, `^not enough positional arguments \{arguments: cat, gotten: 1, needed: 2\}$`).Assert(err)

	err = r.Process([]string{`cat`, `dog`, `mouse`})
	check.MatchError(t, `^too many optional arguments \{arguments: mouse, gotten: 1, maximum: 0\}$`).Assert(err)

	check.MatchError(t, `^must provide a non-nil target pointer for a positional argument$`).
		Panic(func() { r.PosInt(nil) })

	var i3 int
	r.PosInt(&i3)
	err = r.Process([]string{`pig`, `cow`, `42`})
	check.NoError(t).Assert(err)
	check.Equal(t, `pig`).Assert(s1)
	check.Equal(t, `cow`).Assert(s2)
	check.Equal(t, 42).Assert(i3)

	err = r.Process([]string{`cat`, `dog`, `mouse`})
	check.MatchError(t, `^error setting positional argument \{argument: mouse, argument index: 3\}: `+
		`unable to parse value \{input: mouse, type: int\}: `+
		`strconv\.ParseInt: parsing \"mouse\": invalid syntax: invalid syntax$`).Assert(err)
}

func Test_Args_Positional_More(t *testing.T) {
	b4, f5 := false, 0.0
	r := New().PosBool(&b4).PosFloat(&f5)
	err := r.Process([]string{`true`, `-6.45`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b4)
	check.Equal(t, -6.45).Assert(f5)
}

func Test_Args_Flags(t *testing.T) {
	var b1, b2 bool
	r := New().Flag(&b1, `a`, `apple`).Flag(&b2, `v`, `verbose`)

	err := r.Process([]string{})
	check.NoError(t).Assert(err)
	check.False(t).Assert(b1)
	check.False(t).Assert(b2)

	b1, b2 = false, false
	err = r.Process([]string{`-a`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b1)
	check.False(t).Assert(b2)

	b1, b2 = false, false
	err = r.Process([]string{`--verbose`})
	check.NoError(t).Assert(err)
	check.False(t).Assert(b1)
	check.True(t).Assert(b2)

	b1, b2 = false, false
	err = r.Process([]string{`--apple`, `--verbose`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b1)
	check.True(t).Assert(b2)

	b1, b2 = false, false
	err = r.Process([]string{`-av`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b1)
	check.True(t).Assert(b2)

	b1, b2 = false, false
	err = r.Process([]string{`-b`})
	check.MatchError(t, `^unknown short name found in arguments \{argument: -b, `+
		`argument index: 0, name: b, name index: 0\}$`).Assert(err)

	err = r.Process([]string{`--banana`})
	check.MatchError(t, `^unknown long name found in arguments `+
		`\{argument: --banana, argument index: 0\}$`).Assert(err)

	check.MatchError(t, `^must provide a non-nil target pointer for an argument flag$`).
		Panic(func() { r.Flag(nil, `n`, `nope`) })

	var b3 bool
	check.MatchError(t, `^may not create a flag with a short name used by a prior flag \{short name: v\}$`).
		Panic(func() { r.Flag(&b3, `v`, `nope`) })
	check.MatchError(t, `^may not create a flag with a long name used by a prior flag \{long name: verbose\}$`).
		Panic(func() { r.Flag(&b3, `n`, `verbose`) })
	check.MatchError(t, `^may not create a flag with an invalid short name \{short name: "4"\}$`).
		Panic(func() { r.Flag(&b3, `4`, `nope`) })
	check.MatchError(t, `^may not create a flag with an invalid long name \{long name: "42"\}$`).
		Panic(func() { r.Flag(&b3, `n`, `42`) })
	check.MatchError(t, `^may not add a flag without a least one name$`).
		Panic(func() { r.Flag(&b3, ``, ``) })

	r.FlagFunc(func() error { return errors.New(`boom`) }, `b`, `boom`)
	err = r.Process([]string{`-b`})
	check.MatchError(t, `^error setting flag \{argument: -b, argument index: 0, `+
		`name: b, name index: 0\}: boom$`).Assert(err)

	err = r.Process([]string{`--boom`})
	check.MatchError(t, `^error setting flag \{argument: --boom, argument index: 0\}: boom$`).Assert(err)
}

func Test_Args_Flags_Strings(t *testing.T) {
	const (
		lvlInfo = `level_info`
		lvlWarn = `level_warn`
		lvlErr  = `level_err`
	)

	var s1 string
	r := New().FlagStr(&s1, lvlInfo, `i`, `info`).
		FlagStr(&s1, lvlWarn, `w`, `warning`).
		FlagStr(&s1, lvlErr, `e`, `error`)

	err := r.Process([]string{`-i`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlInfo).Assert(s1)

	err = r.Process([]string{`-ie`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlErr).Assert(s1)

	err = r.Process([]string{`-ie`, `--warning`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlWarn).Assert(s1)

	err = r.Process([]string{`-wew`, `--info`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlInfo).Assert(s1)
}

func Test_Args_Flags_Ints(t *testing.T) {
	const (
		lvlInfo = iota
		lvlWarn
		lvlErr
	)

	var i1 int
	r := New().FlagInt(&i1, lvlInfo, `i`, `info`).
		FlagInt(&i1, lvlWarn, `w`, `warning`).
		FlagInt(&i1, lvlErr, `e`, `error`)

	err := r.Process([]string{`-i`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlInfo).Assert(i1)

	err = r.Process([]string{`--warning`})
	check.NoError(t).Assert(err)
	check.Equal(t, lvlWarn).Assert(i1)
}

func Test_Args_Named(t *testing.T) {
	b1, s2, i3, f4 := false, ``, 0, 0.0
	r := New().NamedBool(&b1, `a`, `apple`).
		NamedStr(&s2, `b`, `banana`).
		NamedInt(&i3, `c`, `cat`).
		NamedFloat(&f4, `d`, `dog`)

	err := r.Process([]string{})
	check.NoError(t).Assert(err)
	check.False(t).Assert(b1)
	check.Empty(t).Assert(s2)
	check.Zero(t).Assert(i3)
	check.Zero(t).Assert(f4)

	err = r.Process([]string{`-b`, `fruit`, `-c`, `-42`})
	check.NoError(t).Assert(err)
	check.False(t).Assert(b1)
	check.Equal(t, `fruit`).Assert(s2)
	check.Equal(t, -42).Assert(i3)
	check.Zero(t).Assert(f4)

	b1, s2, i3, f4 = false, ``, 0, 0.0
	err = r.Process([]string{`--apple`, `true`, `--dog`, `3.14`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b1)
	check.Empty(t).Assert(s2)
	check.Zero(t).Assert(i3)
	check.Equal(t, 3.14).Assert(f4)

	b1, s2, i3, f4 = false, ``, 0, 0.0
	err = r.Process([]string{`-b`, `--dog`})
	check.NoError(t).Assert(err)
	check.False(t).Assert(b1)
	check.Equal(t, `--dog`).Assert(s2)
	check.Zero(t).Assert(i3)
	check.Zero(t).Assert(f4)

	b1, s2, i3, f4 = false, ``, 0, 0.0
	err = r.Process([]string{`-ab`, `true`, `dog`})
	check.MatchError(t, `^may not have a short named value anywhere but the `+
		`end of a flag group \{argument: -ab, argument index: 0, name: a, name index: 0\}$`).
		Assert(err)

	err = r.Process([]string{`-c`, `pickle`})
	check.MatchError(t, `^error setting named argument \{argument: -c, `+
		`argument index: 0, name: c, name index: 0\}: unable to parse value `+
		`\{input: pickle, type: int\}: strconv\.ParseInt: parsing "pickle": `+
		`invalid syntax: invalid syntax$`).
		Assert(err)

	err = r.Process([]string{`--cat`, `pickle`})
	check.MatchError(t, `^error setting named argument \{argument: --cat, `+
		`argument index: 0\}: unable to parse value `+
		`\{input: pickle, type: int\}: strconv\.ParseInt: parsing "pickle": `+
		`invalid syntax: invalid syntax$`).
		Assert(err)

	err = r.Process([]string{`-b`})
	check.MatchError(t, `^no value found for short named value at end of `+
		`arguments {argument: -b, argument index: 0, name: b, name index: 0}$`).
		Assert(err)

	err = r.Process([]string{`--banana`})
	check.MatchError(t, `^no value found for long named value at end of `+
		`arguments \{argument: --banana, argument index: 0\}$`).
		Assert(err)

	check.MatchError(t, `^may not create a named input with a short name used `+
		`by a prior named input \{short name: a\}$`).
		Panic(func() { r.NamedBool(&b1, `a`, ``) })

	check.MatchError(t, `^may not create a named input with a long name used `+
		`by a prior named input \{long name: dog\}$`).
		Panic(func() { r.NamedBool(&b1, ``, `dog`) })

	check.MatchError(t, `^may not add a named input without a least one name$`).
		Panic(func() { r.NamedBool(&b1, ``, ``) })

	check.MatchError(t, `^must provide a non-nil target pointer for a named argument$`).
		Panic(func() { r.NamedBool(nil, `e`, ``) })
}

func Test_Args_Optional(t *testing.T) {
	s1, i2, f3, b4 := ``, 0, 0.0, false
	r := New().OptionalStr(&s1).OptionalInt(&i2).OptionalFloat(&f3).OptionalBool(&b4)

	err := r.Process([]string{})
	check.NoError(t).Assert(err)
	check.Empty(t).Assert(s1)
	check.Zero(t).Assert(i2)
	check.Zero(t).Assert(f3)
	check.False(t).Assert(b4)

	err = r.Process([]string{`cat`, `72`, `84.4`})
	check.NoError(t).Assert(err)
	check.Equal(t, `cat`).Assert(s1)
	check.Equal(t, 72).Assert(i2)
	check.Equal(t, 84.4).Assert(f3)
	check.False(t).Assert(b4)

	s1, i2, f3, b4 = ``, 0, 0.0, false
	err = r.Process([]string{`dog`, `-24`, `3.125`, `true`})
	check.NoError(t).Assert(err)
	check.Equal(t, `dog`).Assert(s1)
	check.Equal(t, -24).Assert(i2)
	check.Equal(t, 3.125).Assert(f3)
	check.True(t).Assert(b4)

	err = r.Process([]string{`dog`, `-24`, `3.125`, `true`, `pancake`})
	check.MatchError(t, `^too many optional arguments \{arguments: dog, -24, 3\.125, `+
		`true, pancake, gotten: 5, maximum: 4\}$`).
		Assert(err)

	err = r.Process([]string{`dog`, `pancake`})
	check.MatchError(t, `^error setting optional argument \{argument index: 1, `+
		`arguments: dog, pancake\}: unable to parse value \{input: pancake, type: int\}: `+
		`strconv\.ParseInt: parsing "pancake": invalid syntax: invalid syntax$`).
		Assert(err)

	check.MatchError(t, `^must provide a non-nil target pointer for an optional argument$`).
		Panic(func() { r.OptionalInt(nil) })

	var varStr []string
	check.MatchError(t, `^may not add a variant argument after an optional argument has been added$`).
		Panic(func() { r.VarStr(&varStr) })

	check.MatchError(t, `^may not add a new positional argument after an optional argument has been added$`).
		Panic(func() { r.PosStr(&s1) })
}

func Test_Args_Variant(t *testing.T) {
	var iVar []int
	r := New().VarInt(&iVar)

	err := r.Process([]string{})
	check.NoError(t).Assert(err)
	check.Empty(t).Assert(iVar)

	err = r.Process([]string{`1`, `-5`, `42`})
	check.NoError(t).Assert(err)
	check.Equal(t, []int{1, -5, 42}).Assert(iVar)

	err = r.Process([]string{`1`, `cat`})
	check.MatchError(t, `^error setting variant argument \{arguments: 1, cat\}: `+
		`unable to parse value \{input: cat, type: int\}: strconv\.ParseInt: parsing "cat": `+
		`invalid syntax: invalid syntax$`).Assert(err)
	check.Equal(t, []int{1, -5, 42}).Assert(iVar)

	check.MatchError(t, `^must provide a non-nil target pointer for a variadic argument$`).
		Panic(func() { r.VarStr(nil) })

	var s1 string
	check.MatchError(t, `^may not add a new positional argument after a variant argument has been added$`).
		Panic(func() { r.PosStr(&s1) })

	check.MatchError(t, `^may not add a new variant argument after a variant argument has already been added$`).
		Panic(func() { r.VarInt(&iVar) })

	check.MatchError(t, `^may not add an optional argument after a variant argument has been added$`).
		Panic(func() { r.OptionalStr(&s1) })
}

func Test_Args_Variant_Bool(t *testing.T) {
	var bVar []bool
	r := New().VarBool(&bVar)

	err := r.Process([]string{`1`, `0`, `F`, `T`})
	check.NoError(t).Assert(err)
	check.Equal(t, []bool{true, false, false, true}).Assert(bVar)
}

func Test_Args_Variant_Float(t *testing.T) {
	var fVar []float64
	r := New().VarFloat(&fVar)

	err := r.Process([]string{`2.0`, `2.25`, `3.5`})
	check.NoError(t).Assert(err)
	check.Equal(t, []float64{2.0, 2.25, 3.5}).Assert(fVar)
}

func Test_Args_Combination(t *testing.T) {
	var b1 bool
	var s2, s3, s4, s5, s6 string
	r := New().Flag(&b1, `v`, `verbose`).NamedStr(&s2, `i`, `input`).
		PosStr(&s3).PosStr(&s4).OptionalStr(&s5).OptionalStr(&s6)

	err := r.Process([]string{`-v`, `cat`, `--input`, `keyboard`, `dog`, `pig`})
	check.NoError(t).Assert(err)
	check.True(t).Assert(b1)
	check.Equal(t, `keyboard`).Assert(s2)
	check.Equal(t, `cat`).Assert(s3)
	check.Equal(t, `dog`).Assert(s4)
	check.Equal(t, `pig`).Assert(s5)
	check.Empty(t).Assert(s6)
}

func Test_Args_Struct_Pos(t *testing.T) {
	t0 := struct {
		X     int
		Y     int
		moose int
		Title string
	}{X: 0, Y: 0, moose: 0, Title: ``}
	r := New().Struct(&t0)

	check.MatchError(t, `^not enough positional arguments \{arguments: 12, 34, gotten: 2, needed: 3\}$`).
		Assert(r.Process([]string{`12`, `34`}))

	check.MatchError(t, `^error setting positional argument \{argument: Dog, argument index: 3\}: `+
		`unable to parse value \{input: Dog, type: int\}: strconv.ParseInt: parsing \"Dog\": `+
		`invalid syntax: invalid syntax$`).
		Assert(r.Process([]string{`12`, `Dog`, `Cat`}))

	check.NoError(t).Assert(r.Process([]string{`12`, `34`, `Cat`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Equal(t, 34).Assert(t0.Y)
	check.Equal(t, `Cat`).Assert(t0.Title)

	check.MatchError(t, `^must provide a non-nil pointer to a structure$`).
		Panic(func() { New().Struct(nil) })
	check.MatchError(t, `^must provide a non-nil pointer to a structure$`).
		Panic(func() { New().Struct(t0) })
	check.MatchError(t, `^must provide a non-nil pointer to a structure$`).
		Panic(func() {
			v := 10
			New().Struct(&v)
		})

	t1 := struct {
		X int `args:"  "`
	}{X: 0}
	check.NoError(t).Assert(New().Struct(&t1).Process([]string{`76`}))
	check.Equal(t, 76).Assert(t1.X)
}

func Test_Args_Struct_Types(t *testing.T) {
	t0 := struct {
		B    bool
		Str  string
		I00  int
		I08  int8
		I16  int16
		I32  int32
		I64  int64
		U00  uint
		U08  uint8
		U16  uint16
		U32  uint32
		U64  uint64
		F32  float32
		F64  float64
		C64  complex64
		C128 complex128
	}{B: false, Str: ``, I00: 0, I08: 0, I16: 0, I32: 0, I64: 0, U00: 0, U08: 0, U16: 0, U32: 0, U64: 0, F32: 0, F64: 0, C64: 0, C128: 0}
	check.NoError(t).Assert(New().Struct(&t0).Process([]string{
		`true`, `Hello`,
		`-1`, `-2`, `-3`, `-4`, `-5`,
		`1`, `2`, `3`, `4`, `5`,
		`3.14`, `6.28`, `5+3i`, `6+8i`,
	}))
	check.True(t).Assert(t0.B)
	check.Equal[string](t, `Hello`).Assert(t0.Str)
	check.Equal[int](t, -1).Assert(t0.I00)
	check.Equal[int8](t, -2).Assert(t0.I08)
	check.Equal[int16](t, -3).Assert(t0.I16)
	check.Equal[int32](t, -4).Assert(t0.I32)
	check.Equal[int64](t, -5).Assert(t0.I64)
	check.Equal[uint](t, 1).Assert(t0.U00)
	check.Equal[uint8](t, 2).Assert(t0.U08)
	check.Equal[uint16](t, 3).Assert(t0.U16)
	check.Equal[uint32](t, 4).Assert(t0.U32)
	check.Equal[uint64](t, 5).Assert(t0.U64)
	check.Equal[float32](t, 3.14).Assert(t0.F32)
	check.Equal[float64](t, 6.28).Assert(t0.F64)
	check.Equal[complex64](t, 5+3i).Assert(t0.C64)
	check.Equal[complex128](t, 6+8i).Assert(t0.C128)

	t1 := struct {
		Ptr *int
	}{Ptr: nil}
	check.MatchError(t, `^must provide a non-nil target pointer for a positional argument \{field: Ptr, tag: \}$`).
		Panic(func() { New().Struct(&t1) })

	var i0 int
	t1.Ptr = &i0
	check.NoError(t).Assert(New().Struct(&t1).Process([]string{`1234`}))
	check.Equal(t, 1234).Assert(i0)

	t2 := struct {
		Foo func(x, y int) int
	}{Foo: nil}
	check.MatchError(t, `^unexpected field type for arguments \{field: Foo\}$`).
		Panic(func() { New().Struct(&t2) })
}

func Test_Args_Struct_Skip(t *testing.T) {
	t0 := struct {
		X int
		Y int `args:"skip"`
		Z int
	}{X: 0, Y: 0, Z: 0}
	r := New().Struct(&t0)

	check.NoError(t).Assert(r.Process([]string{`12`, `34`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Zero(t).Assert(t0.Y)
	check.Equal(t, 34).Assert(t0.Z)
}

func Test_Args_Struct_Optional(t *testing.T) {
	t0 := struct {
		X int
		Y int `args:"optional"`
		Z int `args:"optional"`
	}{X: 0, Y: 0, Z: 0}
	r := New().Struct(&t0)

	check.MatchError(t, `^not enough positional arguments \{arguments: , gotten: 0, needed: 1\}$`).
		Assert(r.Process([]string{}))

	t0.X, t0.Y, t0.Z = 0, 0, 0
	check.NoError(t).Assert(r.Process([]string{`42`}))
	check.Equal(t, 42).Assert(t0.X)
	check.Zero(t).Assert(t0.Y)
	check.Zero(t).Assert(t0.Z)

	t0.X, t0.Y, t0.Z = 0, 0, 0
	check.NoError(t).Assert(r.Process([]string{`12`, `34`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Equal(t, 34).Assert(t0.Y)
	check.Zero(t).Assert(t0.Z)

	t0.X, t0.Y, t0.Z = 0, 0, 0
	check.NoError(t).Assert(r.Process([]string{`56`, `78`, `90`}))
	check.Equal(t, 56).Assert(t0.X)
	check.Equal(t, 78).Assert(t0.Y)
	check.Equal(t, 90).Assert(t0.Z)

	check.MatchError(t, `^too many optional arguments \{arguments: 2, 3, 4, gotten: 3, maximum: 2\}$`).
		Assert(r.Process([]string{`1`, `2`, `3`, `4`}))

	check.MatchError(t, `^error setting optional argument \{argument index: 2, arguments: 2, cat\}: `+
		`unable to parse value \{input: cat, type: int\}: strconv.ParseInt: `+
		`parsing "cat": invalid syntax: invalid syntax`).
		Assert(r.Process([]string{`1`, `2`, `cat`}))
}

func Test_Args_Struct_Variadic(t *testing.T) {
	t0 := struct {
		X int
		Y []int
	}{X: 0, Y: nil}
	r := New().Struct(&t0)

	check.NoError(t).Assert(r.Process([]string{`12`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Zero(t).Assert(t0.Y)

	check.NoError(t).Assert(r.Process([]string{`12`, `34`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Equal(t, []int{34}).Assert(t0.Y)

	check.NoError(t).Assert(r.Process([]string{`12`, `34`, `56`, `78`}))
	check.Equal(t, 12).Assert(t0.X)
	check.Equal(t, []int{34, 56, 78}).Assert(t0.Y)

	check.MatchError(t, `error setting variant argument \{arguments: 34, cat, 78\}: `+
		`unable to parse value \{input: cat, type: int\}: strconv.ParseInt: `+
		`parsing "cat": invalid syntax: invalid syntax`).
		Assert(r.Process([]string{`12`, `34`, `cat`, `78`}))

	t1 := struct {
		X int
		Y []int `args:"cat"`
	}{X: 0, Y: nil}
	check.MatchError(t, `^invalid tag on a variadic argument value. May only have the skip tag\. \{field: Y, tag: cat\}$`).
		Panic(func() { New().Struct(&t1) })
}

func Test_Args_Struct_Flags(t *testing.T) {
	t0 := struct {
		X bool `args:"flag,a,apple"`
		Y int  `args:"flag,b,banana,42"`
	}{X: false, Y: 0}
	r := New().Struct(&t0)

	check.NoError(t).Assert(r.Process([]string{}))
	check.False(t).Assert(t0.X)
	check.Zero(t).Assert(t0.Y)

	t0.X, t0.Y = false, 0
	check.NoError(t).Assert(r.Process([]string{`-a`}))
	check.True(t).Assert(t0.X)
	check.Zero(t).Assert(t0.Y)

	t0.X, t0.Y = false, 0
	check.NoError(t).Assert(r.Process([]string{`-ab`}))
	check.True(t).Assert(t0.X)
	check.Equal(t, 42).Assert(t0.Y)

	t0.X, t0.Y = false, 0
	check.NoError(t).Assert(r.Process([]string{`-ba`}))
	check.True(t).Assert(t0.X)
	check.Equal(t, 42).Assert(t0.Y)

	t0.X, t0.Y = false, 0
	check.NoError(t).Assert(r.Process([]string{`-a`, `-b`}))
	check.True(t).Assert(t0.X)
	check.Equal(t, 42).Assert(t0.Y)

	t0.X, t0.Y = false, 0
	check.NoError(t).Assert(r.Process([]string{`--banana`}))
	check.False(t).Assert(t0.X)
	check.Equal(t, 42).Assert(t0.Y)

	t1 := struct {
		X int `args:"flag"`
	}{X: 0}
	check.MatchError(t, `^the tag on a flag must have three or four values, `+
		`i.e. "flag,v,verbose" \{field: X, tag: flag\}$`).
		Panic(func() { New().Struct(&t1) })

	t2 := struct {
		X int `args:"flag,a,apples"`
	}{X: 0}
	check.MatchError(t, `^the tag on a flag must have the fourth value to use as a `+
		`default value when the type is not a bool \{field: X, tag: flag,a,apples\}$`).
		Panic(func() { New().Struct(&t2) })

	t3 := struct {
		X int `args:"flag,a,apples,cat"`
	}{X: 0}
	check.MatchError(t, `^the default value in the tag for a flag could not be parsed `+
		`\{field: X, tag: flag,a,apples,cat\}: unable to parse value \{input: cat, type: int\}: `+
		`strconv.ParseInt: parsing "cat": invalid syntax: invalid syntax$`).
		Panic(func() { New().Struct(&t3) })

	t4 := struct {
		X int `args:"flag,apples,bananas,12"`
	}{X: 0}
	check.MatchError(t, `^may not create a flag with an invalid short name `+
		`\{field: X, short name: "apples", tag: flag,apples,bananas,12\}$`).
		Panic(func() { New().Struct(&t4) })

	t5 := struct {
		X int `args:"flag,,,12"`
	}{X: 0}
	check.MatchError(t, `^may not add a flag without a least one name \{field: X, tag: flag,,,12\}$`).
		Panic(func() { New().Struct(&t5) })
}

func Test_Args_Struct_Named(t *testing.T) {
	t0 := struct {
		X int    `args:"i,input"`
		Y string `args:"o,output"`
	}{X: 0, Y: ``}
	r := New().Struct(&t0)

	check.NoError(t).Assert(r.Process([]string{}))
	check.Zero(t).Assert(t0.X)
	check.Empty(t).Assert(t0.Y)

	t0.X, t0.Y = 0, ``
	check.NoError(t).Assert(r.Process([]string{`-i`, `23`}))
	check.Equal(t, 23).Assert(t0.X)
	check.Empty(t).Assert(t0.Y)

	t0.X, t0.Y = 0, ``
	check.NoError(t).Assert(r.Process([]string{`-o`, `Hello World`}))
	check.Zero(t).Assert(t0.X)
	check.Equal(t, `Hello World`).Assert(t0.Y)

	t0.X, t0.Y = 0, ``
	check.NoError(t).Assert(r.Process([]string{`--input`, `42`}))
	check.Equal(t, 42).Assert(t0.X)
	check.Empty(t).Assert(t0.Y)

	check.MatchError(t, `^error setting named argument \{argument: -i, argument index: 0, `+
		`name: i, name index: 0\}: unable to parse value \{input: cat, type: int\}: `+
		`strconv\.ParseInt: parsing "cat": invalid syntax: invalid syntax$`).
		Assert(r.Process([]string{`-i`, `cat`}))

	t1 := struct {
		X int `args:"i,input,boom"`
	}{X: 0}
	check.MatchError(t, `^the tag on a named input must have two values, `+
		`i\.e\. "i,input" \{field: X, tag: i,input,boom\}$`).
		Panic(func() { New().Struct(&t1) })

	t2 := struct {
		X string `args:"i,"`
		Y string `args:",output"`
	}{X: ``, Y: ``}
	r = New().Struct(&t2)

	check.NoError(t).Assert(r.Process([]string{`-i`, `cat`, `--output`, `dog`}))
	check.Equal(t, `cat`).Assert(t2.X)
	check.Equal(t, `dog`).Assert(t2.Y)

	t3 := struct {
		X string `args:"  i,  input"`
		Y string `args:"  o ,  output"`
	}{X: ``, Y: ``}
	r = New().Struct(&t3)

	check.NoError(t).Assert(r.Process([]string{`-i`, `cat`, `--output`, `dog`}))
	check.Equal(t, `cat`).Assert(t3.X)
	check.Equal(t, `dog`).Assert(t3.Y)
}
