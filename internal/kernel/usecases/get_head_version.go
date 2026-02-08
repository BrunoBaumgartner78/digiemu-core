package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type GetHeadVersion struct {
	Repo ports.UnitRepository
}

func (uc GetHeadVersion) GetHeadVersion(in ports.GetHeadVersionRequest) (ports.GetHeadVersionResponse, error) {
	u, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.GetHeadVersionResponse{}, err
	}
	if !ok {
		return ports.GetHeadVersionResponse{}, domain.ErrUnitNotFound
	}
	if u.HeadVersionID == "" {
		return ports.GetHeadVersionResponse{}, domain.ErrNoVersions
	}

	vs, err := uc.Repo.ListVersionsByUnitID(u.ID)
	if err != nil {
		return ports.GetHeadVersionResponse{}, err
	}
	if len(vs) == 0 {
		return ports.GetHeadVersionResponse{}, domain.ErrNoVersions
	}

	// find head in list (simple + stable; no new repo method needed yet)
	var head *domain.Version
	for i := range vs {
		if vs[i].ID == u.HeadVersionID {
			head = &vs[i]
			break
		}
	}
	if head == nil {
		// Inconsistent repository: head pointer exists but cannot be resolved.
		return ports.GetHeadVersionResponse{}, domain.ErrInconsistentHead
	}

	return ports.GetHeadVersionResponse{
		UnitID:  u.ID,
		Version: toVersionDTO(*head),
	}, nil
}
