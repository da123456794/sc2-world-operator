package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Crossbow stores crossbow draw, loaded arrow and damage state.
type Crossbow struct {
	// Draw is the crossbow draw amount, in range [0,15].
	Draw uint8 `mapstructure:"draw"`
	// ArrowLoaded specifies if an arrow is loaded.
	ArrowLoaded bool `mapstructure:"arrow_loaded"`
	// ArrowType is the loaded arrow type index when arrow_loaded is true.
	ArrowType uint8 `mapstructure:"arrow_type"`
	// Damage is the crossbow damage value, in range [0,255].
	Damage uint8 `mapstructure:"damage"`
}

func (s *Crossbow) Marshal(io bit.IO) {
	arrowType := uint8(0)
	if s.ArrowLoaded {
		arrowType = (s.ArrowType & 0xF) + 1
	}

	io.Uint8(&s.Draw, 4)
	io.Uint8(&arrowType, 4)
	io.Uint8(&s.Damage, 8)
	io.Pad(2)

	if arrowType == 0 {
		s.ArrowLoaded = false
		s.ArrowType = 0
	} else {
		s.ArrowLoaded = true
		s.ArrowType = arrowType - 1
	}
}
