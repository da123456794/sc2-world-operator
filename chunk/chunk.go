package chunk

import "github.com/Yeah114/sc2-world-operator/define"

const (
	XSize      = 16
	YSize      = 256
	ZSize      = 16
	CellCount  = XSize * YSize * ZSize // 65536
	ShaftCount = XSize * ZSize         // 256
)

// Chunk holds a 16×256×16 block grid and a 16×16 shaft (climate) grid.
type Chunk struct {
	cells  cellPalette
	shafts shaftPalette
}

// New returns an empty Chunk (all cells and shafts are zero = air).
func New() *Chunk { return &Chunk{} }

// cellIdx returns the flat index for (x, y, z).
func cellIdx(x, y, z int) int {
	return y + x*YSize + z*YSize*XSize
}

// shaftIdx returns the flat index for shaft column (x, z).
func shaftIdx(x, z int) int {
	return x + z*XSize
}

// Cell returns the full int32 cell value at (x, y, z).
func (c *Chunk) Cell(x, y, z int) int32 {
	return c.cells.At(cellIdx(x, y, z))
}

// SetCell sets the full int32 cell value at (x, y, z).
func (c *Chunk) SetCell(x, y, z int, v int32) {
	c.cells.Set(cellIdx(x, y, z), v)
}

// BlockID returns the block ID (Contents) at (x, y, z).
func (c *Chunk) BlockID(x, y, z int) int32 {
	return define.ExtractContents(c.Cell(x, y, z))
}

// SetBlockID sets the block ID at (x, y, z), preserving light and data.
func (c *Chunk) SetBlockID(x, y, z int, id int32) {
	v := c.Cell(x, y, z)
	c.SetCell(x, y, z, define.ReplaceContents(v, id))
}

// Light returns the light value at (x, y, z).
func (c *Chunk) Light(x, y, z int) int32 {
	return define.ExtractLight(c.Cell(x, y, z))
}

// SetLight sets the light value at (x, y, z), preserving block ID and data.
func (c *Chunk) SetLight(x, y, z int, light int32) {
	v := c.Cell(x, y, z)
	c.SetCell(x, y, z, define.ReplaceLight(v, light))
}

// Data returns the block data at (x, y, z).
func (c *Chunk) Data(x, y, z int) int32 {
	return define.ExtractData(c.Cell(x, y, z))
}

// SetData sets the block data at (x, y, z), preserving block ID and light.
func (c *Chunk) SetData(x, y, z int, data int32) {
	v := c.Cell(x, y, z)
	c.SetCell(x, y, z, define.ReplaceData(v, data))
}

// Shaft returns the raw shaft value for column (x, z).
func (c *Chunk) Shaft(x, z int) int32 {
	return c.shafts.At(shaftIdx(x, z))
}

// SetShaft sets the raw shaft value for column (x, z).
func (c *Chunk) SetShaft(x, z int, v int32) {
	c.shafts.Set(shaftIdx(x, z), v)
}

// Temperature returns the temperature for column (x, z).
func (c *Chunk) Temperature(x, z int) int32 {
	return define.ExtractTemperature(c.Shaft(x, z))
}

// Humidity returns the humidity for column (x, z).
func (c *Chunk) Humidity(x, z int) int32 {
	return define.ExtractHumidity(c.Shaft(x, z))
}
