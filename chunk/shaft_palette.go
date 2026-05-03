package chunk

// shaftPalette stores all 16×16 shaft values (temperature/humidity packed int32)
// using one chunk-wide palette.
//
// Storage modes:
//   - len(palette)==0: implicit all-zero shafts.
//   - len(palette)==1: single value repeated for all shafts.
//   - len(palette)>=2: indexed palette via index8 or index16.
type shaftPalette struct {
	palette []int32
	index8  []uint8
	index16 []uint16
}

// At returns the shaft value at flat shaft index idx.
func (p *shaftPalette) At(idx int) int32 {
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

// Set writes shaft value v at flat shaft index idx.
func (p *shaftPalette) Set(idx int, v int32) {
	if len(p.palette) == 0 {
		if v == 0 {
			return
		}
		p.palette = []int32{0, v}
		p.index8 = make([]uint8, ShaftCount)
		p.index8[idx] = 1
		return
	}

	if len(p.palette) == 1 {
		if p.palette[0] == v {
			return
		}
		p.palette = append(p.palette, v)
		p.index8 = make([]uint8, ShaftCount)
		p.index8[idx] = 1
		return
	}

	if paletteIdx := p.findPalette(v); paletteIdx >= 0 {
		if p.index16 != nil {
			p.index16[idx] = uint16(paletteIdx)
		} else {
			p.index8[idx] = uint8(paletteIdx)
		}
		return
	}

	p.palette = append(p.palette, v)
	newIdx := len(p.palette) - 1

	if p.index16 != nil {
		p.index16[idx] = uint16(newIdx)
		return
	}

	if newIdx <= 255 {
		p.index8[idx] = uint8(newIdx)
		return
	}

	promoted := make([]uint16, ShaftCount)
	for i := range p.index8 {
		promoted[i] = uint16(p.index8[i])
	}
	promoted[idx] = uint16(newIdx)
	p.index16 = promoted
	p.index8 = nil
}

func (p *shaftPalette) findPalette(v int32) int {
	for i, cur := range p.palette {
		if cur == v {
			return i
		}
	}
	return -1
}
