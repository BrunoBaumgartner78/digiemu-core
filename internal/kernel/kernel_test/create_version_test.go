package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateVersion_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()

	createUnit := usecases.CreateUnit{Repo: repo}
	unitOut, err := createUnit.Execute(usecases.CreateUnitInput{
		Key:   "merkblatt-steuern",
		Title: "Merkblatt Steuern",
	})
	if err != nil {
		t.Fatalf("create unit err: %v", err)
	}

	uc := usecases.CreateVersion{Repo: repo}
	out, err := uc.Execute(usecases.CreateVersionInput{
		UnitKey: "merkblatt-steuern",
		Label:   "v1",
		Content: "Inhalt 1",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if out.Version.UnitID != unitOut.Unit.ID {
		t.Fatalf("expected unitID match")
	}

	vs, _ := repo.ListVersionsByUnitID(unitOut.Unit.ID)
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
}

func TestCreateVersion_UnitNotFound(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.CreateVersion{Repo: repo}

	_, err := uc.Execute(usecases.CreateVersionInput{
		UnitKey: "does-not-exist",
		Label:   "v1",
		Content: "x",
	})
	if err != domain.ErrUnitNotFound {
		t.Fatalf("expected ErrUnitNotFound, got %v", err)
	}
}

func TestCreateVersion_Validation(t *testing.T) {
	repo := memory.NewUnitRepo()
	createUnit := usecases.CreateUnit{Repo: repo}
	_, _ = createUnit.Execute(usecases.CreateUnitInput{Key: "abc", Title: "Title"})

	uc := usecases.CreateVersion{Repo: repo}

	_, err := uc.Execute(usecases.CreateVersionInput{UnitKey: "abc", Label: "", Content: "x"})
	if err != domain.ErrInvalidVersionLabel {
		t.Fatalf("expected ErrInvalidVersionLabel, got %v", err)
	}

	_, err = uc.Execute(usecases.CreateVersionInput{UnitKey: "abc", Label: "v1", Content: ""})
	if err != domain.ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}
