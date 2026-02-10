package usecases

import (
	"testing"
)

func TestComputeUncertaintyHash_Determinism(t *testing.T) {
	a := []byte(`{
  "schema_version": "uncertainty/v0",
  "id": "u1",
  "type": "empirical",
  "level": "low",
  "applies_to": {"scope":"version"}
}`)

	b := []byte(`{"applies_to":{"scope":"version"},"id":"u1","level":"low","schema_version":"uncertainty/v0","type":"empirical"}`)

	ha, err := ComputeUncertaintyHash(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hb, err := ComputeUncertaintyHash(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ha != hb {
		t.Fatalf("hash mismatch: %s != %s", ha, hb)
	}
}
