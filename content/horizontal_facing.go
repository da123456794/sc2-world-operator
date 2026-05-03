package content

import "github.com/Yeah114/sc2-world-operator/bit"

// HorizontalFacing stores a horizontal facing value (0-3).
type HorizontalFacing struct {
	// Facing is the horizontal facing direction. Only 2 bits are used.
	Facing uint8 `mapstructure:"facing"`
}

func (s *HorizontalFacing) Marshal(io bit.IO) {
	io.Uint8(&s.Facing, 2)
}
