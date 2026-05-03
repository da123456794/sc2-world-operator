package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PlayerEntity represents a parsed player entity from PlayerEntities/<guid>.json.
// SC2 player JSON uses the same [typeName, value] pair format as
// ValuesDictionary, but stored as a nested JSON object.
type PlayerEntity struct {
	// GUID is the filename stem (the player's permanent GUID).
	GUID string
	// Raw is the parsed "Entity" object for arbitrary field access.
	Raw map[string]json.RawMessage
}

// LoadPlayerEntities reads all files under PlayerEntities/ and returns
// one PlayerEntity per file.
func (w *World) LoadPlayerEntities() ([]*PlayerEntity, error) {
	dir := filepath.Join(w.dir, "PlayerEntities")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("world.LoadPlayerEntities: %w", err)
	}

	var players []*PlayerEntity
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		guid := e.Name()[:len(e.Name())-5] // strip ".json"
		p, err := w.loadPlayerEntity(filepath.Join(dir, e.Name()), guid)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (w *World) loadPlayerEntity(path, guid string) (*PlayerEntity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("world: read player %q: %w", guid, err)
	}
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	// Top-level: {"Entity": {...}}
	var top struct {
		Entity map[string]json.RawMessage `json:"Entity"`
	}
	if err := json.Unmarshal(data, &top); err != nil {
		return nil, fmt.Errorf("world: parse player %q: %w", guid, err)
	}
	return &PlayerEntity{GUID: guid, Raw: top.Entity}, nil
}

// StringField decodes a string ValuesDictionary entry from the entity.
// SC2 JSON format for a string field: "Key": ["string", "value"]
func (p *PlayerEntity) StringField(key string) (string, bool) {
	raw, ok := p.Raw[key]
	if !ok {
		return "", false
	}
	var pair [2]json.RawMessage
	if err := json.Unmarshal(raw, &pair); err != nil {
		return "", false
	}
	var s string
	if err := json.Unmarshal(pair[1], &s); err != nil {
		return "", false
	}
	return s, true
}

// IntField decodes an int ValuesDictionary entry from the entity.
func (p *PlayerEntity) IntField(key string) (int, bool) {
	raw, ok := p.Raw[key]
	if !ok {
		return 0, false
	}
	var pair [2]json.RawMessage
	if err := json.Unmarshal(raw, &pair); err != nil {
		return 0, false
	}
	var n int
	if err := json.Unmarshal(pair[1], &n); err != nil {
		return 0, false
	}
	return n, true
}
