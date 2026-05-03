package region

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/Yeah114/sc2-world-operator/chunk"
	"github.com/Yeah114/sc2-world-operator/define"
)

const (
	regionMagic  uint32 = 0x314E4752                  // "RGN1" as little-endian uint32
	chunkMagic   uint32 = 0x314B4843                  // "CHK1" as little-endian uint32
	dirEntries          = 256                         // 16×16 chunk slots per region file
	dirEntrySize        = 8                           // int32 offset + int32 size (bytes)
	dataOffset          = 4 + dirEntries*dirEntrySize // = 2052
	extraSpace          = 1024                        // padding appended after each chunk block
)

// dirEntry..
type dirEntry struct {
	Offset int32 // absolute byte offset in the file; 0 = not present
	Size   int32 // compressed payload size (NOT including the 4-byte CHK1 magic)
}

// Region wraps a single .dat region file and provides LoadChunk / SaveChunk.
// It is safe for concurrent use.
type Region struct {
	mu   sync.Mutex
	f    *os.File
	path string
}

// Open opens the region file at path in read-write mode, creating it if necessary.
func Open(path string) (*Region, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("region.Open %q: %w", path, err)
	}
	r := &Region{f: f, path: path}

	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	if stat.Size() == 0 {
		if err := r.initHeader(); err != nil {
			_ = f.Close()
			return nil, err
		}
	} else {
		var buf [4]byte
		if _, err := f.ReadAt(buf[:], 0); err != nil {
			_ = f.Close()
			return nil, fmt.Errorf("region.Open %q: read magic: %w", path, err)
		}
		if binary.LittleEndian.Uint32(buf[:]) != regionMagic {
			_ = f.Close()
			return nil, fmt.Errorf("region.Open %q: invalid magic", path)
		}
	}
	return r, nil
}

// Close closes the underlying file.
func (r *Region) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.f.Close()
}

// chunkSlot returns the 0-based slot index for a local chunk coordinate (0-15, 0-15).
func chunkSlot(lx, lz int) int { return lx + lz*16 }

// LocalCoords returns the local chunk position within the region (0-15, 0-15)
// for global chunk coordinates.
func LocalCoords(pos define.ChunkPos) (lx, lz int) {
	lx = int(pos.X()) & 0xF
	lz = int(pos.Z()) & 0xF
	return
}

// RegionCoords returns the region coordinate that contains the given chunk.
func RegionCoords(pos define.ChunkPos) (rx, rz int32) {
	rx = pos.X() >> 4
	rz = pos.Z() >> 4
	return
}

// LoadChunk reads and decodes the chunk at local position (lx, lz) within this region.
// Returns (nil, nil) if the chunk does not exist yet.
func (r *Region) LoadChunk(lx, lz int) (*chunk.Chunk, error) {
	compressed, err := r.LoadChunkPayloadOnly(lx, lz)
	if err != nil || compressed == nil {
		return nil, err
	}
	return chunk.Decode(compressed)
}

// LoadChunkPayloadOnly reads the raw compressed chunk payload at local position
// (lx, lz) without decoding.
// Returns (nil, nil) if the chunk does not exist yet.
func (r *Region) LoadChunkPayloadOnly(lx, lz int) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, err := r.readDirEntry(lx, lz)
	if err != nil {
		return nil, err
	}
	if entry.Offset == 0 {
		return nil, nil
	}

	return r.readChunkData(entry)
}

// SaveChunk encodes and writes chunk c to local position (lx, lz) within this region.
func (r *Region) SaveChunk(lx, lz int, c *chunk.Chunk) error {
	compressed, err := chunk.Encode(c)
	if err != nil {
		return err
	}
	return r.SaveChunkPayloadOnly(lx, lz, compressed)
}

// SaveChunkPayloadOnly writes a raw compressed chunk payload to local position
// (lx, lz) without encoding.
func (r *Region) SaveChunkPayloadOnly(lx, lz int, payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if payload == nil {
		payload = []byte{}
	}
	return r.writeChunkData(lx, lz, payload)
}

// --- internal helpers ---

