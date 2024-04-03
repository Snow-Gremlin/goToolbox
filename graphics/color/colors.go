package color

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

// List of predefined colors.
const (
	AliceBlue            = Color(0xFFF0F8FF)
	AntiqueWhite         = Color(0xFFFAEBD7)
	Aquamarine           = Color(0xFF7FFFD4)
	Azure                = Color(0xFFF0FFFF)
	Beige                = Color(0xFFF5F5DC)
	Bisque               = Color(0xFFFFE4C4)
	Black                = Color(0xFF000000)
	BlanchedAlmond       = Color(0xFFFFEBCD)
	Blue                 = Color(0xFF0000FF)
	BlueViolet           = Color(0xFF8A2BE2)
	Brown                = Color(0xFFA52A2A)
	BurlyWood            = Color(0xFFDEB887)
	CadetBlue            = Color(0xFF5F9EA0)
	Chartreuse           = Color(0xFF7FFF00)
	Chocolate            = Color(0xFFD2691E)
	Coral                = Color(0xFFFF7F50)
	CornflowerBlue       = Color(0xFF6495ED)
	Cornsilk             = Color(0xFFFFF8DC)
	Crimson              = Color(0xFFDC143C)
	Cyan                 = Color(0xFF00FFFF)
	DarkBlue             = Color(0xFF00008B)
	DarkCyan             = Color(0xFF008B8B)
	DarkGoldenrod        = Color(0xFFB8860B)
	DarkGray             = Color(0xFFA9A9A9)
	DarkGreen            = Color(0xFF006400)
	DarkKhaki            = Color(0xFFBDB76B)
	DarkMagenta          = Color(0xFF8B008B)
	DarkOliveGreen       = Color(0xFF556B2F)
	DarkOrange           = Color(0xFFFF8C00)
	DarkOrchid           = Color(0xFF9932CC)
	DarkRed              = Color(0xFF8B0000)
	DarkSalmon           = Color(0xFFE9967A)
	DarkSeaGreen         = Color(0xFF8FBC8B)
	DarkSlateBlue        = Color(0xFF483D8B)
	DarkSlateGray        = Color(0xFF2F4F4F)
	DarkTurquoise        = Color(0xFF00CED1)
	DarkViolet           = Color(0xFF9400D3)
	DeepPink             = Color(0xFFFF1493)
	DeepSkyBlue          = Color(0xFF00BFFF)
	DimGray              = Color(0xFF696969)
	DodgerBlue           = Color(0xFF1E90FF)
	Firebrick            = Color(0xFFB22222)
	FloralWhite          = Color(0xFFFFFAF0)
	ForestGreen          = Color(0xFF228B22)
	Gainsboro            = Color(0xFFDCDCDC)
	GhostWhite           = Color(0xFFF8F8FF)
	Gold                 = Color(0xFFFFD700)
	Goldenrod            = Color(0xFFDAA520)
	Gray                 = Color(0xFF808080)
	Green                = Color(0xFF008000)
	GreenYellow          = Color(0xFFADFF2F)
	Honeydew             = Color(0xFFF0FFF0)
	HotPink              = Color(0xFFFF69B4)
	IndianRed            = Color(0xFFCD5C5C)
	Indigo               = Color(0xFF4B0082)
	Ivory                = Color(0xFFFFFFF0)
	khaki                = Color(0xFFF0E68C)
	Lavender             = Color(0xFFE6E6FA)
	LavenderBlush        = Color(0xFFFFF0F5)
	LawnGreen            = Color(0xFF7CFC00)
	LemonChiffon         = Color(0xFFFFFACD)
	LightBlue            = Color(0xFFADD8E6)
	LightCoral           = Color(0xFFF08080)
	LightCyan            = Color(0xFFE0FFFF)
	LightGoldenrodYellow = Color(0xFFFAFAD2)
	LightGray            = Color(0xFFD3D3D3)
	LightGreen           = Color(0xFF90EE90)
	LightPink            = Color(0xFFFFB6C1)
	LightSalmon          = Color(0xFFFFA07A)
	LightSeaGreen        = Color(0xFF20B2AA)
	LightSkyBlue         = Color(0xFF87CEFA)
	LightSlateGray       = Color(0xFF778899)
	LightSteelBlue       = Color(0xFFB0C4DE)
	LightYellow          = Color(0xFFFFFFE0)
	Lime                 = Color(0xFF00FF00)
	LimeGreen            = Color(0xFF32CD32)
	Linen                = Color(0xFFFAF0E6)
	Magenta              = Color(0xFFFF00FF)
	Maroon               = Color(0xFF800000)
	MediumAquamarine     = Color(0xFF66CDAA)
	MediumBlue           = Color(0xFF0000CD)
	MediumOrchid         = Color(0xFFBA55D3)
	MediumPurple         = Color(0xFF9370DB)
	MediumSeaGreen       = Color(0xFF3CB371)
	MediumSlateBlue      = Color(0xFF7B68EE)
	MediumSpringGreen    = Color(0xFF00FA9A)
	MediumTurquoise      = Color(0xFF48D1CC)
	MediumVioletRed      = Color(0xFFC71585)
	MidnightBlue         = Color(0xFF191970)
	MintCream            = Color(0xFFF5FFFA)
	MistyRose            = Color(0xFFFFE4E1)
	Moccasin             = Color(0xFFFFE4B5)
	NavajoWhite          = Color(0xFFFFDEAD)
	Navy                 = Color(0xFF000080)
	OldLace              = Color(0xFFFDF5E6)
	Olive                = Color(0xFF808000)
	OliveDrab            = Color(0xFF6B8E23)
	Orange               = Color(0xFFFFA500)
	OrangeRed            = Color(0xFFFF4500)
	Orchid               = Color(0xFFDA70D6)
	PaleGoldenrod        = Color(0xFFEEE8AA)
	PaleGreen            = Color(0xFF98FB98)
	PaleTurquoise        = Color(0xFFAFEEEE)
	PaleVioletRed        = Color(0xFFDB7093)
	PapayaWhip           = Color(0xFFFFEFD5)
	PeachPuff            = Color(0xFFFFDAB9)
	Peru                 = Color(0xFFCD853F)
	Pink                 = Color(0xFFFFC0CB)
	Plum                 = Color(0xFFDDA0DD)
	PowderBlue           = Color(0xFFB0E0E6)
	Purple               = Color(0xFF800080)
	Red                  = Color(0xFFFF0000)
	RosyBrown            = Color(0xFFBC8F8F)
	RoyalBlue            = Color(0xFF4169E1)
	SaddleBrown          = Color(0xFF8B4513)
	Salmon               = Color(0xFFFA8072)
	SandyBrown           = Color(0xFFF4A460)
	SeaGreen             = Color(0xFF2E8B57)
	SeaShell             = Color(0xFFFFF5EE)
	Sienna               = Color(0xFFA0522D)
	Silver               = Color(0xFFC0C0C0)
	SkyBlue              = Color(0xFF87CEEB)
	SlateBlue            = Color(0xFF6A5ACD)
	SlateGray            = Color(0xFF708090)
	Snow                 = Color(0xFFFFFAFA)
	SpringGreen          = Color(0xFF00FF7F)
	SteelBlue            = Color(0xFF4682B4)
	Tan                  = Color(0xFFD2B48C)
	Teal                 = Color(0xFF008080)
	Thistle              = Color(0xFFD8BFD8)
	Tomato               = Color(0xFFFF6347)
	Transparent          = Color(0x00000000)
	Turquoise            = Color(0xFF40E0D0)
	Violet               = Color(0xFFEE82EE)
	Wheat                = Color(0xFFF5DEB3)
	White                = Color(0xFFFFFFFF)
	WhiteSmoke           = Color(0xFFF5F5F5)
	Yellow               = Color(0xFFFFFF00)
	YellowGreen          = Color(0xFF9ACD32)
)

