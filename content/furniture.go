package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Furniture stores furniture layout bits.
type Furniture struct {
	// Rotation is the furniture rotation index, in range [0,3].
	Rotation uint8 `mapstructure:"rotation"`
	// Design is the furniture design index.
	Design int32 `mapstructure:"design"`
	// LightEmitter specifies if the furniture emits light.
	LightEmitter bool `mapstructure:"light_emitter"`
	// Shadow is the shadow strength index, in range [0,3].
	Shadow uint8 `mapstructure:"shadow"`
}

func (s *Furniture) Marshal(io bit.IO) {
	design := s.Design & 0x7FF
	designLow := design & 0x3FF
	designHigh := (design >> 10) & 1

	io.Uint8(&s.Rotation, 2)
	io.Int32(&designLow, 10)
	io.Uint8(&s.Shadow, 2)
	io.Bool(&s.LightEmitter)
	io.Int32(&designHigh, 1)
	s.Design = (designHigh << 10) | designLow
}
