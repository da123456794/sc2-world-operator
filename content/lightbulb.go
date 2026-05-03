package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Lightbulb stores lightbulb state.
type Lightbulb struct {
	// Face is the mounting face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// Intensity is the light intensity level, in range [0,15].
	Intensity uint8 `mapstructure:"intensity"`
	// Color is the color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *Lightbulb) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 3)
	io.Uint8(&s.Intensity, 4)
	io.Uint8(&s.Color, 4)
}
