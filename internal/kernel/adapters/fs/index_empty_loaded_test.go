package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexEnsureLoaded_DoesNotRebuildRepeatedly_WhenNoUnits(t *testing.T) {
	dir := t.TempDir()

	// Create empty data layout (no units)
	if err := os.MkdirAll(filepath.Join(dir, "index"), 0o755); err != nil {
		t.Fatalf("mkdir index: %v", err)
	}
	// Write a valid, empty index file (schema OK, no entries).
	idxPath := filepath.Join(dir, "index", "units_by_key.json")
	empty := []byte(`{"schema":"digiemu.index.units_by_key.v1","keys":{}}`)
	if err := os.WriteFile(idxPath, empty, 0o644); err != nil {
		t.Fatalf("write empty index: %v", err)
	}

	s := newIndexStore(dir)

	rebuildCalls := 0
	rebuild := func() error { rebuildCalls++; return nil }

	// First ensureLoaded should load (no rebuild needed)
	if err := s.ensureLoaded(rebuild); err != nil {
		t.Fatalf("ensureLoaded #1: %v", err)
	}
	// Second ensureLoaded must not rebuild again (loaded=true)
	if err := s.ensureLoaded(rebuild); err != nil {
		t.Fatalf("ensureLoaded #2: %v", err)
	}

	if rebuildCalls != 0 {
		t.Fatalf("expected rebuildCalls=0, got %d", rebuildCalls)
	}
}

func TestIndexEnsureLoaded_RebuildsOnce_WhenMissingAndNoUnits(t *testing.T) {
	dir := t.TempDir()

	// No index file exists here
	s := newIndexStore(dir)

	rebuildCalls := 0
	rebuild := func() error {
		rebuildCalls++
		// Rebuild should create the index file even if there are 0 units.
		if err := os.MkdirAll(filepath.Join(dir, "index"), 0o755); err != nil {
			return err
		}
		idxPath := filepath.Join(dir, "index", "units_by_key.json")
		return os.WriteFile(idxPath, []byte(`{"schema":"digiemu.index.units_by_key.v1","keys":{}}`), 0o644)
	}

	if err := s.ensureLoaded(rebuild); err != nil {
		t.Fatalf("ensureLoaded #1: %v", err)
	}
	if err := s.ensureLoaded(rebuild); err != nil {
		t.Fatalf("ensureLoaded #2: %v", err)
	}

	if rebuildCalls != 1 {
		t.Fatalf("expected rebuildCalls=1, got %d", rebuildCalls)
	}
}
