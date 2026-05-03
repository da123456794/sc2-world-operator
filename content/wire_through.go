package content

import "github.com/Yeah114/sc2-world-operator/bit"

// WireThrough stores through-wire routing face.
type WireThrough struct {
	// WiredFace is the routed face direction. Valid decoded values are 0, 1, 4.
	WiredFace uint8 `mapstructure:"wired_face"`
}

func (s *WireThrough) Marshal(io bit.IO) {
	encoded := uint8(2)
	switch s.WiredFace {
	case 0, 2:
		encoded = 0
	case 1, 3:
		encoded = 1
	default:
		encoded = 2
	}

	io.Uint8(&encoded, 2)
	io.Pad(16)

	switch encoded & 0x3 {
	case 0:
		s.WiredFace = 0
	case 1:
		s.WiredFace = 1
	default:
		s.WiredFace = 4
	}
}
