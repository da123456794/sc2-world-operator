package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Grass stores whether the grass block is snowy.
type Grass struct {
	// Snowy specifies if the grass is in snowy variant.
	Snowy bool `mapstructure:"snowy"`
}

func (s *Grass) Marshal(io bit.IO) {
	raw := int32(0)
	if s.Snowy {
		raw = 1
	}

	io.Int32(&raw, bit.BitsBlockExtra)
	s.Snowy = raw != 0
}
