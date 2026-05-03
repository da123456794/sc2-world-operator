package chunk

// cellPalette stores one full 16×256×16 chunk using a small palette.
//
// Storage modes:
//   - len(palette)==0: implicit all-zero chunk.
//   - len(palette)==1: single value repeated for all cells.
//   - len(palette)>=2: indexed palette via index8 or index16.
type cellPalette struct {
	palette []int32
	index8  []uint8
	index16 []uint16
}

// At returns the cell value at local flat index idx.
func (p *cellPalette) At(idx int) int32 {
	if len(p.palette) == 0 {
		return 0
	}
	if len(p.palette) == 1 {
		return p.palette[0]
	}
	if p.index16 != nil {
		return p.palette[p.index16[idx]]
	}
	return p.palette[p.index8[idx]]
}

// Set writes value v at local flat index idx.
func (p *cellPalette) Set(idx int, v int32) {
	if len(p.palette) == 0 {
		if v == 0 {
			return
		}
		p.palette = []int32{0, v}
		p.index8 = make([]uint8, CellCount)
		p.index8[idx] = 1
		return
	}

	if len(p.palette) == 1 {
		if p.palette[0] == v {
			return
		}
		p.palette = append(p.palette, v)
		p.index8 = make([]uint8, CellCount)
		p.index8[idx] = 1
		return
	}

	if p.index16 != nil {
		if paletteIdx := p.findPalette(v); paletteIdx >= 0 {
			p.index16[idx] = uint16(paletteIdx)
			return
		}
		p.palette = append(p.palette, v)
		p.index16[idx] = uint16(len(p.palette) - 1)
		return
	}

	if paletteIdx := p.findPalette(v); paletteIdx >= 0 {
		p.index8[idx] = uint8(paletteIdx)
		return
	}

	p.palette = append(p.palette, v)
	newIdx := len(p.palette) - 1
	if newIdx <= 255 {
		p.index8[idx] = uint8(newIdx)
		return
	}

	// Promote index table to uint16 once palette no longer fits uint8.
	promoted := make([]uint16, CellCount)
	for i := range p.index8 {
		promoted[i] = uint16(p.index8[i])
	}
	promoted[idx] = uint16(newIdx)
	p.index16 = promoted
	p.index8 = nil
}

func (p *cellPalette) findPalette(v int32) int {
	for i, cur := range p.palette {
		if cur == v {
			return i
		}
	}
	return -1
}
