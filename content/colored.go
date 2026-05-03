package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Colored stores one color index.
type Colored struct {
	// Color is the color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *Colored) Marshal(io bit.IO) {
	io.Uint8(&s.Color, 4)
}
