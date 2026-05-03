package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Soil stores soil hydration and nitrogen.
type Soil struct {
	// Hydration is the hydration level, in range [0,15].
	Hydration uint8 `mapstructure:"hydration"`
	// Nitrogen is the nitrogen level, in range [0,15].
	Nitrogen uint8 `mapstructure:"nitrogen"`
}

func (s *Soil) Marshal(io bit.IO) {
	io.Uint8(&s.Hydration, 4)
	io.Uint8(&s.Nitrogen, 4)
}
