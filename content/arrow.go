package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Arrow stores arrow type state.
type Arrow struct {
	// ArrowType is the arrow type index (low 4 bits).
	ArrowType uint8 `mapstructure:"arrow_type"`
}

func (s *Arrow) Marshal(io bit.IO) {
	io.Uint8(&s.ArrowType, 4)
	io.Pad(14)
}