// initHeader writes the RGN1 magic and empty 256-entry directory to a new file.
func (r *Region) initHeader() error {
	var buf [4 + dirEntries*dirEntrySize]byte
	binary.LittleEndian.PutUint32(buf[:4], regionMagic)
	// directory entries are all zeros (Offset=0 means not present)
	_, err := r.f.WriteAt(buf[:], 0)
	return err
}

// readDirEntry reads the directory entry for local chunk (lx, lz).
func (r *Region) readDirEntry(lx, lz int) (dirEntry, error) {
	slot := chunkSlot(lx, lz)
	pos := int64(4 + slot*dirEntrySize)
	var buf [dirEntrySize]byte
	if _, err := r.f.ReadAt(buf[:], pos); err != nil {
		return dirEntry{}, fmt.Errorf("region: read dir entry: %w", err)
	}
	return dirEntry{
		Offset: int32(binary.LittleEndian.Uint32(buf[0:])),
		Size:   int32(binary.LittleEndian.Uint32(buf[4:])),
	}, nil
}

// writeDirEntry writes a directory entry for local chunk (lx, lz).
func (r *Region) writeDirEntry(lx, lz int, e dirEntry) error {
	slot := chunkSlot(lx, lz)
	pos := int64(4 + slot*dirEntrySize)
	var buf [dirEntrySize]byte
	binary.LittleEndian.PutUint32(buf[0:], uint32(e.Offset))
	binary.LittleEndian.PutUint32(buf[4:], uint32(e.Size))
	_, err := r.f.WriteAt(buf[:], pos)
	return err
}

// readAllDirEntries reads all 256 directory entries.
func (r *Region) readAllDirEntries() ([dirEntries]dirEntry, error) {
	var buf [dirEntries * dirEntrySize]byte
	if _, err := r.f.ReadAt(buf[:], 4); err != nil && err != io.EOF {
		return [dirEntries]dirEntry{}, fmt.Errorf("region: read dir: %w", err)
	}
	var entries [dirEntries]dirEntry
	for i := range entries {
		entries[i].Offset = int32(binary.LittleEndian.Uint32(buf[i*8:]))
		entries[i].Size = int32(binary.LittleEndian.Uint32(buf[i*8+4:]))
	}
	return entries, nil
}

// readChunkData reads the compressed payload described by entry.
func (r *Region) readChunkData(e dirEntry) ([]byte, error) {
	// Layout at e.Offset: [CHK1 magic (4 bytes)] [compressed data (e.Size bytes)]
	var magic [4]byte
	if _, err := r.f.ReadAt(magic[:], int64(e.Offset)); err != nil {
		return nil, fmt.Errorf("region: read chunk magic: %w", err)
	}
	if binary.LittleEndian.Uint32(magic[:]) != chunkMagic {
		return nil, fmt.Errorf("region: invalid chunk magic at offset %d", e.Offset)
	}
	buf := make([]byte, e.Size)
	if _, err := r.f.ReadAt(buf, int64(e.Offset)+4); err != nil {
		return nil, fmt.Errorf("region: read chunk data: %w", err)
	}
	return buf, nil
}

// writeChunkData writes compressed payload data for local chunk (lx, lz).
func (r *Region) writeChunkData(lx, lz int, data []byte) error {
	entries, err := r.readAllDirEntries()
	if err != nil {
		return err
	}
	slot := chunkSlot(lx, lz)
	newSize := int32(len(data))

	existing := entries[slot]
	if existing.Offset > 0 {
		// Check available space: distance to next entry minus 4 (CHK1 magic).
		available := r.availableSpace(entries, slot)
		if newSize <= available {
			// Write in place.
			return r.writeBlock(existing.Offset, data, lx, lz, newSize)
		}
		// Not enough space: full reorganise.
		return r.reorganise(entries, slot, data)
	}

	// New entry: append at end of file.
	fileEnd, err := r.f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("region: seek end: %w", err)
	}
	newOffset := int32(fileEnd)
	// Extend file to hold CHK1 + data + extraSpace.
	totalSpace := int64(4 + newSize + extraSpace)
	if err := r.f.Truncate(fileEnd + totalSpace); err != nil {
		return fmt.Errorf("region: truncate: %w", err)
	}
	return r.writeBlock(newOffset, data, lx, lz, newSize)
}

