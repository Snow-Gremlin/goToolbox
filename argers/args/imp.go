package args

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/argers"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type readerImp struct {
	shortFlags map[string]argers.FlagHandle
	longFlags  map[string]argers.FlagHandle
	shortNamed map[string]argers.ArgHandle
	longNamed  map[string]argers.ArgHandle
	pos        []argers.ArgHandle
	optionals  []argers.ArgHandle
	variant    argers.VarHandle
}

var (
	shortNameMatch = utils.LazyMatcher(`^[A-Za-z]$`)
	longNameMatch  = utils.LazyMatcher(`^[A-Za-z]\w*$`)
	shortArgMatch  = utils.LazyMatcher(`^-[A-Za-z]+$`)
	longArgMatch   = utils.LazyMatcher(`^--[A-Za-z]\w*$`)
)

func (imp *readerImp) validateNames(usage, short, long string) {
	hasName := false
	if len(short) > 0 {
		if !shortNameMatch(short) {
			panic(terror.New(`may not create a `+usage+` with an invalid short name`).
				With(`short name`, strconv.Quote(short)))
		}
		if _, exists := imp.shortFlags[short]; exists {
			panic(terror.New(`may not create a `+usage+` with a short name used by a prior flag`).
				With(`short name`, short))
		}
		if _, exists := imp.shortNamed[short]; exists {
			panic(terror.New(`may not create a `+usage+` with a short name used by a prior named input`).
				With(`short name`, short))
		}
		hasName = true
	}
	if len(long) > 0 {
		if !longNameMatch(long) {
			panic(terror.New(`may not create a `+usage+` with an invalid long name`).
				With(`long name`, strconv.Quote(long)))
		}
		if _, exists := imp.longFlags[long]; exists {
			panic(terror.New(`may not create a `+usage+` with a long name used by a prior flag`).
				With(`long name`, long))
		}
		if _, exists := imp.longNamed[long]; exists {
			panic(terror.New(`may not create a `+usage+` with a long name used by a prior named input`).
				With(`long name`, long))
		}
		hasName = true
	}
	if !hasName {
		panic(terror.New(`may not add a ` + usage + ` without a least one name`))
	}
}

func (imp *readerImp) FlagFunc(handle argers.FlagHandle, short, long string) argers.Reader {
	imp.validateNames(`flag`, short, long)
	if imp.shortFlags == nil {
		imp.shortFlags = map[string]argers.FlagHandle{}
		imp.longFlags = map[string]argers.FlagHandle{}
	}
	imp.shortFlags[short] = handle
	imp.longFlags[long] = handle
	return imp
}

func (imp *readerImp) Flag(target *bool, short, long string) argers.Reader {
	return Flag(imp, target, true, short, long)
}

func (imp *readerImp) FlagStr(target *string, value, short, long string) argers.Reader {
	return Flag(imp, target, value, short, long)
}

func (imp *readerImp) FlagInt(target *int, value int, short, long string) argers.Reader {
	return Flag(imp, target, value, short, long)
}

func (imp *readerImp) NamedFunc(handle argers.ArgHandle, short, long string) argers.Reader {
	imp.validateNames(`named input`, short, long)
	if imp.shortNamed == nil {
		imp.shortNamed = map[string]argers.ArgHandle{}
		imp.longNamed = map[string]argers.ArgHandle{}
	}
	imp.shortNamed[short] = handle
	imp.longNamed[long] = handle
	return imp
}

func (imp *readerImp) NamedBool(target *bool, short, long string) argers.Reader {
	return Named(imp, target, short, long)
}

func (imp *readerImp) NamedStr(target *string, short, long string) argers.Reader {
	return Named(imp, target, short, long)
}

func (imp *readerImp) NamedInt(target *int, short, long string) argers.Reader {
	return Named(imp, target, short, long)
}

func (imp *readerImp) NamedFloat(target *float64, short, long string) argers.Reader {
	return Named(imp, target, short, long)
}

