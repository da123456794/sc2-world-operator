package chunk

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"

	"github.com/Yeah114/sc2-world-operator/define"
)

// Encode compresses a Chunk into a region-file chunk payload (DEFLATE-compressed RLE data).
//
// The output can be written to a region file after the CHK1 magic header.
func Encode(c *Chunk) ([]byte, error) {
	var uncompressed []byte

	// Encode 256 shaft bytes (outer=x, inner=z).
	// Each byte: high nibble = temperature, low nibble = humidity.
	for x := 0; x < XSize; x++ {
		for z := 0; z < ZSize; z++ {
			shaft := c.Shaft(x, z)
			temp := byte(define.ExtractTemperature(shaft))
			humid := byte(define.ExtractHumidity(shaft))
			uncompressed = append(uncompressed, (temp<<4)|humid)
		}
	}

	// RLE encode cells (traversal: outer=y, mid=z, inner=x).
	// A run is broken when the cell value (with light cleared) changes or
	// count reaches maxRLECount.
	runVal := int32(0)
	runLen := 0

	appendRLE := func(val int32, count int) {
		cellVal := define.ReplaceLight(val, 0)
		if count < 16 {
			word := define.ReplaceLight(cellVal, int32(count-1))
			var buf [4]byte
			binary.LittleEndian.PutUint32(buf[:], uint32(word))
			uncompressed = append(uncompressed, buf[:]...)
		} else {
			word := define.ReplaceLight(cellVal, 15)
			var buf [4]byte
			binary.LittleEndian.PutUint32(buf[:], uint32(word))
			uncompressed = append(uncompressed, buf[:]...)
			uncompressed = append(uncompressed, byte(count-16))
		}
	}

	first := true
	for y := 0; y < YSize; y++ {
		for z := 0; z < ZSize; z++ {
			for x := 0; x < XSize; x++ {
				// Strip light bits; they are repurposed for RLE count.
				cell := define.ReplaceLight(c.Cell(x, y, z), 0)
				if first {
					runVal = cell
					runLen = 1
					first = false
					continue
				}
				if cell == runVal && runLen < maxRLECount {
					runLen++
					if runLen == maxRLECount {
						appendRLE(runVal, runLen)
						runLen = 0
						first = true
					}
				} else {
					if runLen > 0 {
						appendRLE(runVal, runLen)
					}
					runVal = cell
					runLen = 1
				}
			}
		}
	}
	if runLen > 0 {
		appendRLE(runVal, runLen)
	}

	// DEFLATE compress (BestSpeed = level 1).
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("chunk.Encode: flate.NewWriter: %w", err)
	}
	if _, err := w.Write(uncompressed); err != nil {
		return nil, fmt.Errorf("chunk.Encode: flate write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("chunk.Encode: flate close: %w", err)
	}
	return buf.Bytes(), nil
}
