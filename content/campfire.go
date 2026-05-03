package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Campfire stores campfire burn level state.
type Campfire struct {
	// Level is the campfire level value, in range [0,15].
	Level uint8 `mapstructure:"level"`
}

func (s *Campfire) Marshal(io bit.IO) {
	io.Uint8(&s.Level, 4)
	io.Pad(14)
}
