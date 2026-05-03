package content

const (
	FacingSouth  uint8 = iota // +Z
	FacingWest                // -X
	FacingNorth               // -Z
	FacingEast                // +X
	FacingTop                 // +Y
	FacingBottom              // -Y
)
