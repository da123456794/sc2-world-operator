package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Fence stores fence variant and paint color state.
type Fence struct {
	// Variant is the fence variant value, in range [0,15].
	Variant uint8 `mapstructure:"variant"`
	// Colored specifies if paint color is enabled.
	Colored bool `mapstructure:"colored"`
	// Color is the paint color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *Fence) Marshal(io bit.IO) {
	color := s.Color & 0xF
	if !s.Colored {
		color = 0
	}

	io.Uint8(&s.Variant, 4)
	io.Bool(&s.Colored)
	io.Uint8(&color, 4)
	io.Pad(9)

	if s.Colored {
		s.Color = color
	} else {
		s.Color = 0
	}
}
