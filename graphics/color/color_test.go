package color

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Color_Parse(t *testing.T) {
	check := func(input string, exp Color) {
		value, err := Parse(input)
		check.NoError(t).With(`input`, input).With(`expected`, exp).Require(err)
		check.Equal(t, exp).With(`input`, input).Assert(value)
	}

	check(`black`, Black)
	check(`0`, Black)
	check(`FF000000`, Black)
	check(`0x000000`, Black)
	check(`0x00_00_00`, Black)
	check(`0xFF000000`, Black)
	check(`rgb(0, 0, 0)`, Black)
	check(`rgba(0, 0, 0, 1.0)`, Black)
	check(`#000`, Black)
	check(`#000000`, Black)
	check(`#000000FF`, Black)

	// Red Component
	check(`red`, Red)
	check(`FF0000`, Red)
	check(`FFFF0000`, Red)
	check(`0xFF0000`, Red)
	check(`0xFF_00_00`, Red)
	check(`0xFFFF0000`, Red)
	check(`rgb(255, 0, 0)`, Red)
	check(`rgba(255, 0, 0, 1.0)`, Red)
	check(`#F00`, Red)
	check(`#FF0000`, Red)
	check(`#FF0000FF`, Red)

	// Green Component
	check(`lime`, Lime)
	check(`FF00`, Lime)
	check(`FF00FF00`, Lime)
	check(`0x00FF00`, Lime)
	check(`0x00_FF_00`, Lime)
	check(`0xFF00FF00`, Lime)
	check(`rgb(0, 255, 0)`, Lime)
	check(`rgba(0, 255, 0, 1.0)`, Lime)
	check(`#0F0`, Lime)
	check(`#00FF00`, Lime)
	check(`#00FF00FF`, Lime)

	// Green Component
	check(`blue`, Blue)
	check(`FF`, Blue)
	check(`FF0000FF`, Blue)
	check(`0x0000FF`, Blue)
	check(`0x00_00_FF`, Blue)
	check(`0xFF0000FF`, Blue)
	check(`rgb(0, 0, 255)`, Blue)
	check(`rgba(0, 0, 255, 1.0)`, Blue)
	check(`#00F`, Blue)
	check(`#0000FF`, Blue)
	check(`#0000FFFF`, Blue)

	check(`darkolivegreen`, DarkOliveGreen)
	check(`DarkOliveGreen`, DarkOliveGreen)
	check(`DARKOLIVEGREEN`, DarkOliveGreen)
	check(`dark olive green`, DarkOliveGreen)
	check(`dark olivE grEEn`, DarkOliveGreen)
	check(`DarkOliveGreen  `, DarkOliveGreen)
	check(`  DarkOliveGreen`, DarkOliveGreen)
	check(`dark_olive_green`, DarkOliveGreen)
	check(`dar_koli_vegre_en`, DarkOliveGreen)

	check(`aqua`, Cyan)
	check(`cyan`, Cyan)
	check(`fuchsia`, Magenta)
	check(`magenta`, Magenta)
}

func Test_Color_FromHSL(t *testing.T) {
	check := func(h, s, l float64, exp Color) {
		c := FromHSL(h, s, l)
		check.Equal(t, exp).
			With(`hue`, h).
			With(`saturation`, s).
			With(`lightness`, l).
			Assert(c)
	}

	check(0.0, 0.0, 0.0, Black)
	check(0.0, 0.0, 100.0, White)
	check(0.0, 100.0, 50.0, Red)
	check(120.0, 100.0, 50.0, Lime)
	check(240.0, 100.0, 50.0, Blue)
	check(60.0, 100.0, 50.0, Yellow)
	check(180.0, 100.0, 50.0, Cyan)
	check(300.0, 100.0, 50.0, Magenta)
	check(0.0, 0.0, 75.0, 0xBFBFBF)
	check(0.0, 0.0, 50.0, Gray)
	check(0.0, 100.0, 25.0, Maroon)
	check(60.0, 100.0, 25.0, Olive)
	check(120.0, 100.0, 25.0, Green)
	check(300.0, 100.0, 25.0, Purple)
	check(180.0, 100.0, 25.0, Teal)
	check(240.0, 100.0, 25.0, Navy)
}
