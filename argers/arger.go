package argers

type (
	// FlagHandle is the handle for flag arguments.
	FlagHandle func() error

	// ArgHandle is the handle for named, positional, and optional arguments.
	ArgHandle func(value string) error

	// VarHandle is the handle for a variant argument.
	VarHandle func(value []string) error
)

// Reader is an argument reader.
//
// This helps parsing arguments for complicated applications that
// may have several tools built into one application.
type Reader interface {
	interface {
		// FlagFunc adds a flag to this reader.
		//
		// When the flag by the short name (e.g. `-v`) or long name (e.g. `--verbose`)
		// is in the arguments, the given handle will be called.
		// Flags may be grouped in the arguments (e.g. `-bvf`).
		//
		// The returned reader is the receiver.
		FlagFunc(handle FlagHandle, short, long string) Reader

		// Flag adds a bool flag to this reader.
		//
		// When the flag by the short name (e.g. `-v`) or long name (e.g. `--verbose`)
		// is in the arguments, the given target will be set to true,
		// otherwise the target is not modified.
		// Flags may be grouped in the arguments (e.g. `-bvf`).
		//
		// The returned reader is the receiver.
		Flag(target *bool, short, long string) Reader

		// FlagStr adds a string flag to this reader.
		//
		// When the flag by the short name (e.g. `-f`) or long name (e.g. `--fast`)
		// is in the arguments, the given target will be set to the given value,
		// otherwise the target is not modified.
		// This can be used to have several flags with different meanings
		// that target the same string.
		// Flags may be grouped in the arguments (e.g. `-bvf`).
		//
		// The returned reader is the receiver.
		FlagStr(target *string, value, short, long string) Reader

		// FlagStr adds a int flag to this reader.
		//
		// When the flag by the short name (e.g. `-f`) or long name (e.g. `--fast`)
		// is in the arguments, the given target will be set to the given value,
		// otherwise the target is not modified.
		// This can be used to have several flags with different meanings
		// that target the same string.
		// Flags may be grouped in the arguments (e.g. `-bvf`).
		//
		// The returned reader is the receiver.
		FlagInt(target *int, value int, short, long string) Reader
	}

	interface {
		// NamedFunc adds a named value with a value to this reader.
		//
		// When the named value by the short name (e.g. `-o file.txt`) or long name
		// (e.g. `--out file.txt`) is in the arguments, the given handle will be
		// called with the value of the following argument.
		// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
		//
		// The returned reader is the receiver.
		NamedFunc(handle ArgHandle, short, long string) Reader

		// NamedBool adds a bool named value with a value to this reader.
		//
		// When the named value by the short name (e.g. `-o file.txt`) or long name
		// (e.g. `--out file.txt`) is in the arguments, the given bool will be
		// set to the value of the following argument.
		// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
		//
		// The returned reader is the receiver.
		NamedBool(target *bool, short, long string) Reader

		// NamedStr adds a string named value with a value to this reader.
		//
		// When the named value by the short name (e.g. `-o file.txt`) or long name
		// (e.g. `--out file.txt`) is in the arguments, the given string will be
		// set to the value of the following argument.
		// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
		//
		// The returned reader is the receiver.
		NamedStr(target *string, short, long string) Reader

		// NamedInt adds an int named value with a value to this reader.
		//
		// When the named value by the short name (e.g. `-o file.txt`) or long name
		// (e.g. `--out file.txt`) is in the arguments, the given int will be
		// set to the value of the following argument.
		// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
		//
		// The returned reader is the receiver.
		NamedInt(target *int, short, long string) Reader

		// NamedFloat adds a float named value with a value to this reader.
		//
		// When the named value by the short name (e.g. `-o file.txt`) or long name
		// (e.g. `--out file.txt`) is in the arguments, the given float will be
		// set to the value of the following argument.
		// Flags must be the last in a group (e.g. `-vo file.txt`) or not in a group.
		//
		// The returned reader is the receiver.
		NamedFloat(target *float64, short, long string) Reader
	}

	interface {
		// PosFunc adds a positional argument to this reader.
		//
		// After all the flags and named values have been removed, the remaining
		// arguments are read in order where the first added positional argument
		// then the next positional and so on.
		// The given handle is called with the given argument at its position.
		//
		// The returned reader is the receiver.
		PosFunc(handle ArgHandle) Reader

		// PosBool adds a bool positional argument to this reader.
		//
		// After all the flags and named values have been removed, the remaining
		// arguments are read in order where the first added positional argument
		// then the next positional and so on.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		PosBool(target *bool) Reader

		// PosStr adds a string positional argument to this reader.
		//
		// After all the flags and named values have been removed, the remaining
		// arguments are read in order where the first added positional argument
		// then the next positional and so on.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		PosStr(target *string) Reader

		// PosInt adds an int positional argument to this reader.
		//
		// After all the flags and named values have been removed, the remaining
		// arguments are read in order where the first added positional argument
		// then the next positional and so on.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		PosInt(target *int) Reader

		// PosFloat adds a float positional argument to this reader.
		//
		// After all the flags and named values have been removed, the remaining
		// arguments are read in order where the first added positional argument
		// then the next positional and so on.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		PosFloat(target *float64) Reader
	}

	interface {
		// OptionalFunc adds an optional argument to this reader.
		//
		// After all the flags, named values, and positional arguments
		// are read, any remaining arguments will be set to the given
		// optional in the order that they were added.
		// Optional arguments may only be added after positional arguments
		// and may not be used with variant arguments.
		// The given handle is called with the given argument at its position.
		//
		// The returned reader is the receiver.
		OptionalFunc(handle ArgHandle) Reader

		// OptionalBool adds a bool optional argument to this reader.
		//
		// After all the flags, named values, and positional arguments
		// are read, any remaining arguments will be set to the given
		// optional in the order that they were added.
		// Optional arguments may only be added after positional arguments
		// and may not be used with variant arguments.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		OptionalBool(target *bool) Reader

		// OptionalStr adds a string optional argument to this reader.
		//
		// After all the flags, named values, and positional arguments
		// are read, any remaining arguments will be set to the given
		// optional in the order that they were added.
		// Optional arguments may only be added after positional arguments
		// and may not be used with variant arguments.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		OptionalStr(target *string) Reader

		// OptionalInt adds an int optional argument to this reader.
		//
		// After all the flags, named values, and positional arguments
		// are read, any remaining arguments will be set to the given
		// optional in the order that they were added.
		// Optional arguments may only be added after positional arguments
		// and may not be used with variant arguments.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		OptionalInt(target *int) Reader

		// OptionalFloat adds a float optional argument to this reader.
		//
		// After all the flags, named values, and positional arguments
		// are read, any remaining arguments will be set to the given
		// optional in the order that they were added.
		// Optional arguments may only be added after positional arguments
		// and may not be used with variant arguments.
		// The given target is set to the given argument at its position.
		//
		// The returned reader is the receiver.
		OptionalFloat(target *float64) Reader
	}

	interface {
		// VarFunc adds a variant argument to this reader.
		//
		// After all the flags, named values, and positional arguments are read,
		// any remaining arguments will be set to this variant argument.
		// Only one variant may be added and it may only be added after
		// positional arguments. This may not be used with optional arguments.
		// The given handle is called with the given argument.
		//
		// The returned reader is the receiver.
		VarFunc(handle VarHandle) Reader

		// VarBool adds a bool variant argument to this reader.
		//
		// After all the flags, named values, and positional arguments are read,
		// any remaining arguments will be set to this variant argument.
		// Only one variant may be added and it may only be added after
		// positional arguments. This may not be used with optional arguments.
		// The given target is set to the given argument.
		//
		// The returned reader is the receiver.
		VarBool(target *[]bool) Reader

		// VarStr adds a string variant argument to this reader.
		//
		// After all the flags, named values, and positional arguments are read,
		// any remaining arguments will be set to this variant argument.
		// Only one variant may be added and it may only be added after
		// positional arguments. This may not be used with optional arguments.
		// The given target is set to the given argument.
		//
		// The returned reader is the receiver.
		VarStr(target *[]string) Reader

		// VarInt adds an int variant argument to this reader.
		//
		// After all the flags, named values, and positional arguments are read,
		// any remaining arguments will be set to this variant argument.
		// Only one variant may be added and it may only be added after
		// positional arguments. This may not be used with optional arguments.
		// The given target is set to the given argument.
		//
		// The returned reader is the receiver.
		VarInt(target *[]int) Reader

		// VarFloat adds a float variant argument to this reader.
		//
		// After all the flags, named values, and positional arguments are read,
		// any remaining arguments will be set to this variant argument.
		// Only one variant may be added and it may only be added after
		// positional arguments. This may not be used with optional arguments.
		// The given target is set to the given argument.
		//
		// The returned reader is the receiver.
		VarFloat(target *[]float64) Reader
	}

	// Struct adds arguments using the fields of the given target to this reader.
	//
	// The given target must be a non-nil pointer to a struct that has fields
	// that the arguments will write to.
	//  - Only public fields will be written to.
	//  - The fields are added in the order they are defined in the struct.
	//  - By default, the public field is a positional argument unless the field
	//    is a slice of values then it is a variadic argument.
	//  - To not add a public field, add the tag, `args:"skip"`.
	//  - To make a field an optional argument, add the tag, `args:"optional"`.
	//	- To make a flag argument, add the tag `args:"flag,<short>,<long>,<value>"`
	//    where `<short>` is the short name (e.g. "v"), `<long>` is
	//    the long name (e.g. "verbose"), and `<value>` is the value to assign
	//    to the field when that flag argument is added.
	//    If the field is a boolean the value may be left off and the value is
	//    defaulted to `true`.
	//  - To make a named input argument, add the tag `args:"<short>,<long>"`
	//    where `<short>` is the short name (e.g. "i") and `<long>` is
	//    the long name (e.g. "input").
	//  - Fields may be a non-nil pointer and still function the same as if
	//    it isn't a pointer.
	//
	// Example:
	//    type argObject struct {
	//		PosOne int
	//      PosTwo string
	//      Input  string `args:"i, input"`
	//      Output string `args:"o, output"`
	//      None   bool   `args:"skip"`
	//      Extra  int    `args:"optional"`
	//    }
	Struct(target any) Reader

	// Process reads the given arguments and calls the appropriate argument
	// handles. Typically this will be given a subset of os.Args().
	// Returns any errors in the arguments.
	Process(args []string) error
}
