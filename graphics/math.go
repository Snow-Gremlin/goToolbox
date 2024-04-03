package graphics

import (
	"cmp"
	"math"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Clamp will keep the given value between min and max.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value <= min {
		return min
	}
	if value >= max {
		return max
	}
	return value
}

// Mod gets the modulo of the given value and the given divisor.
// The remainder after the division is returned.
func Mod[T utils.NumConstraint](value, divisor T) T {
	switch v := any(value).(type) {
	case int:
		return T(v % int(divisor))
	case int8:
		return T(v % int8(divisor))
	case int16:
		return T(v % int16(divisor))
	case int32:
		return T(v % int32(divisor))
	case int64:
		return T(v % int64(divisor))
	case uint:
		return T(v % uint(divisor))
	case uint8:
		return T(v % uint8(divisor))
	case uint16:
		return T(v % uint16(divisor))
	case uint32:
		return T(v % uint32(divisor))
	case uint64:
		return T(v % uint64(divisor))
	case float32:
		return T(float32(math.Mod(float64(v), float64(divisor))))
	case float64:
		return T(math.Mod(v, float64(divisor)))
	default:
		return utils.Zero[T]()
	}
}

// Wrap will keep wrap the given value into the given min and max.
// For example if the range is 10 to 20: given 15 returns 15,
// given 24 returns 14, given 6 returns 16.
func Wrap[T utils.NumConstraint](value, min, max T) T {
	return Mod(value-min, max-min) + min
}

// Lerp gets the linear interpolation between `a` and `b` as scaled by `i`.
//
// The `i` is interpolation factor. 0.0 or less will return `a`.
// 1.0 or more will return `b`. Between 0.0 and 1.0 will be
// a scaled mixture of the two values.
func Lerp(a, b, i float64) float64 {
	if i <= 0.0 {
		return a
	}
	if i >= 1.0 {
		return b
	}
	return a*i + b*(1.0-i)
}
