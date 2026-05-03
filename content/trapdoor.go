package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Trapdoor represents the state of a trapdoor block.
type Trapdoor struct {
	// Facing is the direction the trapdoor is facing. Only 2 bits are used.
	Facing uint8 `mapstructure:"facing"`
	// Open specifies if the trapdoor is opened.
	Open bool `mapstructure:"open"`
	// UpsideDown specifies if the trapdoor is placed upside down.
	UpsideDown bool `mapstructure:"upside_down"`
}

func (t *Trapdoor) Marshal(io bit.IO) {
	io.Uint8(&t.Facing, 2)
	io.Bool(&t.Open)
	io.Bool(&t.UpsideDown)
}
