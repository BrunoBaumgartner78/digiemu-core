package usecases

import (
	"testing"
)

func TestComputeClaimSetHash_Determinism(t *testing.T) {
	a := []byte(`{
  "schema_version": "claimset/v0",
  "version_id": "ver_1",
  "claims": [
    {"id":"c1","text":"First"},
    {"id":"c2","text":"Second"}
  ]
}`)

	b := []byte(`{"version_id":"ver_1","claims":[{"text":"First","id":"c1"},{"id":"c2","text":"Second"}],"schema_version":"claimset/v0"}`)

	ha, err := ComputeClaimSetHash(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hb, err := ComputeClaimSetHash(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ha != hb {
		t.Fatalf("hash mismatch: %s != %s", ha, hb)
	}
}
