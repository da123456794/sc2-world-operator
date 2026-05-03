package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Bedrock stores bedrock variant data.
type Bedrock struct {
	// Variant is the bedrock variant value.
	Variant int32 `mapstructure:"variant"`
}

func (s *Bedrock) Marshal(io bit.IO) {
	io.Int32(&s.Variant, 18)
}
