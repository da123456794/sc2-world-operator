package content

import "github.com/Yeah114/sc2-world-operator/bit"

// WallFace stores one horizontal face (0-3) value.
type WallFace struct {
	// Face is the horizontal face index. 0=south, 1=west, 2=north, 3=east.
	Face uint8 `mapstructure:"face"`
}

func (s *WallFace) Marshal(io bit.IO) {
	io.Uint8(&s.Face, 2)
}
