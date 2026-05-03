package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Piston stores piston block state.
type Piston struct {
	// Extended specifies if the piston is extended.
	Extended bool `mapstructure:"extended"`
	// Face is the push face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// MaxExtension is the configured max extension, in range [0,7].
	MaxExtension uint8 `mapstructure:"max_extension"`
	// PullCount is the configured pull count, in range [0,7].
	PullCount uint8 `mapstructure:"pull_count"`
	// Speed is the configured speed value, in range [0,7].
	Speed uint8 `mapstructure:"speed"`
}

func (s *Piston) Marshal(io bit.IO) {
	io.Bool(&s.Extended)
	io.Pad(2)
	io.Uint8(&s.Face, 3)
	io.Uint8(&s.MaxExtension, 3)
	io.Uint8(&s.PullCount, 3)
	io.Uint8(&s.Speed, 3)
	io.Pad(1)
}
