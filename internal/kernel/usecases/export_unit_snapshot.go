package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type ExportUnitSnapshot struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLogByUnitReader // optional; required only if IncludeAudit=true
}

func (uc ExportUnitSnapshot) ExportUnitSnapshot(in ports.ExportUnitSnapshotRequest) (ports.ExportUnitSnapshotResponse, error) {
	u, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.ExportUnitSnapshotResponse{}, err
	}
	if !ok {
		return ports.ExportUnitSnapshotResponse{}, domain.ErrUnitNotFound
	}

	vs, err := uc.Repo.ListVersionsByUnitID(u.ID)
	if err != nil {
		return ports.ExportUnitSnapshotResponse{}, err
	}

	outVers := make([]ports.VersionDTO, 0, len(vs))
	for _, v := range vs {
		outVers = append(outVers, toVersionDTO(v))
	}

	resp := ports.ExportUnitSnapshotResponse{
		Unit: ports.UnitDTO{
			ID:            u.ID,
			Key:           u.Key,
			Title:         u.Title,
			Description:   u.Description,
			HeadVersionID: u.HeadVersionID,
		},
		Versions: outVers,
	}

	// v0.2.5: snapshot hash over unit + versions (canonical, deterministic)
	resp.SnapshotHash = sha256HexFromLines(snapshotCanonicalLines(resp.Unit, resp.Versions))

	if in.IncludeAudit {
		if uc.Audit == nil {
			return ports.ExportUnitSnapshotResponse{}, domain.ErrAuditNotConfigured
		}
		evs, err := uc.Audit.ListByUnitID(u.ID)
		if err != nil {
			return ports.ExportUnitSnapshotResponse{}, err
		}
		resp.Audit = evs

		// v0.2.5: audit hash (canonical, deterministic)
		audLines, err := auditCanonicalLines(evs)
		if err != nil {
			return ports.ExportUnitSnapshotResponse{}, err
		}
		resp.AuditHash = sha256HexFromLines(audLines)
	}

	return resp, nil
}
