package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProjectInfo contains selected fields from Project.json.
// Project.json uses the same [typeName, value] pair format as
// ValuesDictionary MessagePack, but serialised as JSON (UTF-8 BOM).
type ProjectInfo struct {
	// Raw is the full parsed JSON object for fields not explicitly decoded.
	Raw map[string]json.RawMessage
}

// LoadProject reads and parses Project.json from the world directory.
func (w *World) LoadProject() (*ProjectInfo, error) {
	path := filepath.Join(w.dir, "Project.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("world.LoadProject: %w", err)
	}
	// Strip UTF-8 BOM if present.
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("world.LoadProject: unmarshal: %w", err)
	}
	return &ProjectInfo{Raw: raw}, nil
}

// StringField decodes a string-type ValuesDictionary field from the project.
// SC2 JSON format: "Key": ["string", "value"]
func (p *ProjectInfo) StringField(key string) (string, bool) {
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

// IntField decodes an int-type ValuesDictionary field from the project.
// SC2 JSON format: "Key": ["int", number]
func (p *ProjectInfo) IntField(key string) (int, bool) {
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
