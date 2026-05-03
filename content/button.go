package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Button stores button state.
type Button struct {
	// Face is the mounting face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// Voltage is the output voltage level, in range [0,15].
	Voltage uint8 `mapstructure:"voltage"`
}

func (s *Button) Marshal(io bit.IO) {
	voltage := uint8(15 - int32(s.Voltage&0xF))

	io.Uint8(&s.Face, 3)
	io.Pad(1)
	io.Uint8(&voltage, 4)
	s.Voltage = uint8(15 - int32(voltage&0xF))
}
