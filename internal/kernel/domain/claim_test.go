package domain

import "testing"

func TestClaimSet_ValidateMinimal_Valid(t *testing.T) {
	cs := &ClaimSet{
		SchemaVersion: ClaimSetSchemaV0,
		VersionID:     "ver_1",
		Claims: []Claim{
			{ID: "c1", Text: "First claim"},
			{ID: "c2", Text: "Second claim"},
		},
		Relations: []ClaimRelation{
			{Type: RelationContradicts, FromClaimID: "c1", ToClaimID: "c2"},
		},
	}
	if err := cs.ValidateMinimal(); err != nil {
		t.Fatalf("expected valid claimset, got error: %v", err)
	}
}

func TestClaimSet_ValidateMinimal_InvalidSchema(t *testing.T) {
	cs := &ClaimSet{SchemaVersion: "wrong", VersionID: "ver_1"}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for wrong schema_version")
	}
}

func TestClaimSet_ValidateMinimal_MissingFields(t *testing.T) {
	cs := &ClaimSet{SchemaVersion: ClaimSetSchemaV0}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for missing version_id")
	}

	cs = &ClaimSet{
		SchemaVersion: ClaimSetSchemaV0,
		VersionID:     "v",
		Claims:        []Claim{{ID: "", Text: ""}},
	}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for missing claim id/text")
	}
}

func TestClaimSet_ValidateMinimal_DuplicateIDs(t *testing.T) {
	cs := &ClaimSet{
		SchemaVersion: ClaimSetSchemaV0,
		VersionID:     "v",
		Claims: []Claim{
			{ID: "c1", Text: "A"},
			{ID: "c1", Text: "B"},
		},
	}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for duplicate claim id")
	}
}

func TestClaimSet_ValidateMinimal_InvalidRelation(t *testing.T) {
	cs := &ClaimSet{
		SchemaVersion: ClaimSetSchemaV0,
		VersionID:     "v",
		Claims:        []Claim{{ID: "c1", Text: "A"}},
		Relations:     []ClaimRelation{{Type: RelationContradicts, FromClaimID: "c1", ToClaimID: "missing"}},
	}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for relation referencing unknown claim id")
	}

	cs = &ClaimSet{
		SchemaVersion: ClaimSetSchemaV0,
		VersionID:     "v",
		Claims:        []Claim{{ID: "c1", Text: "A"}, {ID: "c2", Text: "B"}},
		Relations:     []ClaimRelation{{Type: RelationType("UNKNOWN"), FromClaimID: "c1", ToClaimID: "c2"}},
	}
	if err := cs.ValidateMinimal(); err == nil {
		t.Fatal("expected error for unsupported relation type")
	}
}
