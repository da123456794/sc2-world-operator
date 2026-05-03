package content

import "github.com/Yeah114/sc2-world-operator/bit"

// PaintedCube stores painted cube color state.
type PaintedCube struct {
	// Colored specifies if paint color is enabled.
	Colored bool `mapstructure:"colored"`
	// Color is the paint color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *PaintedCube) Marshal(io bit.IO) {
	colored := s.Colored
	color := s.Color & 0xF
	if !colored {
		color = 0
	}

	io.Bool(&colored)
	io.Uint8(&color, 4)
	io.Pad(13)

	s.Colored = colored
	if colored {
		s.Color = color
	} else {
		s.Color = 0
	}
}
