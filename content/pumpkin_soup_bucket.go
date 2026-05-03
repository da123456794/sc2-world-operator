package content

import "github.com/Yeah114/sc2-world-operator/bit"

// PumpkinSoupBucket stores pumpkin soup bucket data and damage state.
type PumpkinSoupBucket struct {
	// Data is the low 4-bit bucket data value.
	Data uint8 `mapstructure:"data"`
	// Damage is the 12-bit damage value used by item logic.
	Damage int32 `mapstructure:"damage"`
	// Extra stores the top 2 bits of block extra data.
	Extra uint8 `mapstructure:"extra"`
}

func (s *PumpkinSoupBucket) Marshal(io bit.IO) {
	io.Uint8(&s.Data, 4)
	io.Int32(&s.Damage, 12)
	io.Uint8(&s.Extra, 2)
}
