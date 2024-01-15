package readonlyVariantList

import (
	"reflect"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	// CountFunc is a function signature to get a count of elements in a collection.
	CountFunc func() int

	// GetFunc is a function signature to get an element in a collection
	// at a given index. The index should only be between zero inclusively,
	// and the count exclusively gotten from the accompanying CountFunc.
	GetFunc[T any] func(int) T

	// OnChangeFunc is a function signature for getting an event that indicates
	// when the source changed.
	//
	// This will be optional since not all data sources know when they
	// have been changed and are able to emit an event.
	OnChangeFunc func() events.Event[collections.ChangeArgs]
)

// From creates a new readonly variant list which is sourced by the two given functions.
//
// The first function returns the count of elements in the list.
// The second function gets an element from the list at the given index.
// The get function shouldn't fail for any index between zero inclusively and length exclusively.
// This returns nil if either function is nil.
func From[T any](count CountFunc, get GetFunc[T], onChange OnChangeFunc) collections.ReadonlyList[T] {
	if count == nil || get == nil {
		return (*readonlyVariantListImp[T])(nil)
	}
	return &readonlyVariantListImp[T]{
		countHandle:    count,
		getHandle:      get,
		onChangeHandle: onChange,
	}
}

// Cast will create a readonly variant list which will perform a cast
// on each value. This will panic if any value can not be cast.
// This may take an optional selector to use instead of a default cast.
func Cast[TOut, TIn any](source collections.ReadonlyList[TIn], selector ...collections.Selector[TIn, TOut]) collections.ReadonlyList[TOut] {
	switch len(selector) {
	case 0:
		return From[TOut](
			source.Count,
			func(i int) TOut { return any(source.Get(i)).(TOut) },
			source.OnChange)
	case 1:
		sel := selector[0]
		return From[TOut](
			source.Count,
			func(i int) TOut { return sel(source.Get(i)) },
			source.OnChange)
	default:
		panic(terror.InvalidArgCount(1, len(selector), `selector`))
	}
}

// Wrap tries to create a new readonly variant list using the giving value
// as the source of data. If the value can not be used, nil is returned.
//
// The value maybe a string, any slice, any array, or from a struct,
// an interface which has `Count() int` and `Get(int) X`, or
// a struct or an interface which has `ToSlice() []X` or `Byte() []byte`.
// Except for the `ToSlice` or `Bytes` structs, which will create a slice and
// use that for the list, the rest will update the elements in the list as the
// underlying source is changed, if it can be changed.
func Wrap(source any) collections.ReadonlyList[any] {
	if utils.IsNil(source) {
		return (*readonlyVariantListImp[any])(nil)
	}

	switch v := source.(type) {
	case string:
		return fromString(v)
	case []byte:
		return fromSlice(v)
	case []int:
		return fromSlice(v)
	case []string:
		return fromSlice(v)
	case []rune:
		return fromSlice(v)
	case []any:
		return fromSlice(v)
	}

	val := reflect.ValueOf(source)
	switch val.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		return fromArrayValue(val)
	case reflect.Map:
		return fromMapValue(val)
	default:
		if list, ok := fromObject(source, val); ok {
			return list
		}
		return fromSingleValue(source)
	}
}

func fromString(str string) collections.ReadonlyList[any] {
	return From(
		func() int { return len(str) },
		func(i int) any { return str[i] },
		nil)
}

func fromSlice[E any, S ~[]E](s S) collections.ReadonlyList[any] {
	return From(
		func() int { return len(s) },
		func(i int) any { return s[i] },
		nil)
}

func fromArrayValue(val reflect.Value) collections.ReadonlyList[any] {
	return From(
		val.Len,
		func(i int) any { return val.Index(i).Interface() },
		nil)
}

func fromObject(source any, val reflect.Value) (collections.ReadonlyList[any], bool) {
	if count := countMethod(source); count != nil {
		if get := reflectGetMethod(val); get != nil {
			return From(count, get, onChange(source)), true
		}
	}

	if v, ok := source.(interface{ Bytes() []byte }); ok {
		return fromSlice(v.Bytes()), true
	}

	if list, ok := fromSliceable(val); ok {
		return list, true
	}

	return nil, false
}

func countMethod(source any) CountFunc {
	switch v := source.(type) {
	case interface{ Count() int }:
		return v.Count
	case interface{ Length() int }:
		return v.Length
	case interface{ Len() int }:
		return v.Len
	default:
		return nil
	}
}

func reflectGetMethod(val reflect.Value) GetFunc[any] {
	if get := val.MethodByName(`Get`); !utils.IsZero(get) {
		t := get.Type()
		if t.NumIn() == 1 && t.In(0) == utils.TypeOf[int]() && t.NumOut() == 1 {
			return func(i int) any {
				result := get.Call([]reflect.Value{reflect.ValueOf(i)})
				return result[0].Interface()
			}
		}
	}
	return nil
}

func onChange(source any) OnChangeFunc {
	if c, ok := source.(collections.OnChanger); ok {
		return c.OnChange
	}
	return nil
}

func fromSliceable(val reflect.Value) (collections.ReadonlyList[any], bool) {
	if toSlice := val.MethodByName(`ToSlice`); !utils.IsZero(toSlice) {
		t := toSlice.Type()
		if t.NumIn() == 0 && t.NumOut() == 1 {
			switch t.Out(0).Kind() {
			case reflect.Array, reflect.Slice, reflect.String:
				slice := toSlice.Call([]reflect.Value{})
				return fromArrayValue(slice[0]), true
			}
		}
	}
	return nil, false
}

func fromMapValue(val reflect.Value) collections.ReadonlyList[any] {
	keys := val.MapKeys()
	return From(
		func() int { return len(keys) },
		func(i int) any {
			key := keys[i].Interface()
			value := val.MapIndex(keys[i]).Interface()
			return tuple2.New(key, value)
		},
		nil)
}

func fromSingleValue(value any) collections.ReadonlyList[any] {
	return From[any](
		func() int { return 1 },
		func(i int) any { return value },
		nil)
}
