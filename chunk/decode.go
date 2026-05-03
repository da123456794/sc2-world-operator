package chunk

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Yeah114/sc2-world-operator/define"
)

// maxRLECount is the maximum count a single RLE entry can represent (271).
const maxRLECount = 271

// Decode decompresses a region-file chunk payload (raw DEFLATE-compressed RLE data)
// into a Chunk. buf is the compressed bytes read after the CHK1 magic.
//
// The binary format:
//  1. DEFLATE-decompress buf to get the raw stream.
//  2. First 256 bytes: shaft bytes, one per (x, z) column (outer=x, inner=z).
//     Each byte: high nibble = temperature, low nibble = humidity.
//  3. Remainder: RLE-encoded cells in traversal order (outer=y, mid=z, inner=x).
//     RLE entry (4 or 5 bytes):
//     - word (int32 LE): bits 0-9 = block ID, bits 14-31 = data, bits 10-13 = count-1 (if < 15)
//     - if bits 10-13 == 15: 5th byte = count - 16 (total count = byte + 16)
func Decode(buf []byte) (*Chunk, error) {
	r := flate.NewReader(bytes.NewReader(buf))
	defer r.Close()
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("chunk.Decode: deflate: %w", err)
	}

	if len(raw) < ShaftCount {
		return nil, fmt.Errorf("chunk.Decode: decompressed data too short (%d bytes)", len(raw))
	}

	c := New()
	pos := 0

	// Decode 256 shaft bytes (outer=x, inner=z).
	for x := 0; x < XSize; x++ {
		for z := 0; z < ZSize; z++ {
			b := raw[pos]
			pos++
			humidity := int32(b & 0xF)
			temperature := int32(b >> 4)
			shaft := define.ReplaceHumidity(0, humidity)
			shaft = define.ReplaceTemperature(shaft, temperature)
			c.SetShaft(x, z, shaft)
		}
	}

	// RLE decode cells (traversal: outer=y, mid=z, inner=x).
	cx, cy, cz := 0, 0, 0
	total := 0

	for pos < len(raw) {
		if pos+4 > len(raw) {
			return nil, fmt.Errorf("chunk.Decode: RLE entry truncated at pos=%d", pos)
		}
		word := int32(binary.LittleEndian.Uint32(raw[pos:]))
		pos += 4

		light := define.ExtractLight(word)
		cellVal := define.ReplaceLight(word, 0)

		var count int32
		if light < 15 {
			count = light + 1
		} else {
			if pos >= len(raw) {
				return nil, fmt.Errorf("chunk.Decode: RLE extended count byte missing")
			}
			count = int32(raw[pos]) + 16
			pos++
		}

		for k := int32(0); k < count; k++ {
			if total >= CellCount {
				return nil, fmt.Errorf("chunk.Decode: more cells than expected (%d)", CellCount)
			}
			c.SetCell(cx, cy, cz, cellVal)
			total++
			cx++
			if cx >= XSize {
				cx = 0
				cz++
				if cz >= ZSize {
					cz = 0
					cy++
				}
			}
		}
	}

	if cx != 0 || cy != YSize || cz != 0 {
		return nil, fmt.Errorf("chunk.Decode: corrupt data: ended at x=%d y=%d z=%d (total=%d)", cx, cy, cz, total)
	}
	return c, nil
}
