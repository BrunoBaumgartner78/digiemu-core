package fs_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/domain"
)

func TestFSRepo_CreateUnit_SaveLoad_AddVersion(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	u, err := domain.NewUnit("u-key", "Unit Title", "")
	if err != nil {
		t.Fatalf("new unit err: %v", err)
	}

	if err := repo.SaveUnit(u); err != nil {
		t.Fatalf("save unit err: %v", err)
	}

	// file exists and valid JSON
	p := filepath.Join(dir, "units", u.ID+".json")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("read file err: %v", err)
	}
	var rec map[string]interface{}
	if err := json.Unmarshal(b, &rec); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	// find by key
	got, ok, err := repo.FindUnitByKey("u-key")
	if err != nil || !ok {
		t.Fatalf("find by key failed: %v %v", ok, err)
	}
	if got.ID != u.ID {
		t.Fatalf("unexpected id")
	}

	// add version
	v, err := domain.NewVersion(u.ID, "v1", "content")
	if err != nil {
		t.Fatalf("new version err: %v", err)
	}
	if err := repo.SaveVersion(v); err != nil {
		t.Fatalf("save version err: %v", err)
	}

	vs, err := repo.ListVersionsByUnitID(u.ID)
	if err != nil {
		t.Fatalf("list versions err: %v", err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
}

func TestFSRepo_FindUnitByID_NotFound(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	_, ok, err := repo.FindUnitByID("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ok {
		t.Fatalf("expected not found")
	}
}
