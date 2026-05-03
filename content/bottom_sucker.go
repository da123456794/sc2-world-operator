package content

import "github.com/Yeah114/sc2-world-operator/bit"

// BottomSucker stores shared sea urchin / starfish state.
type BottomSucker struct {
	// Level is the embedded water fluid level value.
	Level uint8 `mapstructure:"level"`
	// Top specifies if the embedded water state is top-surface fluid.
	Top bool `mapstructure:"top"`
	// Face is the attached face index. 0=south, 1=west, 2=north, 3=east, 4=up, 5=down.
	Face uint8 `mapstructure:"face"`
	// Subvariant is the subvariant index, in range [0,3].
	Subvariant uint8 `mapstructure:"subvariant"`
}

func (s *BottomSucker) Marshal(io bit.IO) {
	io.Uint8(&s.Level, 4)
	io.Bool(&s.Top)
	io.Pad(3)
	io.Uint8(&s.Face, 3)
	io.Uint8(&s.Subvariant, 2)
	io.Pad(5)
}
