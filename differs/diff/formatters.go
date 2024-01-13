package diff

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/differs"
	"github.com/Snow-Gremlin/goToolbox/differs/step"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

const (
	startChange   = `<<<<<<<<`
	middleChange  = `========`
	endChange     = `>>>>>>>>`
	equalPrefix   = ` `
	addedPrefix   = `+`
	removedPrefix = `-`
)

var prefixes = map[step.Step]string{
	step.Equal:   equalPrefix,
	step.Added:   addedPrefix,
	step.Removed: removedPrefix,
}

// EnumerateValues will iterate the given result from a diff of the two given
// slices where the enumeration contains the correct values for each step.
//
// Typically this is used for lines of a document. The number of A and B values
// must be the same as used when creating the given results.
func EnumerateValues[T comparable, S ~[]T](result differs.Result, a, b S) collections.Enumerator[collections.Tuple2[step.Step, S]] {
	if aLen, bLen := len(a), len(b); aLen != result.ACount() || bLen != result.BCount() {
		e := terror.New(`must have the same number of values as the result was created with`).
			With(`required A count`, result.ACount()).
			With(`required B count`, result.BCount())
		if aLen != result.ACount() {
			e = e.With(`incorrect A count`, aLen)
		}
		if bLen != result.BCount() {
			e = e.With(`incorrect B count`, bLen)
		}
		panic(e)
	}
	return enumerator.New(func() collections.Iterator[collections.Tuple2[step.Step, S]] {
		it := result.Enumerate().Iterate()
		aIndex, bIndex := 0, 0
		return iterator.New(func() (collections.Tuple2[step.Step, S], bool) {
			if it.Next() {
				stepType, count := it.Current().Values()
				var values S
				switch stepType {
				case step.Equal:
					values = a[aIndex : aIndex+count]
					aIndex += count
					bIndex += count
				case step.Added:
					values = b[bIndex : bIndex+count]
					bIndex += count
				case step.Removed:
					values = a[aIndex : aIndex+count]
					aIndex += count
				}
				return tuple2.New(stepType, values), true
			}
			return utils.Zero[collections.Tuple2[step.Step, S]](), false
		})
	})
}

// PlusMinus gets the labelled difference between the two slices.
// It formats the results by prepending a "+" to new strings in [b],
// a "-" for any to removed strings from [a], and " " if the strings are the same.
func PlusMinus(path differs.Result, a, b []string) []string {
	result := make([]string, path.Total())
	index := 0
	EnumerateValues(path, a, b).Foreach(func(t collections.Tuple2[step.Step, []string]) {
		stepType, values := t.Values()
		prefix := prefixes[stepType]
		for _, v := range values {
			result[index] = prefix + v
			index++
		}
	})
	return result
}

// Inline gets the difference with all values that are of the same step type
// joined into one string prefixed with a symbol indicating the step type.
func Inline[T comparable, S ~[]T](path differs.Result, a, b S, separator string) []string {
	result := make([]string, path.Count())
	index := 0
	EnumerateValues(path, a, b).Foreach(func(t collections.Tuple2[step.Step, S]) {
		stepType, values := t.Values()
		result[index] = prefixes[stepType] + strings.Join(utils.Strings(values), separator)
		index++
	})
	return result
}

// Merge gets the labelled difference between the two slices
// using a similar output to the git merge differences output.
func Merge(path differs.Result, a, b []string) []string {
	result := make([]string, 0, path.Total()+path.Count()*2+1)
	prevState := step.Equal
	EnumerateValues(path, a, b).Foreach(func(t collections.Tuple2[step.Step, []string]) {
		stepType, values := t.Values()
		switch stepType {
		case step.Equal:
			switch prevState {
			case step.Added:
				result = append(result, endChange)
			case step.Removed:
				result = append(result, middleChange, endChange)
			}

		case step.Added:
			switch prevState {
			case step.Equal:
				result = append(result, startChange, middleChange)
			case step.Removed:
				result = append(result, middleChange)
			}

		case step.Removed:
			switch prevState {
			case step.Equal:
				result = append(result, startChange)
			case step.Added:
				result = append(result, endChange, startChange)
			}
		}
		result = append(result, values...)
		prevState = stepType
	})

	switch prevState {
	case step.Added:
		result = append(result, endChange)
	case step.Removed:
		result = append(result, middleChange, endChange)
	}
	return result
}
