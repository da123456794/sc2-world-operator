package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Cairn stores cairn level data.
type Cairn struct {
	// Level is the cairn level value used by SC2 drop logic.
	Level int32 `mapstructure:"level"`
}

func (s *Cairn) Marshal(io bit.IO) {
	io.Int32(&s.Level, 18)
}
