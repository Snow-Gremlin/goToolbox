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
	check.MatchError(t, `^too many arguments \{arguments: cat, gotten: 1, maximum: 0\}$`).
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
	check.MatchError(t, `^too many arguments \{arguments: mouse, gotten: 1, maximum: 0\}$`).Assert(err)

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
	check.MatchError(t, `^too many arguments \{arguments: dog, -24, 3\.125, `+
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
	check.MatchError(t, `error setting variant argument \{arguments: 1, cat\}: `+
		`unable to parse value \{input: cat, type: int\}: strconv\.ParseInt: parsing "cat": `+
		`invalid syntax: invalid syntax`).Assert(err)
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
