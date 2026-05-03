package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Gravestone stores gravestone variant and rotation.
type Gravestone struct {
	// Variant is the gravestone variant index.
	Variant uint8 `mapstructure:"variant"`
	// Rotation is the gravestone rotation index.
	Rotation uint8 `mapstructure:"rotation"`
}

func (s *Gravestone) Marshal(io bit.IO) {
	io.Uint8(&s.Variant, 3)
	io.Uint8(&s.Rotation, 1)
}
