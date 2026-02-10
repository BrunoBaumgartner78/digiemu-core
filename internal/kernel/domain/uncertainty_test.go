package domain

import "testing"

func TestUncertainty_ValidateMinimal_OK(t *testing.T) {
	u := Uncertainty{
		SchemaVersion: UncertaintySchemaV0,
		ID:            "u1",
		Type:          "empirical",
		Level:         "low",
		AppliesTo:     UncertaintyAppliesTo{Scope: ScopeVersion},
	}
	if err := u.ValidateMinimal(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUncertainty_ValidateMinimal_MissingClaimID(t *testing.T) {
	u := Uncertainty{
		SchemaVersion: UncertaintySchemaV0,
		ID:            "u2",
		Type:          "empirical",
		Level:         "low",
		AppliesTo:     UncertaintyAppliesTo{Scope: ScopeClaim},
	}
	if err := u.ValidateMinimal(); err == nil {
		t.Fatalf("expected error for missing claim id")
	}
}

func TestUncertainty_ValidateMinimal_WrongSchema(t *testing.T) {
	u := Uncertainty{
		SchemaVersion: "wrong/v0",
		ID:            "u3",
		Type:          "empirical",
		Level:         "low",
		AppliesTo:     UncertaintyAppliesTo{Scope: ScopeVersion},
	}
	if err := u.ValidateMinimal(); err == nil {
		t.Fatalf("expected error for wrong schema_version")
	}
}
