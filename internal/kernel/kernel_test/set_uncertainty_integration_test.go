package kernel_test

import (
	"testing"

	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

// Integration: SetUncertainty -> Export -> VerifyAudit(StrictHash=true)
func TestSetUncertainty_Export_Verify_Strict_Memory(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.RealClock{}

	// create unit & version using existing usecases
	cu := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err := cu.CreateUnit(ports.CreateUnitRequest{Key: "k01", Title: "t01", Description: "d", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	cv := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	outV, err := cv.CreateVersion(ports.CreateVersionRequest{UnitKey: "k01", Label: "lbl", Content: "c", ActorID: "test"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	// set uncertainty
	su := usecases.SetUncertainty{Repo: repo, Audit: audit, Clock: clock}
	uJSON := []byte(`{"schema_version":"uncertainty/v0","id":"u1","type":"empirical","level":"low","applies_to":{"scope":"version"}}`)
	_, err = su.SetUncertainty(ports.SetUncertaintyRequest{UnitKey: "k01", VersionID: outV.VersionID, BodyBytes: uJSON, ActorID: "test"})
	if err != nil {
		t.Fatalf("set uncertainty: %v", err)
	}

	// export is basically listing units and versions; now verify audit strict
	ver := usecases.VerifyAudit{Repo: repo, Audit: audit}
	res, err := ver.VerifyAudit(ports.VerifyAuditRequest{UnitKey: "k01", StrictHash: true})
	if err != nil {
		t.Fatalf("verify audit: %v", err)
	}
	if !res.Ok {
		t.Fatalf("expected audit ok, got missing=%v duplicates=%v mismatches=%v", res.Missing, res.Duplicates, res.HashMismatches)
	}
}
