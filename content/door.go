package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Door stores shared state for door-like blocks.
type Door struct {
	// Facing is the horizontal facing direction. Only 2 bits are used.
	Facing uint8 `mapstructure:"facing"`
	// Open specifies if the door is opened.
	Open bool `mapstructure:"open"`
	// RightHanded specifies if the hinge is on the right side.
	RightHanded bool `mapstructure:"right_handed"`
}

func (s *Door) Marshal(io bit.IO) {
	leftHanded := !s.RightHanded

	io.Uint8(&s.Facing, 2)
	io.Bool(&s.Open)
	io.Bool(&leftHanded)
	s.RightHanded = !leftHanded
}
