package args

import (
	"github.com/Snow-Gremlin/goToolbox/argers"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new argument reader.
func New() argers.Reader {
	return &readerImp{
		shortFlags: nil,
		longFlags:  nil,
		shortNamed: nil,
		longNamed:  nil,
		pos:        nil,
		optionals:  nil,
		variant:    nil,
	}
}

// Flag adds a flag to the given reader.
//
// When the flag by the short name (e.g. `-v`) or long name (e.g. `--verbose`)
// is in the arguments, the given target will be assigned to the given value,
// otherwise the target is not modified.
// Flags may be grouped in the arguments (e.g. `-bvf`).
//
// The returned reader is the same as the given reader.
func Flag[T any](r argers.Reader, target *T, value T, short, long string) argers.Reader {
	if target == nil {
		panic(terror.New(`must provide a non-nil target pointer for an argument flag`))
	}
	return r.FlagFunc(func() error {
		*target = value
		return nil
	}, short, long)
}

// Named adds a named value with a value to the given reader.
//
// When the named value by the short name (e.g. `-o file.txt`) or long name
// (e.g. `--out file.txt`) is in the arguments, the given target will be
// assigned to the value of the following argument.
// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
//
// The returned reader is the same as the given reader.
func Named[T utils.ParsableConstraint](r argers.Reader, target *T, short, long string) argers.Reader {
	if target == nil {
		panic(terror.New(`must provide a non-nil target pointer for a named argument`))
	}
	return r.NamedFunc(func(s string) error {
		value, err := utils.Parse[T](s)
		if err != nil {
			return err
		}
		*target = value
		return nil
	}, short, long)
}

// Pos adds a positional argument to the given reader.
//
// After all the flags and named values have been removed, the remaining
// arguments are read in order where the first added positional argument
// then the next positional and so on.
// The given target is set to the given argument at its position.
//
// The returned reader is the same as the given reader.
func Pos[T utils.ParsableConstraint](r argers.Reader, target *T) argers.Reader {
	if target == nil {
		panic(terror.New(`must provide a non-nil target pointer for a positional argument`))
	}
	return r.PosFunc(func(s string) error {
		value, err := utils.Parse[T](s)
		if err != nil {
			return err
		}
		*target = value
		return nil
	})
}

// Optional adds an optional argument to the given reader.
//
// After all the flags, named values, and positional arguments are read,
// any remaining arguments will be set to the given optional in the order
// that they were added.
// Optional arguments may only be added after positional arguments and may
// not be used with variant arguments.
// The given target is set to the given argument at its position.
//
// The returned reader is the same as the given reader.
func Optional[T utils.ParsableConstraint](r argers.Reader, target *T) argers.Reader {
	if target == nil {
		panic(terror.New(`must provide a non-nil target pointer for an optional argument`))
	}
	return r.OptionalFunc(func(s string) error {
		value, err := utils.Parse[T](s)
		if err != nil {
			return err
		}
		*target = value
		return nil
	})
}

// Var adds a variant argument to the given reader.
//
// After all the flags, named values, and positional arguments are read,
// any remaining arguments will be set to this variant argument.
// Only one variant may be added and it may only be added after
// positional arguments. This may not be used with optional arguments.
// The target is set to the given arguments.
//
// The returned reader is the same as the given reader.
func Var[T utils.ParsableConstraint, S ~[]T](r argers.Reader, target *S) argers.Reader {
	if target == nil {
		panic(terror.New(`must provide a non-nil target pointer for a variadic argument`))
	}
	return r.VarFunc(func(strings []string) error {
		values := make(S, len(strings))
		for i, s := range strings {
			value, err := utils.Parse[T](s)
			if err != nil {
				return err
			}
			values[i] = value
		}
		*target = values
		return nil
	})
}
