package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/graphics"
)

// Color is a 32-bit ARGB color.
type Color uint32

// toByte converts the real number between 0.0 to 1.0 into
// a byte value for that number between 0 and 255 then shifted by given amount.
func toByte(value float64, shift int) uint32 {
	return graphics.Clamp(uint32(value*255.0), 0, 255) << shift
}

// fromByte converts a byte value at the given shift offset from 0 to 255
// into a real number between 0.0 to 1.0.
func fromByte(value uint32, shift int) float64 {
	return float64((value>>shift)&0xFF) / 255.0
}

// New creates a new color with the given components.
//
// The components are an optional alpha, red, green, then blue.
// The components will be clamped between 0.0 and 1.0 and granulated to 32-bits.
//
// There must be 3 or 4 components, otherwise this will panic.
// If only 3 components are given, alpha is set to opaque.
func New(components ...float64) Color {
	switch len(components) {
	case 3:
		return Color(uint32(Black) | toByte(components[0], 16) | toByte(components[1], 8) | toByte(components[2], 0))
	case 4:
		return Color(toByte(components[0], 24) | toByte(components[1], 16) | toByte(components[2], 8) | toByte(components[3], 0))
	default:
		panic(fmt.Errorf(`too many components given to New. Must be 3 or 4 but got %d`, len(components)))
	}
}

// FromHSV creates a color from hue, saturation, and value.
//
// The hue is in degrees wrapped between 0.0 to 360.0.
// Saturation and value is between 0.0 and 100.0.
// Alpha is between 0.0 and 1.0.
//
// This may only accept zero or one alpha value, otherwise this will panic.
// If an alpha isn't given, alpha is set to opaque.
func FromHSV(hue, saturation, value float64, alpha ...float64) Color {
	a := 1.0
	if len(alpha) == 1 {
		a = alpha[0]
	} else if len(alpha) > 0 {
		panic(fmt.Errorf(`too many alpha given to FromHSV. Must be 0 or 1 but got %d`, len(alpha)))
	}

	hue = graphics.Wrap(hue, 0.0, 360.0)
	index, f := math.Modf(hue / 60.0) // sector 0 to 5
	saturation = graphics.Clamp(saturation/100.0, 0.0, 1.0)
	value = graphics.Clamp(value/100.0, 0.0, 1.0)
	p := value * (1.0 - saturation)
	q := value * (1.0 - saturation*f)
	t := value * (1.0 - saturation*(1.0-f))

	switch index {
	case 0:
		return New(a, value, t, p)
	case 1:
		return New(a, q, value, p)
	case 2:
		return New(a, p, value, t)
	case 3:
		return New(a, p, q, value)
	case 4:
		return New(a, t, p, value)
	default:
		return New(a, value, p, q)
	}
}

// FromHSL creates a color from hue, saturation, and lighting.
//
// The hue is in degrees wrapped between 0.0 to 360.0.
// Saturation and lighting is between 0.0 and 100.0.
// Alpha is between 0.0 and 1.0.
func FromHSL(hue, saturation, lightness float64, alpha ...float64) Color {
	a := 1.0
	if len(alpha) == 1 {
		a = alpha[0]
	} else if len(alpha) > 0 {
		panic(fmt.Errorf(`too many alpha given to FromHSL. Must be 0 or 1 but got %d`, len(alpha)))
	}

	hue = graphics.Wrap(hue/360.0, 0.0, 1.0)
	saturation = graphics.Clamp(saturation/100.0, 0.0, 1.0)
	lightness = graphics.Clamp(lightness/100.0, 0.0, 1.0)

	if saturation <= 0.0 {
		return New(a, lightness, lightness, lightness)
	}

	var max float64
	if lightness < 0.5 {
		max = lightness * (1.0 + saturation)
	} else {
		max = lightness + saturation - lightness*saturation
	}
	min := 2.0*lightness - max
	width := (max - min) * 6.0

	third := 1.0 / 3.0
	f := func(t float64) float64 {
		if t < 0.0 {
			t += 1.0
		}
		if t > 1.0 {
			t -= 1.0
		}
		if t*6.0 < 1.0 {
			return min + width*t
		}
		if t < 0.5 {
			return max
		}
		if t*3.0 < 2.0 {
			return min + width*(2.0/3.0-t)
		}
		return min
	}
	return New(a, f(hue+third), f(hue), f(hue-third))
}

