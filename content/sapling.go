package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Sapling stores sapling tree type index.
type Sapling struct {
	// TreeType is the sapling type index.
	TreeType uint8 `mapstructure:"tree_type"`
}

func (s *Sapling) Marshal(io bit.IO) {
	io.Uint8(&s.TreeType, 4)
}
