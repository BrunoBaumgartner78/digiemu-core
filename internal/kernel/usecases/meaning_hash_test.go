package usecases

import (
	"encoding/json"
	"testing"

	"digiemu-core/internal/kernel/domain"
)

func TestComputeMeaningHash_Stability(t *testing.T) {
	a := `{"schema_version":"meaning/v1","title":"T","purpose":"P","scope":{"audience":["dev"]}}`
	b := `{
  "purpose":"P",
  "schema_version":"meaning/v1",
  "scope": { "audience": [ "dev" ] },
  "title":"T"
}`
	var ma, mb domain.Meaning
	if err := json.Unmarshal([]byte(a), &ma); err != nil {
		t.Fatalf("unmarshal a: %v", err)
	}
	if err := json.Unmarshal([]byte(b), &mb); err != nil {
		t.Fatalf("unmarshal b: %v", err)
	}
	ha, err := ComputeMeaningHash(ma)
	if err != nil {
		t.Fatalf("compute a: %v", err)
	}
	hb, err := ComputeMeaningHash(mb)
	if err != nil {
		t.Fatalf("compute b: %v", err)
	}
	if ha != hb {
		t.Fatalf("expected identical hashes, got %s != %s", ha, hb)
	}
}
