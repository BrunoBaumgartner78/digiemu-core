package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type ListVersions struct {
	Repo ports.UnitRepository
}

func (uc ListVersions) ListVersions(in ports.ListVersionsRequest) (ports.ListVersionsResponse, error) {
	u, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.ListVersionsResponse{}, err
	}
	if !ok {
		return ports.ListVersionsResponse{}, domain.ErrUnitNotFound
	}

	vs, err := uc.Repo.ListVersionsByUnitID(u.ID)
	if err != nil {
		return ports.ListVersionsResponse{}, err
	}

	out := make([]ports.VersionDTO, 0, len(vs))
	if !in.NewestFirst {
		for _, v := range vs {
			out = append(out, toVersionDTO(v))
		}
	} else {
		for i := len(vs) - 1; i >= 0; i-- {
			out = append(out, toVersionDTO(vs[i]))
		}
	}

	return ports.ListVersionsResponse{
		UnitID:   u.ID,
		Versions: out,
	}, nil
}

func toVersionDTO(v domain.Version) ports.VersionDTO {
	return ports.VersionDTO{
		ID:            v.ID,
		UnitID:        v.UnitID,
		Label:         v.Label,
		Content:       v.Content,
		PrevVersionID: v.PrevVersionID,
		ContentHash:   v.ContentHash,
		CreatedAtUnix: v.CreatedAtUnix,
		ActorID:       v.ActorID,
	}
}
