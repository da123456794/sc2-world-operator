package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Leaves stores leaf season and shake state.
type Leaves struct {
	// Season is the leaves season index.
	Season uint8 `mapstructure:"season"`
	// Shaken specifies if the leaves are in shaken state.
	Shaken bool `mapstructure:"shaken"`
}

func (s *Leaves) Marshal(io bit.IO) {
	io.Uint8(&s.Season, 2)
	io.Pad(3)
	io.Bool(&s.Shaken)
}
