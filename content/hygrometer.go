package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Hygrometer stores hygrometer facing state.
type Hygrometer struct {
	// Face is the output connector face index. 0=south, 1=west, 2=north, 3=east.
	Face uint8 `mapstructure:"face"`
}

func (s *Hygrometer) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 2)
	io.Pad(16)
}
