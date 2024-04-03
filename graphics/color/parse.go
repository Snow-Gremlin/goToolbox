package color

import (
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

const (
	hashCap = `^#[0-9a-fA-F]{3,8}$`
	hexCap  = `\s*([0-9a-fA-F]+)\s*`
	rgbCap  = `^rgb\s*\(` + hexCap + `,` + hexCap + `,` + hexCap + `\)$`
	rgbaCap = `^rgba\s*\(` + hexCap + `,` + hexCap + `,` + hexCap + `,\s*([0-9.]+)\s*\)$`
)

var (
	hashMatch = utils.LazyRegex(hashCap)
	rgbMatch  = utils.LazyRegex(rgbCap)
	rgbaMatch = utils.LazyRegex(rgbaCap)
)

// Parse will determine the color via the given string.
func Parse(input string) (Color, error) {
	input = strings.TrimSpace(input)

	// Try for `#` hex numbers
	if hashMatch().MatchString(input) {
		return parseHashColor(input)
	}

	// Try for `rgb(RR, GG, BB)`
	if match := rgbMatch().FindStringSubmatch(input); len(match) == 4 {
		return parseComponents(input, ``, match[1], match[2], match[3])
	}

	// Try for `rgba(RR, GG, BB, alpha)`
	if match := rgbaMatch().FindStringSubmatch(input); len(match) == 5 {
		return parseComponents(input, match[4], match[1], match[2], match[3])
	}

	// TODO: add hsl(120, 100%, 50%) and hsla(120, 100%, 50%, 0.3)

	// Try for named color
	low := strings.ReplaceAll(input, ` `, ``)
	low = strings.ReplaceAll(low, `_`, ``)
	low = strings.ToLower(low)
	if color, has := fromLowerCase[low]; has {
		return color, nil
	}

	// Try for hex number without `#`, e.g. 0xAARRGGBB
	low, _ = strings.CutPrefix(low, `0x`)
	if value64, err := strconv.ParseUint(low, 16, 32); err == nil {
		value := uint32(value64)
		if len(low) <= 6 {
			value |= uint32(Black)
		}
		return Color(value), nil
	}

	return Transparent, terror.New(`unable to parse color for unknown format`).
		With(`input`, input)
}

func parseHashColor(input string) (Color, error) {
	hash := input[1:]
	if value, err := strconv.ParseUint(hash, 16, 32); err == nil {
		switch len(hash) {
		case 3: // #RGB
			value = (value&0xF00)<<8 | (value&0x0F0)<<4 | value&0x00F
			return Color(uint64(Black) | value | value<<4), nil
		case 6: // #RRGGBB
			return Color(uint64(Black) | value), nil
		case 8: // #RRGGBBAA
			return Color((value&0xFF)<<24 | (value&0xFFFFFF00)>>8), nil
		}
	}
	return Transparent, terror.New(`unable to parse hash color. Expected 3, 6, or 8 hexadecimal digits`).
		With(`input`, input)
}

func parseComponents(input, alpha, red, green, blue string) (Color, error) {
	var value uint32
	if len(alpha) > 0 {
		a, err := strconv.ParseFloat(alpha, 64)
		if err != nil {
			return Transparent, terror.New(`alpha component not parsable. Expected floating point value between 0.0 to 1.0`).
				With(`input`, input).
				With(`alpha`, alpha)
		}
		value |= toByte(a, 24)
	} else {
		value |= uint32(Black)
	}

	if len(red) > 0 && red != `0` {
		c, err := strconv.ParseUint(red, 10, 8)
		if err != nil {
			return Transparent, terror.New(`red component not parsable. Expected base-10 value between 0 and 255`).
				With(`input`, input).
				With(`red`, red)
		}
		value |= uint32(c&0xFF) << 16
	}

	if len(green) > 0 && green != `0` {
		c, err := strconv.ParseUint(green, 10, 8)
		if err != nil {
			return Transparent, terror.New(`green component not parsable. Expected base-10 value between 0 and 255`).
				With(`input`, input).
				With(`green`, green)
		}
		value |= uint32(c&0xFF) << 8
	}

	if len(blue) > 0 && blue != `0` {
		c, err := strconv.ParseUint(blue, 10, 8)
		if err != nil {
			return Transparent, terror.New(`blue component not parsable. Expected base-10 value between 0 and 255`).
				With(`input`, input).
				With(`blue`, blue)
		}
		value |= uint32(c & 0xFF)
	}

	return Color(value), nil
}
