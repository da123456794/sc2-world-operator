package content

import "github.com/Yeah114/sc2-world-operator/bit"

// JackOLantern stores jack-o-lantern rotation state.
type JackOLantern struct {
	// Rotation is the horizontal rotation index, in range [0,3].
	Rotation uint8 `mapstructure:"rotation"`
}

func (s *JackOLantern) Marshal(io bit.IO) {
	io.Uint8(&s.Rotation, 2)
	io.Pad(16)
}
