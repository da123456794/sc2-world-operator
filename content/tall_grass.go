package content

import "github.com/Yeah114/sc2-world-operator/bit"

// TallGrass stores tall grass variant.
type TallGrass struct {
	// Small specifies if the tall grass uses the small variant.
	Small bool `mapstructure:"small"`
}

func (s *TallGrass) Marshal(io bit.IO) {
	io.Pad(3)
	io.Bool(&s.Small)
}
