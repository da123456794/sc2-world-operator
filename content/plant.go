package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Plant stores growth stage and wild flag.
type Plant struct {
	// Wild specifies if the plant is naturally spawned.
	Wild bool `mapstructure:"wild"`
	// Size is the plant growth stage, in range [0,15].
	Size uint8 `mapstructure:"size"`
}

func (s *Plant) Marshal(io bit.IO) {
	io.Bool(&s.Wild)
	io.Uint8(&s.Size, 4)
}
