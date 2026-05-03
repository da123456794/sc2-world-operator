package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Pumpkin stores pumpkin growth and damage state.
type Pumpkin struct {
	// Size is the visible pumpkin size, in range [0,7].
	Size uint8 `mapstructure:"size"`
	// Dead specifies if the pumpkin is dead.
	Dead bool `mapstructure:"dead"`
	// Damage is the 1-bit pumpkin damage state.
	Damage uint8 `mapstructure:"damage"`
}

func (s *Pumpkin) Marshal(io bit.IO) {
	encodedSize := uint8(7 - (s.Size & 7))

	io.Uint8(&encodedSize, 3)
	io.Bool(&s.Dead)
	io.Uint8(&s.Damage, 1)
	io.Pad(13)

	s.Size = 7 - (encodedSize & 7)
}
