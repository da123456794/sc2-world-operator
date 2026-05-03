package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Water stores water fluid state.
type Water struct {
	// Level is the fluid level value.
	Level uint8 `mapstructure:"level"`
	// Top specifies if this water cell is a top surface fluid cell.
	Top bool `mapstructure:"top"`
}

func (s *Water) Marshal(io bit.IO) {
	io.Uint8(&s.Level, 4)
	io.Bool(&s.Top)
	io.Pad(13)
}