// String gets the 0xAARRGGBB formatted string.
func (c Color) String() string {
	s := strings.ToUpper(strconv.FormatUint(uint64(c), 16))
	if count := len(s); count < 8 {
		s = strings.Repeat(`0`, 8-count) + s
	}
	return `0x` + s
}

// Name gets the name for this color if the color is named.
func (c Color) Name() (string, bool) {
	name, hasName := toName[c]
	return name, hasName
}

// Alpha component between 0.0 (transparent) and 1.0 (opaque) inclusively.
func (c Color) Alpha() float64 {
	return fromByte(uint32(c), 24)
}

// Red component between 0.0 and 1.0 inclusively.
func (c Color) Red() float64 {
	return fromByte(uint32(c), 16)
}

// Green component between 0.0 and 1.0 inclusively.
func (c Color) Green() float64 {
	return fromByte(uint32(c), 8)
}

// Blue component between 0.0 and 1.0 inclusively.
func (c Color) Blue() float64 {
	return fromByte(uint32(c), 0)
}

// SetAlpha creates a new color the same as this color with the alpha changed.
// The alpha component between 0.0 (transparent) and 1.0 (opaque) inclusively.
func (c Color) SetAlpha(alpha float64) Color {
	return Color(uint32(c)&0x00FFFFFF | toByte(alpha, 24))
}

// SetRed creates a new color the same as this color with the red changed.
// The red component between 0.0 and 1.0 inclusively.
func (c Color) SetRed(red float64) Color {
	return Color(uint32(c)&0xFF00FFFF | toByte(red, 16))
}

// SetGreen creates a new color the same as this color with the green changed.
// The green component between 0.0 and 1.0 inclusively.
func (c Color) SetGreen(green float64) Color {
	return Color(uint32(c)&0xFFFF00FF | toByte(green, 8))
}

// SetBlue creates a new color the same as this color with the blue changed.
// The blue component between 0.0 and 1.0 inclusively.
func (c Color) SetBlue(blue float64) Color {
	return Color(uint32(c)&0xFFFFFF00 | toByte(blue, 0))
}

// Slice gets the color components in the order red, green, blue, then alpha.
func (c Color) Slice() []float64 {
	return []float64{c.Red(), c.Green(), c.Blue(), c.Alpha()}
}

// Invert will invert the color, creating the complement color and inverted translucency.
func (c Color) Invert() Color {
	return Color(^uint(c))
}

// Lerp creates the linear interpolation between this color and the `other` color.
//
// The `i` is interpolation factor. 0.0 or less will return this color.
// 1.0 or more will return the `other` color. Between 0.0 and 1.0 will be
// a scaled mixture of the two colors.
func (c Color) Lerp(other Color, i float64) Color {
	return New(
		graphics.Lerp(c.Alpha(), other.Alpha(), i),
		graphics.Lerp(c.Red(), other.Red(), i),
		graphics.Lerp(c.Green(), other.Green(), i),
		graphics.Lerp(c.Blue(), other.Blue(), i))
}

// Add creates a new color as the sum of this color and the `other` color.
//
// The color components will saturate at 1.0 so are limited to 1.0.
func (c Color) Add(other Color) Color {
	return New(
		c.Alpha()+other.Alpha(),
		c.Red()+other.Red(),
		c.Green()+other.Green(),
		c.Blue()+other.Blue())
}

// Sub creates a new color as the difference of this color and the `other` color.
//
// The color components will deplete at 0.0 so are limited to 0.0.
func (c Color) Sub(other Color) Color {
	return New(
		c.Alpha()-other.Alpha(),
		c.Red()-other.Red(),
		c.Green()-other.Green(),
		c.Blue()-other.Blue())
}

// Scale creates a new color scaled by the given `scalar`.
func (c Color) Scale(scalar float64) Color {
	return New(
		c.Alpha()*scalar,
		c.Red()*scalar,
		c.Green()*scalar,
		c.Blue()*scalar)
}

// Equals returns true if this object and the given object are equal.
func (c Color) Equals(other any) bool {
	c2, ok := other.(Color)
	return ok && uint32(c) == uint32(c2)
}
