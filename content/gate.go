package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Gate stores shared gate orientation.
type Gate struct {
	// Rotation is the in-plane rotation, in range [0,3].
	Rotation uint8 `mapstructure:"rotation"`
	// Face is the mounting face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
}

func (s *Gate) Marshal(io bit.IO) {
	io.Uint8(&s.Rotation, 2)
	io.Uint8(&s.Face, 3)
}
