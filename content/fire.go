package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Fire stores fire face mask data.
type Fire struct {
	// SideMask is the side-face fire mask.
	// Bit 0=south, 1=west, 2=north, 3=east. In SC2, 0 means all side faces are on fire.
	SideMask uint8 `mapstructure:"side_mask"`
}

func (s *Fire) Marshal(io bit.IO) {
	io.Uint8(&s.SideMask, 4)
	io.Pad(14)
}
