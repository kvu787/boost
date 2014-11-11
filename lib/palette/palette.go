package palette

import (
	sf "bitbucket.org/kvu787/gosfml2"
	"math/rand"
)

var (
	RED         sf.Color = sf.Color{156, 56, 57, 255}
	GREEN       sf.Color = sf.ColorGreen()
	BLUE        sf.Color = sf.Color{61, 87, 122, 255}
	YELLOW      sf.Color = sf.ColorYellow()
	GRAY        sf.Color = sf.Color{59, 60, 54, 255}
	LIGHT_GRAY  sf.Color = sf.Color{160, 160, 160, 255}
	WHITE       sf.Color = sf.ColorWhite()
	DARK_GREEN  sf.Color = sf.Color{25, 112, 8, 255}
	LASER_BLUE  sf.Color = sf.Color{12, 231, 242, 255}
	BLACK       sf.Color = sf.ColorBlack()
	LIGHT_BROWN sf.Color = sf.Color{185, 122, 87, 255}
	DARK_BROWN  sf.Color = sf.Color{94, 47, 0, 255}
	TRANSPARENT sf.Color = sf.Color{0, 0, 0, 0}
)

func RandomColor() sf.Color {
	return sf.Color{
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		255}
}

func SetAlpha(color sf.Color, alpha uint8) sf.Color {
	newColor := color
	newColor.A = alpha
	return newColor
}
