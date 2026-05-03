package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Clothing stores clothing index and color.
type Clothing struct {
	// Index is the clothing item index.
	Index int32 `mapstructure:"index"`
	// Color is the clothing color index, in range [0,15].
	Color uint8 `mapstructure:"color"`
}

func (s *Clothing) Marshal(io bit.IO) {
	io.Int32(&s.Index, 8)
	io.Pad(4)
	io.Uint8(&s.Color, 4)
}
