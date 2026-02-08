package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type GetVersion struct {
	Repo ports.UnitRepository
}

func (uc GetVersion) GetVersion(in ports.GetVersionRequest) (ports.GetVersionResponse, error) {
	v, ok, err := uc.Repo.FindVersionByID(in.VersionID)
	if err != nil {
		return ports.GetVersionResponse{}, err
	}
	if !ok {
		return ports.GetVersionResponse{}, domain.ErrVersionNotFound
	}
	return ports.GetVersionResponse{Version: toVersionDTO(v)}, nil
}
