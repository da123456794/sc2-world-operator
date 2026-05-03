package define

// ChunkPos is a 2D chunk coordinate.
// Each chunk covers a 16×16 column of blocks in the X-Z plane.
type ChunkPos [2]int32

// X returns the X component.
func (p ChunkPos) X() int32 { return p[0] }

// Z returns the Z component.
func (p ChunkPos) Z() int32 { return p[1] }

// BlockPos is an absolute 3D block coordinate.
type BlockPos [3]int32

// X returns the X component.
func (p BlockPos) X() int32 { return p[0] }

// Y returns the Y component.
func (p BlockPos) Y() int32 { return p[1] }

// Z returns the Z component.
func (p BlockPos) Z() int32 { return p[2] }

// ChunkPosFromBlock derives the chunk coordinate that contains the given block.
func ChunkPosFromBlock(b BlockPos) ChunkPos {
	return ChunkPos{b[0] >> 4, b[2] >> 4}
}

// LocalX returns the block's X offset within its chunk (0-15).
func LocalX(b BlockPos) int { return int(b[0] & 0xF) }

// LocalZ returns the block's Z offset within its chunk (0-15).
func LocalZ(b BlockPos) int { return int(b[2] & 0xF) }
