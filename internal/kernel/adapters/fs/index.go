package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// indexStore keeps a persistent mapping from unit key -> unit id.
// Stored on disk at: <data>/index/units_by_key.json
type indexStore struct {
	mu   sync.Mutex
	base string

	// loaded indicates whether the index has been loaded (even if empty)
	loaded bool

	// in-memory maps
	unitIDByKey map[string]string
}

type unitsByKeyFile struct {
	Schema string            `json:"schema"`
	Keys   map[string]string `json:"keys"` // key -> unitID
}

const unitsByKeySchema = "digiemu.index.units_by_key.v1"

func newIndexStore(basePath string) *indexStore {
	return &indexStore{
		base:        basePath,
		loaded:      false,
		unitIDByKey: map[string]string{},
	}
}

func (s *indexStore) indexDir() string {
	return filepath.Join(s.base, "index")
}

func (s *indexStore) unitsByKeyPath() string {
	return filepath.Join(s.indexDir(), "units_by_key.json")
}

// ensureLoaded loads the index; if missing/corrupt, it rebuilds from disk.
// safe to call often (fast-path if already loaded).
func (s *indexStore) ensureLoaded(rebuild func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// IMPORTANT: empty index is still a valid loaded state.
	if s.loaded {
		return nil
	}

	// 1) Try load. If it succeeds, mark loaded and stop.
	if err := s.loadLocked(); err == nil {
		s.loaded = true
		return nil
	}

	// 2) If load failed and we have a rebuild, try rebuild then load once.
	if rebuild != nil {
		if err := rebuild(); err != nil {
			return err
		}
		if err := s.loadLocked(); err != nil {
			return err
		}
		s.loaded = true
		return nil
	}

	// 3) No rebuild available -> return the load error from above by re-loading once
	// (keeps behavior simple and deterministic).
	if err := s.loadLocked(); err != nil {
		return err
	}
	s.loaded = true
	return nil
}

// loadLocked attempts to load the index file into memory. It must be called while
// holding s.mu or it will lock internally. It guarantees s.unitIDByKey is non-nil
// on success (may be empty map).
func (s *indexStore) loadLocked() error {
	// reuse loadUnitsByKey to parse and validate file, then assign
	m, err := s.loadUnitsByKey()
	if err != nil {
		return err
	}
	if m == nil {
		s.unitIDByKey = make(map[string]string)
	} else {
		s.unitIDByKey = m
	}
	return nil
}

// getUnitIDByKey returns (id, ok). Call ensureLoaded beforehand if you want auto-rebuild.
func (s *indexStore) getUnitIDByKey(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.unitIDByKey[key]
	return id, ok
}

// listUnitIDs returns sorted unit IDs (stable order).
func (s *indexStore) listUnitIDs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := make([]string, 0, len(s.unitIDByKey))
	for _, id := range s.unitIDByKey {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// upsertUnitKey updates index in memory and persists best-effort.
func (s *indexStore) upsertUnitKey(key, unitID string) {
	s.mu.Lock()
	if s.unitIDByKey == nil {
		s.unitIDByKey = map[string]string{}
	}
	s.unitIDByKey[key] = unitID
	snapshot := make(map[string]string, len(s.unitIDByKey))
	for k, v := range s.unitIDByKey {
		snapshot[k] = v
	}
	s.mu.Unlock()

	// Best-effort write; errors are not fatal for core ops.
	_ = s.saveUnitsByKeyAtomic(snapshot)
}

func (s *indexStore) loadUnitsByKey() (map[string]string, error) {
	p := s.unitsByKeyPath()
	b, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, errors.New("index file missing")
	}
	if err != nil {
		return nil, err
	}

	// Allow empty/partial file to be treated as corrupt
	if len(strings.TrimSpace(string(b))) == 0 {
		return nil, errors.New("index file empty")
	}

	var f unitsByKeyFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("index json invalid: %w", err)
	}
	if f.Schema != unitsByKeySchema {
		return nil, fmt.Errorf("index schema mismatch: %s", f.Schema)
	}
	if f.Keys == nil {
		return nil, errors.New("index keys missing")
	}

	// Validate content: non-empty keys/ids
	out := make(map[string]string, len(f.Keys))
	for k, id := range f.Keys {
		k = strings.TrimSpace(k)
		id = strings.TrimSpace(id)
		if k == "" || id == "" {
			return nil, errors.New("index contains empty key or id")
		}
		out[k] = id
	}
	return out, nil
}

func (s *indexStore) saveUnitsByKeyAtomic(m map[string]string) error {
	if err := os.MkdirAll(s.indexDir(), 0o755); err != nil {
		return err
	}

	payload := unitsByKeyFile{
		Schema: unitsByKeySchema,
		Keys:   m,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	dst := s.unitsByKeyPath()
	tmp := dst + ".tmp"

	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}
