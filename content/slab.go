package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Slab stores whether a slab is the top half.
type Slab struct {
	// Top specifies if the slab occupies the top half.
	Top bool `mapstructure:"top"`
}

func (s *Slab) Marshal(io bit.IO) {
	io.Bool(&s.Top)
}