func (imp *readerImp) PosFunc(handle argers.ArgHandle) argers.Reader {
	if len(imp.optionals) > 0 {
		panic(terror.New(`may not add a new positional argument after an optional argument has been added`))
	}
	if imp.variant != nil {
		panic(terror.New(`may not add a new positional argument after a variant argument has been added`))
	}
	imp.pos = append(imp.pos, handle)
	return imp
}

func (imp *readerImp) PosBool(target *bool) argers.Reader {
	return Pos(imp, target)
}

func (imp *readerImp) PosStr(target *string) argers.Reader {
	return Pos(imp, target)
}

func (imp *readerImp) PosInt(target *int) argers.Reader {
	return Pos(imp, target)
}

func (imp *readerImp) PosFloat(target *float64) argers.Reader {
	return Pos(imp, target)
}

func (imp *readerImp) OptionalFunc(handle argers.ArgHandle) argers.Reader {
	if imp.variant != nil {
		panic(terror.New(`may not add an optional argument after a variant argument has been added`))
	}
	imp.optionals = append(imp.optionals, handle)
	return imp
}

func (imp *readerImp) OptionalBool(target *bool) argers.Reader {
	return Optional(imp, target)
}

func (imp *readerImp) OptionalStr(target *string) argers.Reader {
	return Optional(imp, target)
}

func (imp *readerImp) OptionalInt(target *int) argers.Reader {
	return Optional(imp, target)
}

func (imp *readerImp) OptionalFloat(target *float64) argers.Reader {
	return Optional(imp, target)
}

func (imp *readerImp) VarFunc(handle argers.VarHandle) argers.Reader {
	if len(imp.optionals) > 0 {
		panic(terror.New(`may not add a variant argument after an optional argument has been added`))
	}
	if imp.variant != nil {
		panic(terror.New(`may not add a new variant argument after a variant argument has already been added`))
	}
	imp.variant = handle
	return imp
}

func (imp *readerImp) VarBool(target *[]bool) argers.Reader {
	return Var(imp, target)
}

func (imp *readerImp) VarStr(target *[]string) argers.Reader {
	return Var(imp, target)
}

func (imp *readerImp) VarInt(target *[]int) argers.Reader {
	return Var(imp, target)
}

func (imp *readerImp) VarFloat(target *[]float64) argers.Reader {
	return Var(imp, target)
}

func (imp *readerImp) Struct(target any) argers.Reader {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		panic(terror.New(`must provide a non-nil pointer to a structure`))
	}
	s := v.Elem()
	st := s.Type()
	for i, count := 0, s.NumField(); i < count; i++ {
		imp.structField(s.Field(i), st.Field(i))
	}
	return imp
}

func (imp *readerImp) structField(f reflect.Value, ft reflect.StructField) {
	if !ft.IsExported() {
		return
	}
	tag := strings.TrimSpace(ft.Tag.Get(`args`))
	if tag == `skip` {
		return
	}

	ftt, name := ft.Type, ft.Name
	var target any
	if ftt.Kind() == reflect.Pointer {
		ftt = ftt.Elem()
		target = f.Interface()
	} else {
		target = f.Addr().Interface()
	}

	isVar := false
	if ftt.Kind() == reflect.Slice {
		isVar = true
		ftt = ftt.Elem()
	}

	switch ftt.Kind() {
	case reflect.Bool:
		addStructField[bool](imp, target, name, isVar, tag)
	case reflect.String:
		addStructField[string](imp, target, name, isVar, tag)
	case reflect.Int:
		addStructField[int](imp, target, name, isVar, tag)
	case reflect.Int8:
		addStructField[int8](imp, target, name, isVar, tag)
	case reflect.Int16:
		addStructField[int16](imp, target, name, isVar, tag)
	case reflect.Int32:
		addStructField[int32](imp, target, name, isVar, tag)
	case reflect.Int64:
		addStructField[int64](imp, target, name, isVar, tag)
	case reflect.Uint:
		addStructField[uint](imp, target, name, isVar, tag)
	case reflect.Uint8:
		addStructField[uint8](imp, target, name, isVar, tag)
	case reflect.Uint16:
		addStructField[uint16](imp, target, name, isVar, tag)
	case reflect.Uint32:
		addStructField[uint32](imp, target, name, isVar, tag)
	case reflect.Uint64:
		addStructField[uint64](imp, target, name, isVar, tag)
	case reflect.Float32:
		addStructField[float32](imp, target, name, isVar, tag)
	case reflect.Float64:
		addStructField[float64](imp, target, name, isVar, tag)
	case reflect.Complex64:
		addStructField[complex64](imp, target, name, isVar, tag)
	case reflect.Complex128:
		addStructField[complex128](imp, target, name, isVar, tag)
	default:
		panic(terror.New(`unexpected field type for arguments`).
			With(`field`, name))
	}
}

