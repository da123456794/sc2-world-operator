package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Musket stores musket loading, bullet and damage state.
type Musket struct {
	// LoadState is musket load state. 0=empty, 1=gunpowder, 2=wad, 3=loaded.
	LoadState uint8 `mapstructure:"load_state"`
	// HammerState specifies if the hammer is cocked.
	HammerState bool `mapstructure:"hammer_state"`
	// BulletLoaded specifies if a bullet is loaded.
	BulletLoaded bool `mapstructure:"bullet_loaded"`
	// BulletType is the loaded bullet type index when bullet_loaded is true.
	BulletType uint8 `mapstructure:"bullet_type"`
	// Damage is the musket damage value, in range [0,255].
	Damage uint8 `mapstructure:"damage"`
}

func (s *Musket) Marshal(io bit.IO) {
	bulletType := uint8(0)
	if s.BulletLoaded {
		bulletType = (s.BulletType & 0xF) + 1
	}

	io.Uint8(&s.LoadState, 2)
	io.Bool(&s.HammerState)
	io.Pad(1)
	io.Uint8(&bulletType, 4)
	io.Uint8(&s.Damage, 8)
	io.Pad(2)

	if bulletType == 0 {
		s.BulletLoaded = false
		s.BulletType = 0
	} else {
		s.BulletLoaded = true
		s.BulletType = bulletType - 1
	}
}
