package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Magma stores magma fluid state.
type Magma struct {
	// Level is the fluid level value (low 4 bits).
	Level uint8 `mapstructure:"level"`
	// Top specifies if this magma cell is a top surface fluid cell.
	Top bool `mapstructure:"top"`
}

func (s *Magma) Marshal(io bit.IO) {
	io.Uint8(&s.Level, 4)
	io.Bool(&s.Top)
	io.Pad(13)
}
