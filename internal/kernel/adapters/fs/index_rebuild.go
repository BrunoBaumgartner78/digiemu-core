package fs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// rebuildUnitsByKey scans <data>/units/*.json and rebuilds key->id map.
// It is intentionally forgiving: skips unreadable/broken unit files.
func rebuildUnitsByKey(basePath string) (map[string]string, error) {
	unitsDir := filepath.Join(basePath, "units")
	entries, err := os.ReadDir(unitsDir)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		p := filepath.Join(unitsDir, e.Name())
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var rec UnitRecord
		if err := json.Unmarshal(b, &rec); err != nil {
			continue
		}
		key := strings.TrimSpace(rec.Key)
		id := strings.TrimSpace(rec.ID)
		if key == "" || id == "" {
			continue
		}
		out[key] = id
	}
	return out, nil
}