// writeBlock writes CHK1 + data at offset and updates the directory entry.
func (r *Region) writeBlock(offset int32, data []byte, lx, lz int, size int32) error {
	var magic [4]byte
	binary.LittleEndian.PutUint32(magic[:], chunkMagic)
	if _, err := r.f.WriteAt(magic[:], int64(offset)); err != nil {
		return fmt.Errorf("region: write chunk magic: %w", err)
	}
	if _, err := r.f.WriteAt(data, int64(offset)+4); err != nil {
		return fmt.Errorf("region: write chunk data: %w", err)
	}
	return r.writeDirEntry(lx, lz, dirEntry{Offset: offset, Size: size})
}

// availableSpace returns the number of bytes available for in-place data
// replacement at entries[slot], measured as distance to the next entry minus 4.
func (r *Region) availableSpace(entries [dirEntries]dirEntry, slot int) int32 {
	cur := entries[slot].Offset
	minNext := int32(0)
	for i, e := range entries {
		if i == slot || e.Offset <= cur {
			continue
		}
		diff := e.Offset - cur
		if minNext == 0 || diff < minNext {
			minNext = diff
		}
	}
	if minNext == 0 {
		// Last entry in file — no upper bound.
		return 0x7FFFFFFF
	}
	return minNext - 4
}

// reorganise rewrites the entire region file with compacted layout,
// updating entries[slot] to the new data.
func (r *Region) reorganise(entries [dirEntries]dirEntry, slot int, newData []byte) error {
	// Read all existing payloads.
	payloads := make([][]byte, dirEntries)
	for i, e := range entries {
		if e.Offset == 0 {
			continue
		}
		data, err := r.readChunkData(e)
		if err != nil {
			return fmt.Errorf("region: reorganise read slot %d: %w", i, err)
		}
		payloads[i] = data
	}
	payloads[slot] = newData

	// Rebuild the file in a temporary buffer.
	// Layout: [magic 4] [dir 2048] [blocks...]
	var header [4 + dirEntries*dirEntrySize]byte
	binary.LittleEndian.PutUint32(header[:4], regionMagic)

	cursor := int32(dataOffset)
	type block struct {
		offset int32
		data   []byte
	}
	var blocks []block
	newEntries := [dirEntries]dirEntry{}

	for i, payload := range payloads {
		if payload == nil {
			continue
		}
		size := int32(len(payload))
		newEntries[i] = dirEntry{Offset: cursor, Size: size}
		blocks = append(blocks, block{offset: cursor, data: payload})
		cursor += 4 + size + extraSpace
	}

	// Write directory into header.
	for i, e := range newEntries {
		binary.LittleEndian.PutUint32(header[4+i*8:], uint32(e.Offset))
		binary.LittleEndian.PutUint32(header[4+i*8+4:], uint32(e.Size))
	}

	// Truncate and rewrite.
	if err := r.f.Truncate(0); err != nil {
		return fmt.Errorf("region: reorganise truncate: %w", err)
	}
	if _, err := r.f.WriteAt(header[:], 0); err != nil {
		return fmt.Errorf("region: reorganise write header: %w", err)
	}
	for _, b := range blocks {
		var magic [4]byte
		binary.LittleEndian.PutUint32(magic[:], chunkMagic)
		if _, err := r.f.WriteAt(magic[:], int64(b.offset)); err != nil {
			return fmt.Errorf("region: reorganise write magic: %w", err)
		}
		if _, err := r.f.WriteAt(b.data, int64(b.offset)+4); err != nil {
			return fmt.Errorf("region: reorganise write data: %w", err)
		}
	}
	// Extend file to include extraSpace after last block.
	if err := r.f.Truncate(int64(cursor)); err != nil {
		return fmt.Errorf("region: reorganise final truncate: %w", err)
	}
	return nil
}
