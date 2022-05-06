package main

import (
	"fmt"
	"image/color"
)

func ParseStringToColor(s string) (c color.RGBA) {
	c.R = 0x00
	c.B = 0x00
	c.G = 0x00
	c.A = 0xFF
	switch len(s) {
	case 7:
		_, _ = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, _ = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)

		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		return
	}
	return
}

func ParseColorToString(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
