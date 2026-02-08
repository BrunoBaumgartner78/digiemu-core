package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/domain"
)

func TestIndex_Rebuild_WhenMissing(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	// Save a unit without any index existing
	u := domain.Unit{ID: "unit_1", Key: "abc", Title: "Title"}
	if err := repo.SaveUnit(u); err != nil {
		t.Fatalf("SaveUnit: %v", err)
	}

	// Remove index folder to force rebuild path
	_ = os.RemoveAll(filepath.Join(dir, "index"))

	got, ok, err := repo.FindUnitByKey("abc")
	if err != nil {
		t.Fatalf("FindUnitByKey: %v", err)
	}
	if !ok || got.ID != "unit_1" {
		t.Fatalf("expected unit_1, got ok=%v id=%s", ok, got.ID)
	}
}

func TestIndex_Rebuild_WhenCorrupt(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	u := domain.Unit{ID: "unit_1", Key: "abc", Title: "Title"}
	if err := repo.SaveUnit(u); err != nil {
		t.Fatalf("SaveUnit: %v", err)
	}

	// Write a corrupt index file
	idxDir := filepath.Join(dir, "index")
	if err := os.MkdirAll(idxDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	idxPath := filepath.Join(idxDir, "units_by_key.json")
	if err := os.WriteFile(idxPath, []byte("{ this is not json"), 0o644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}

	got, ok, err := repo.FindUnitByKey("abc")
	if err != nil {
		t.Fatalf("FindUnitByKey: %v", err)
	}
	if !ok || got.ID != "unit_1" {
		t.Fatalf("expected unit_1 after rebuild, got ok=%v id=%s", ok, got.ID)
	}
}

func TestIndex_ListUnits_Works_WhenIndexEmpty(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	// No units, no index -> ListUnits should not error
	us, err := repo.ListUnits()
	if err != nil {
		t.Fatalf("ListUnits: %v", err)
	}
	if len(us) != 0 {
		t.Fatalf("expected 0 units, got %d", len(us))
	}
}
