package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Dispenser stores dispenser block state.
type Dispenser struct {
	// Facing is the dispenser output face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Facing uint8 `mapstructure:"facing"`
	// AcceptsDrops specifies if the dispenser accepts dropped items.
	AcceptsDrops bool `mapstructure:"accepts_drops"`
}

func (s *Dispenser) Marshal(io bit.IO) {
	io.Uint8(&s.Facing, 3)
	io.Pad(1)
	io.Bool(&s.AcceptsDrops)
}
