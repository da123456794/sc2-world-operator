package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Battery stores battery voltage.
type Battery struct {
	// Voltage is the stored voltage level, in range [0,15].
	Voltage uint8 `mapstructure:"voltage"`
}

func (s *Battery) Marshal(io bit.IO) {
	io.Uint8(&s.Voltage, 4)
}