// AllNamed gets a list of all colors with names.
func AllNamed() []Color {
	return utils.SortedKeys(toName, func(x, y Color) int {
		return strings.Compare(toName[x], toName[y])
	})
}

var toName = map[Color]string{
	AliceBlue:            `AliceBlue`,
	AntiqueWhite:         `AntiqueWhite`,
	Aquamarine:           `Aquamarine`,
	Azure:                `Azure`,
	Beige:                `Beige`,
	Bisque:               `Bisque`,
	Black:                `Black`,
	BlanchedAlmond:       `BlanchedAlmond`,
	Blue:                 `Blue`,
	BlueViolet:           `BlueViolet`,
	Brown:                `Brown`,
	BurlyWood:            `BurlyWood`,
	CadetBlue:            `CadetBlue`,
	Chartreuse:           `Chartreuse`,
	Chocolate:            `Chocolate`,
	Coral:                `Coral`,
	CornflowerBlue:       `CornflowerBlue`,
	Cornsilk:             `Cornsilk`,
	Crimson:              `Crimson`,
	Cyan:                 `Cyan`,
	DarkBlue:             `DarkBlue`,
	DarkCyan:             `DarkCyan`,
	DarkGoldenrod:        `DarkGoldenrod`,
	DarkGray:             `DarkGray`,
	DarkGreen:            `DarkGreen`,
	DarkKhaki:            `DarkKhaki`,
	DarkMagenta:          `DarkMagenta`,
	DarkOliveGreen:       `DarkOliveGreen`,
	DarkOrange:           `DarkOrange`,
	DarkOrchid:           `DarkOrchid`,
	DarkRed:              `DarkRed`,
	DarkSalmon:           `DarkSalmon`,
	DarkSeaGreen:         `DarkSeaGreen`,
	DarkSlateBlue:        `DarkSlateBlue`,
	DarkSlateGray:        `DarkSlateGray`,
	DarkTurquoise:        `DarkTurquoise`,
	DarkViolet:           `DarkViolet`,
	DeepPink:             `DeepPink`,
	DeepSkyBlue:          `DeepSkyBlue`,
	DimGray:              `DimGray`,
	DodgerBlue:           `DodgerBlue`,
	Firebrick:            `Firebrick`,
	FloralWhite:          `FloralWhite`,
	ForestGreen:          `ForestGreen`,
	Gainsboro:            `Gainsboro`,
	GhostWhite:           `GhostWhite`,
	Gold:                 `Gold`,
	Goldenrod:            `Goldenrod`,
	Gray:                 `Gray`,
	Green:                `Green`,
	GreenYellow:          `GreenYellow`,
	Honeydew:             `Honeydew`,
	HotPink:              `HotPink`,
	IndianRed:            `IndianRed`,
	Indigo:               `Indigo`,
	Ivory:                `Ivory`,
	khaki:                `khaki`,
	Lavender:             `Lavender`,
	LavenderBlush:        `LavenderBlush`,
	LawnGreen:            `LawnGreen`,
	LemonChiffon:         `LemonChiffon`,
	LightBlue:            `LightBlue`,
	LightCoral:           `LightCoral`,
	LightCyan:            `LightCyan`,
	LightGoldenrodYellow: `LightGoldenrodYellow`,
	LightGray:            `LightGray`,
	LightGreen:           `LightGreen`,
	LightPink:            `LightPink`,
	LightSalmon:          `LightSalmon`,
	LightSeaGreen:        `LightSeaGreen`,
	LightSkyBlue:         `LightSkyBlue`,
	LightSlateGray:       `LightSlateGray`,
	LightSteelBlue:       `LightSteelBlue`,
	LightYellow:          `LightYellow`,
	Lime:                 `Lime`,
	LimeGreen:            `LimeGreen`,
	Linen:                `Linen`,
	Magenta:              `Magenta`,
	Maroon:               `Maroon`,
	MediumAquamarine:     `MediumAquamarine`,
	MediumBlue:           `MediumBlue`,
	MediumOrchid:         `MediumOrchid`,
	MediumPurple:         `MediumPurple`,
	MediumSeaGreen:       `MediumSeaGreen`,
	MediumSlateBlue:      `MediumSlateBlue`,
	MediumSpringGreen:    `MediumSpringGreen`,
	MediumTurquoise:      `MediumTurquoise`,
	MediumVioletRed:      `MediumVioletRed`,
	MidnightBlue:         `MidnightBlue`,
	MintCream:            `MintCream`,
	MistyRose:            `MistyRose`,
	Moccasin:             `Moccasin`,
	NavajoWhite:          `NavajoWhite`,
	Navy:                 `Navy`,
	OldLace:              `OldLace`,
	Olive:                `Olive`,
	OliveDrab:            `OliveDrab`,
	Orange:               `Orange`,
	OrangeRed:            `OrangeRed`,
	Orchid:               `Orchid`,
	PaleGoldenrod:        `PaleGoldenrod`,
	PaleGreen:            `PaleGreen`,
	PaleTurquoise:        `PaleTurquoise`,
	PaleVioletRed:        `PaleVioletRed`,
	PapayaWhip:           `PapayaWhip`,
	PeachPuff:            `PeachPuff`,
	Peru:                 `Peru`,
	Pink:                 `Pink`,
	Plum:                 `Plum`,
	PowderBlue:           `PowderBlue`,
	Purple:               `Purple`,
	Red:                  `Red`,
	RosyBrown:            `RosyBrown`,
	RoyalBlue:            `RoyalBlue`,
	SaddleBrown:          `SaddleBrown`,
	Salmon:               `Salmon`,
	SandyBrown:           `SandyBrown`,
	SeaGreen:             `SeaGreen`,
	SeaShell:             `SeaShell`,
	Sienna:               `Sienna`,
	Silver:               `Silver`,
	SkyBlue:              `SkyBlue`,
	SlateBlue:            `SlateBlue`,
	SlateGray:            `SlateGray`,
	Snow:                 `Snow`,
	SpringGreen:          `SpringGreen`,
	SteelBlue:            `SteelBlue`,
	Tan:                  `Tan`,
	Teal:                 `Teal`,
	Thistle:              `Thistle`,
	Tomato:               `Tomato`,
	Transparent:          `Transparent`,
	Turquoise:            `Turquoise`,
	Violet:               `Violet`,
	Wheat:                `Wheat`,
	White:                `White`,
	WhiteSmoke:           `WhiteSmoke`,
	Yellow:               `Yellow`,
	YellowGreen:          `YellowGreen`,
}

var fromLowerCase = func() map[string]Color {
	byName := make(map[string]Color, len(toName)+2)
	for c, s := range toName {
		byName[strings.ToLower(s)] = c
	}
	byName[`aqua`] = Cyan
	byName[`fuchsia`] = Magenta
	return byName
}()
