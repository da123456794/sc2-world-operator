package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Magnet stores magnet orientation state.
type Magnet struct {
	// Axis is magnet axis orientation. 0=x-axis, 1=z-axis.
	Axis bool `mapstructure:"axis"`
}

func (s *Magnet) Marshal(io bit.IO) {
	io.Bool(&s.Axis)
	io.Pad(17)
}
