package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Wood stores wood cut-face state.
type Wood struct {
	// CutFace is the cut-face direction. Valid decoded values are 0, 1, 4.
	CutFace uint8 `mapstructure:"cut_face"`
}

func (s *Wood) Marshal(io bit.IO) {
	encoded := uint8(0)
	switch s.CutFace {
	case 0, 2:
		encoded = 1
	case 1, 3:
		encoded = 2
	default:
		encoded = 0
	}

	io.Uint8(&encoded, 2)
	io.Pad(16)

	switch encoded & 0x3 {
	case 0:
		s.CutFace = 4
	case 1:
		s.CutFace = 0
	default:
		s.CutFace = 1
	}
}
