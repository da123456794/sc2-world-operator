package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Fireworks stores fireworks state.
type Fireworks struct {
	// Color is the fireworks color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
	// Altitude is the fireworks altitude level, in range [0,15].
	Altitude uint8 `mapstructure:"altitude"`
	// Flickering specifies if the fireworks has flickering effect.
	Flickering bool `mapstructure:"flickering"`
}

func (s *Fireworks) Marshal(io bit.IO) {
	io.Uint8(&s.Color, 4)
	io.Uint8(&s.Altitude, 4)
	io.Bool(&s.Flickering)
}
