package printer

import (
	"cligram/internal/theme/tokens"

	"github.com/gookit/color"
)

var (
	Warn    = color.HEX(tokens.ColorWarn500)
	Error   = color.HEX(tokens.ColorNegative500)
	Success = color.HEX(tokens.ColorPositive500)
)
