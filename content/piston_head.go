package content

import "github.com/Yeah114/sc2-world-operator/bit"

// PistonHead stores piston head block state.
type PistonHead struct {
	// Face is the piston head facing index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// IsShaft specifies if this head part is a shaft segment.
	IsShaft bool `mapstructure:"is_shaft"`
}

func (s *PistonHead) Marshal(io bit.IO) {
	io.Bool(&s.IsShaft)
	io.Pad(2)
	io.Uint8(&s.Face, 3)
}
