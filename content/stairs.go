package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Stairs stores shared state for stair blocks.
type Stairs struct {
	// Facing is the horizontal facing direction. Only 2 bits are used.
	Facing uint8 `mapstructure:"facing"`
	// UpsideDown specifies if the stair is upside down.
	UpsideDown bool `mapstructure:"upside_down"`
}

func (s *Stairs) Marshal(io bit.IO) {
	io.Uint8(&s.Facing, 2)
	io.Bool(&s.UpsideDown)
}
