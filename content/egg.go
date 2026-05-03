package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Egg stores egg block/item state.
type Egg struct {
	// Cooked specifies if the egg is cooked.
	Cooked bool `mapstructure:"cooked"`
	// Laid specifies if the egg has been laid.
	Laid bool `mapstructure:"laid"`
	// EggType is the egg type index (12 bits).
	EggType int32 `mapstructure:"egg_type"`
	// Damage is the 1-bit egg damage state.
	Damage uint8 `mapstructure:"damage"`
}

func (s *Egg) Marshal(io bit.IO) {
	io.Bool(&s.Cooked)
	io.Bool(&s.Laid)
	io.Pad(2)
	io.Int32(&s.EggType, 12)
	io.Uint8(&s.Damage, 1)
	io.Pad(1)
}
