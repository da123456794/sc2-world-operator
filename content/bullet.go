package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Bullet stores bullet type state.
type Bullet struct {
	// BulletType is the bullet type index (low 4 bits).
	BulletType uint8 `mapstructure:"bullet_type"`
}

func (s *Bullet) Marshal(io bit.IO) {
	io.Uint8(&s.BulletType, 4)
	io.Pad(14)
}
