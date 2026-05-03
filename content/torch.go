package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Torch stores torch mounting mode.
// Face: 0=south, 1=west, 2=north, 3=east, 4=floor.
type Torch struct {
	// Face is the mounting face index. 0-3 are walls, 4 is floor.
	Face uint8 `mapstructure:"face"`
}

func (s *Torch) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 3)
}
