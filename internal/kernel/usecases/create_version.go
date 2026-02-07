package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// CreateVersion is the usecase implementation.
type CreateVersion struct {
	Repo ports.UnitRepository
}

// CreateVersion implements ports.CreateVersionUsecase
func (uc CreateVersion) CreateVersion(in ports.CreateVersionRequest) (ports.CreateVersionResponse, error) {
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.CreateVersionResponse{}, err
	}
	if !ok {
		return ports.CreateVersionResponse{}, domain.ErrUnitNotFound
	}

	v, err := domain.NewVersion(unit.ID, in.Label, in.Content)
	if err != nil {
		return ports.CreateVersionResponse{}, err
	}

	if err := uc.Repo.SaveVersion(v); err != nil {
		return ports.CreateVersionResponse{}, err
	}

	return ports.CreateVersionResponse{
		VersionID: v.ID,
		UnitID:    v.UnitID,
		Label:     v.Label,
		Content:   v.Content,
	}, nil
}
