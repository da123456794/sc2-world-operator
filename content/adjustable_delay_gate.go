package content

import "github.com/Yeah114/sc2-world-operator/bit"

// AdjustableDelayGate stores adjustable delay gate state.
type AdjustableDelayGate struct {
	// Rotation is the in-plane rotation, in range [0,3].
	Rotation uint8 `mapstructure:"rotation"`
	// Face is the mounting face index.
	Face uint8 `mapstructure:"face"`
	// Delay is the gate delay value, using 10 bits.
	Delay int32 `mapstructure:"delay"`
}

func (s *AdjustableDelayGate) Marshal(io bit.IO) {
	io.Uint8(&s.Rotation, 2)
	io.Uint8(&s.Face, 3)
	io.Int32(&s.Delay, 10)
}
