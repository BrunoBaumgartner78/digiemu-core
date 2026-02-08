package usecases

import (
	"encoding/json"
	"errors"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// SetClaims implements persisting a ClaimSet for a specific version and
// emitting a CLAIM_SET audit event. It validates schema and referential
// integrity using domain.ValidateMinimal.
type SetClaims struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

func (uc SetClaims) SetClaims(in ports.SetClaimsRequest) (ports.SetClaimsResponse, error) {
	if uc.Repo == nil {
		return ports.SetClaimsResponse{}, domain.ErrUnitNotFound
	}
	if uc.Audit == nil {
		return ports.SetClaimsResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.SetClaimsResponse{}, domain.ErrClockNotConfigured
	}

	// resolve unit by key
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if !ok {
		return ports.SetClaimsResponse{}, domain.ErrUnitNotFound
	}

	// determine target version
	verID := in.VersionID
	if verID == "" {
		verID = unit.HeadVersionID
	}
	if verID == "" {
		return ports.SetClaimsResponse{}, domain.ErrVersionNotFound
	}

	// validate version exists
	_, found, err := uc.Repo.FindVersionByID(verID)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if !found {
		return ports.SetClaimsResponse{}, domain.ErrVersionNotFound
	}

	// max size check (reuse 64KB like meaning)
	if len(in.BodyBytes) > 64*1024 {
		return ports.SetClaimsResponse{}, errors.New("claimset.json too large")
	}

	// unmarshal into domain.ClaimSet and validate schema_version
	var cs domain.ClaimSet
	if err := json.Unmarshal(in.BodyBytes, &cs); err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if cs.SchemaVersion != "claimset/v0" {
		return ports.SetClaimsResponse{}, errors.New("unsupported schema_version")
	}

	// validate minimal referential integrity
	if err := cs.ValidateMinimal(); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	// compute canonical hash
	ch, err := ComputeClaimSetHashFromStruct(cs)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}

	// persist via repo (persistence-only)
	if err := uc.Repo.SaveClaimSet(unit.ID, verID, cs, ch); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	// append audit event
	ev := domain.AuditEvent{
		Schema:    "digiemu.audit.v1",
		ID:        domain.NewID("evt"),
		Type:      "CLAIM_SET",
		AtUnix:    uc.Clock.NowUnix(),
		ActorID:   in.ActorID,
		UnitID:    unit.ID,
		VersionID: verID,
		Data: domain.ClaimSetData{
			UnitID:       unit.ID,
			VersionID:    verID,
			ClaimSetHash: ch,
			ClaimSetPath: unit.ID + "." + verID + ".claimset.json",
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	return ports.SetClaimsResponse{UnitID: unit.ID, VersionID: verID, ClaimSetHash: ch}, nil
}
