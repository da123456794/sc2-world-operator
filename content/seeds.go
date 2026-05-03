package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Seeds stores seed type state.
type Seeds struct {
	// SeedType is the seed type index, in range [0,7].
	SeedType uint8 `mapstructure:"seed_type"`
}

func (s *Seeds) Marshal(io bit.IO) {
	io.Uint8(&s.SeedType, 3)
	io.Pad(15)
}
