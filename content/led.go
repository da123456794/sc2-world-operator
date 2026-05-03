package content

import "github.com/Yeah114/sc2-world-operator/bit"

// LED stores shared state for LED-like blocks.
type LED struct {
	// Face is the mounting face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// Color is the LED color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *LED) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 3)
	io.Uint8(&s.Color, 4)
}
