package ansi

import "strings"

const (
	Black         = 30
	Red           = 31
	Green         = 32
	Yellow        = 33
	Blue          = 34
	Magenta       = 35
	Cyan          = 36
	White         = 37
	BrightBlack   = 90
	BrightRed     = 91
	BrightGreen   = 92
	BrightYellow  = 93
	BrightBlue    = 94
	BrightMagenta = 95
	BrightCyan    = 96
	BrightWhite   = 97
)

var colorCodes = map[string]int{
	"black":          Black,
	"red":            Red,
	"green":          Green,
	"yellow":         Yellow,
	"blue":           Blue,
	"magenta":        Magenta,
	"cyan":           Cyan,
	"white":          White,
	"bright-black":   BrightBlack,
	"gray":           BrightBlack,
	"grey":           BrightBlack,
	"bright-red":     BrightRed,
	"bright-green":   BrightGreen,
	"bright-yellow":  BrightYellow,
	"bright-blue":    BrightBlue,
	"bright-magenta": BrightMagenta,
	"bright-cyan":    BrightCyan,
	"bright-white":   BrightWhite,
}

func ColorCode(color string) (int, bool) {
	if color == "" {
		return 0, true
	}
	code, ok := colorCodes[strings.ToLower(strings.TrimSpace(color))]
	return code, ok
}
