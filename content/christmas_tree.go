package content

import "github.com/Yeah114/sc2-world-operator/bit"

// ChristmasTree stores whether the tree is lit.
type ChristmasTree struct {
	// Lit specifies if the christmas tree lights are on.
	Lit bool `mapstructure:"lit"`
}

func (s *ChristmasTree) Marshal(io bit.IO) {
	io.Bool(&s.Lit)
}