func addStructField[T utils.ParsableConstraint](imp *readerImp, target any, name string, isVar bool, tag string) {
	defer func() {
		if r := recover(); r != nil {
			panic(terror.RecoveredPanic(r).
				With(`tag`, tag).
				With(`field`, name))
		}
	}()

	if isVar {
		if len(tag) > 0 {
			panic(terror.New(`invalid tag on a variadic argument value. May only have the skip tag.`))
		}
		Var[T](imp, target.(*[]T))
		return
	}

	if tag == `optional` {
		Optional[T](imp, target.(*T))
		return
	}

	if hasNames, isFlag, short, long, defVal := parseFieldTag[T](tag); hasNames {
		if isFlag {
			Flag[T](imp, target.(*T), defVal, short, long)
			return
		}
		Named[T](imp, target.(*T), short, long)
		return
	}

	Pos[T](imp, target.(*T))
}

func parseFieldTag[T utils.ParsableConstraint](tag string) (hasNames, isFlag bool, short, long string, defVal T) {
	var tags []string
	if len(tag) > 0 {
		tags = enumerator.Split(tag, `,`).Trim().ToSlice()
	}
	count := len(tags)
	if count <= 0 {
		return false, false, ``, ``, utils.Zero[T]()
	}

	if tags[0] == `flag` {
		if count < 3 || count > 4 {
			panic(terror.New(`the tag on a flag must have three or four values, i.e. "flag,v,verbose"`))
		}

		if count == 4 {
			var err error
			defVal, err = utils.Parse[T](tags[3])
			if err != nil {
				panic(terror.New(`the default value in the tag for a flag could not be parsed`, err))
			}
		} else {
			if utils.TypeOf[T]().Kind() != reflect.Bool {
				panic(terror.New(`the tag on a flag must have the fourth value ` +
					`to use as a default value when the type is not a bool`))
			}
			defVal = any(true).(T)
		}

		return true, true, tags[1], tags[2], defVal
	}

	if count != 2 {
		panic(terror.New(`the tag on a named input must have two values, i.e. "i,input"`))
	}
	return true, false, tags[0], tags[1], defVal
}

func (imp *readerImp) Process(args []string) error {
	argList := list.With(args...)
	if err := imp.consumeNamedInput(argList); err != nil {
		return err
	}
	if err := imp.consumePosArgs(argList); err != nil {
		return err
	}
	if err := imp.consumeVarArgs(argList); err != nil {
		return err
	}
	return nil
}

func (imp *readerImp) consumeNamedInput(argList collections.List[string]) error {
	for i := 0; i < argList.Count(); i++ {
		arg := argList.Get(i)
		if shortArgMatch(arg) {
			if err := imp.consumeShortNamedInput(argList, i, arg); err != nil {
				return err
			}
			argList.Remove(i, 1)
			i--
			continue
		}
		if longArgMatch(arg) {
			if err := imp.consumeLongNamedInput(argList, i, arg); err != nil {
				return err
			}
			argList.Remove(i, 1)
			i--
			continue
		}
	}
	return nil
}

