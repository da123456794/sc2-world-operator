package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Thermometer stores thermometer facing state.
type Thermometer struct {
	// Face is the output connector face index. 0=south, 1=west, 2=north, 3=east.
	Face uint8 `mapstructure:"face"`
}

func (s *Thermometer) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 2)
	io.Pad(16)
}
