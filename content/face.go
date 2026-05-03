package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Face stores one full face (0-5) value.
type Face struct {
	// Face is the face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
}

func (s *Face) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 3)
}