func (imp *readerImp) consumeShortNamedInput(argList collections.List[string], index int, arg string) error {
	maxIndex := len(arg) - 1
	for i := 1; i <= maxIndex; i++ {
		short := string(arg[i])

		if handle, has := imp.shortFlags[short]; has {
			if err := handle(); err != nil {
				return terror.New(`error setting flag`, err).
					With(`name`, short).
					With(`name index`, i-1).
					With(`argument`, arg).
					With(`argument index`, index)
			}
			continue
		}

		if named, has := imp.shortNamed[short]; has {
			if i != maxIndex {
				return terror.New(`may not have a short named value anywhere but the end of a flag group`).
					With(`name`, short).
					With(`name index`, i-1).
					With(`argument`, arg).
					With(`argument index`, index)
			}

			if index+1 >= argList.Count() {
				return terror.New(`no value found for short named value at end of arguments`).
					With(`name`, short).
					With(`name index`, i-1).
					With(`argument`, arg).
					With(`argument index`, index)
			}

			value := argList.Get(index + 1)
			argList.Remove(index+1, 1)
			if err := named(value); err != nil {
				return terror.New(`error setting named argument`, err).
					With(`name`, short).
					With(`name index`, i-1).
					With(`argument`, arg).
					With(`argument index`, index)
			}
			break
		}

		return terror.New(`unknown short name found in arguments`).
			With(`name`, short).
			With(`name index`, i-1).
			With(`argument`, arg).
			With(`argument index`, index)
	}
	return nil
}

func (imp *readerImp) consumeLongNamedInput(argList collections.List[string], index int, arg string) error {
	long := arg[2:]

	if handle, has := imp.longFlags[long]; has {
		if err := handle(); err != nil {
			return terror.New(`error setting flag`, err).
				With(`argument`, arg).
				With(`argument index`, index)
		}
		return nil
	}

	if named, has := imp.longNamed[long]; has {

		if index+1 >= argList.Count() {
			return terror.New(`no value found for long named value at end of arguments`).
				With(`argument`, arg).
				With(`argument index`, index)
		}

		value := argList.Get(index + 1)
		argList.Remove(index+1, 1)
		if err := named(value); err != nil {
			return terror.New(`error setting named argument`, err).
				With(`argument`, arg).
				With(`argument index`, index)
		}
		return nil
	}

	return terror.New(`unknown long name found in arguments`).
		With(`argument`, arg).
		With(`argument index`, index)
}

func (imp *readerImp) consumePosArgs(argList collections.List[string]) error {
	posCount := len(imp.pos)
	if posCount <= 0 {
		return nil
	}

	if argList.Count() < posCount {
		return terror.New(`not enough positional arguments`).
			With(`needed`, posCount).
			With(`gotten`, argList.Count()).
			With(`arguments`, argList.String())
	}

	for i, arg := range argList.TakeFront(posCount).ToSlice() {
		if err := imp.pos[i](arg); err != nil {
			return terror.New(`error setting positional argument`, err).
				With(`argument`, arg).
				With(`argument index`, posCount)
		}
	}
	return nil
}

func (imp *readerImp) consumeVarArgs(argList collections.List[string]) error {
	if argList.Empty() {
		return nil
	}

	if imp.variant != nil {
		if err := imp.variant(argList.ToSlice()); err != nil {
			return terror.New(`error setting variant argument`, err).
				With(`arguments`, argList)
		}
		return nil
	}

	optCount := len(imp.optionals)
	if argList.Count() > optCount {
		return terror.New(`too many optional arguments`).
			With(`maximum`, optCount).
			With(`gotten`, argList.Count()).
			With(`arguments`, argList.String())
	}

	for i, arg := range argList.ToSlice() {
		if err := imp.optionals[i](arg); err != nil {
			return terror.New(`error setting optional argument`, err).
				With(`arguments`, argList).
				With(`argument index`, i+len(imp.pos))
		}
	}
	return nil
}
