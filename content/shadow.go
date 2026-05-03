package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Shadow stores explicit shadow strength data.
type Shadow struct {
	// Strength is shadow strength, clamped to [-128, 128].
	Strength int32 `mapstructure:"strength"`
}

func (s *Shadow) Marshal(io bit.IO) {
	shadowStrength := s.Strength
	if shadowStrength < -128 {
		shadowStrength = -128
	}
	if shadowStrength > 128 {
		shadowStrength = 128
	}
	encoded := shadowStrength + 128

	io.Int32(&encoded, 18)
	s.Strength = encoded - 128
}
