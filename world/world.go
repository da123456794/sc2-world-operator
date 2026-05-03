package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Yeah114/sc2-world-operator/chunk"
	"github.com/Yeah114/sc2-world-operator/define"
	"github.com/Yeah114/sc2-world-operator/region"
)

// World implements a world provider for the SurvivalCraft 2 world format
type World struct {
	dir     string
	mu      sync.Mutex
	regions map[[2]int32]*region.Region
}

// Open opens the world at worldDir.
// If the directory does not exist, a new world scaffold is created.
func Open(worldDir string) (*World, error) {
	info, err := os.Stat(worldDir)
	if os.IsNotExist(err) {
		if err := initWorldDir(worldDir); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("world.Open: %w", err)
	}

	if info == nil {
		info, err = os.Stat(worldDir)
		if err != nil {
			return nil, fmt.Errorf("world.Open: %w", err)
		}
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("world.Open: %q is not a directory", worldDir)
	}

	if err := ensureWorldSubDirs(worldDir); err != nil {
		return nil, err
	}
	if err := ensureProjectFile(worldDir); err != nil {
		return nil, err
	}

	return &World{
		dir:     worldDir,
		regions: make(map[[2]int32]*region.Region),
	}, nil
}

func initWorldDir(worldDir string) error {
	if err := os.MkdirAll(worldDir, 0755); err != nil {
		return fmt.Errorf("world.Open: create world dir: %w", err)
	}
	if err := ensureWorldSubDirs(worldDir); err != nil {
		return err
	}
	if err := ensureProjectFile(worldDir); err != nil {
		return err
	}
	return nil
}

func ensureWorldSubDirs(worldDir string) error {
	for _, dirName := range []string{"Regions", "PlayerEntities"} {
		path := filepath.Join(worldDir, dirName)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("world.Open: create %s: %w", dirName, err)
		}
	}
	return nil
}

func ensureProjectFile(worldDir string) error {
	projectPath := filepath.Join(worldDir, "Project.json")
	if _, err := os.Stat(projectPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("world.Open: stat project file: %w", err)
	}

	minimalProject := map[string]any{}
	data, err := json.Marshal(minimalProject)
	if err != nil {
		return fmt.Errorf("world.Open: marshal default project: %w", err)
	}
	if err := os.WriteFile(projectPath, data, 0644); err != nil {
		return fmt.Errorf("world.Open: create default project file: %w", err)
	}
	return nil
}

// Close closes all open region files.
func (w *World) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var firstErr error
	for k, r := range w.regions {
		if err := r.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(w.regions, k)
	}
	return firstErr
}

// LoadChunk reads and decodes the chunk at pos.
// Returns (nil, nil) if the chunk has not been generated yet.
func (w *World) LoadChunk(pos define.ChunkPos) (*chunk.Chunk, error) {
	r, lx, lz, err := w.regionFor(pos)
	if err != nil {
		return nil, err
	}
	return r.LoadChunk(lx, lz)
}

// LoadChunkPayloadOnly reads a raw compressed chunk payload without decoding.
// Returns (nil, nil) if the chunk has not been generated yet.
func (w *World) LoadChunkPayloadOnly(pos define.ChunkPos) ([]byte, error) {
	r, lx, lz, err := w.regionFor(pos)
	if err != nil {
		return nil, err
	}
	return r.LoadChunkPayloadOnly(lx, lz)
}

// SaveChunk encodes and writes c at pos.
func (w *World) SaveChunk(pos define.ChunkPos, c *chunk.Chunk) error {
	r, lx, lz, err := w.regionFor(pos)
	if err != nil {
		return err
	}
	return r.SaveChunk(lx, lz, c)
}

// SaveChunkPayloadOnly writes a raw compressed chunk payload without encoding.
func (w *World) SaveChunkPayloadOnly(pos define.ChunkPos, payload []byte) error {
	r, lx, lz, err := w.regionFor(pos)
	if err != nil {
		return err
	}
	return r.SaveChunkPayloadOnly(lx, lz, payload)
}

// regionFor returns the Region (opening/creating if necessary) and the local
// chunk coords for global chunk pos.
func (w *World) regionFor(pos define.ChunkPos) (*region.Region, int, int, error) {
	rx, rz := region.RegionCoords(pos)
	lx, lz := region.LocalCoords(pos)

	w.mu.Lock()
	defer w.mu.Unlock()

	key := [2]int32{rx, rz}
	r, ok := w.regions[key]
	if !ok {
		path := filepath.Join(w.dir, "Regions", fmt.Sprintf("Region %d,%d.dat", rx, rz))
		// Ensure Regions directory exists.
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, 0, 0, fmt.Errorf("world: mkdir regions: %w", err)
		}
		var err error
		r, err = region.Open(path)
		if err != nil {
			return nil, 0, 0, err
		}
		w.regions[key] = r
	}
	return r, lx, lz, nil
}
