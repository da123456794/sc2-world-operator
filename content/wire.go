package content

import "github.com/Yeah114/sc2-world-operator/bit"

// Wire stores each face connectivity bit.
type Wire struct {
	// ConnectionSouth specifies if a wire connection exists on the south face.
	ConnectionSouth bool `mapstructure:"connection_south"`
	// ConnectionWest specifies if a wire connection exists on the west face.
	ConnectionWest bool `mapstructure:"connection_west"`
	// ConnectionNorth specifies if a wire connection exists on the north face.
	ConnectionNorth bool `mapstructure:"connection_north"`
	// ConnectionEast specifies if a wire connection exists on the east face.
	ConnectionEast bool `mapstructure:"connection_east"`
	// ConnectionTop specifies if a wire connection exists on the top face.
	ConnectionTop bool `mapstructure:"connection_top"`
	// ConnectionBottom specifies if a wire connection exists on the bottom face.
	ConnectionBottom bool `mapstructure:"connection_bottom"`
}

func (s *Wire) Marshal(io bit.IO) {
	io.Bool(&s.ConnectionSouth)
	io.Bool(&s.ConnectionWest)
	io.Bool(&s.ConnectionNorth)
	io.Bool(&s.ConnectionEast)
	io.Bool(&s.ConnectionTop)
	io.Bool(&s.ConnectionBottom)
}
