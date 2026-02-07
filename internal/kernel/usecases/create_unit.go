package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// CreateUnit is the usecase implementation.
type CreateUnit struct {
	Repo ports.UnitRepository
}

// CreateUnit implements ports.CreateUnitUsecase
func (uc CreateUnit) CreateUnit(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
	u, err := domain.NewUnit(in.Key, in.Title, in.Description)
	if err != nil {
		return ports.CreateUnitResponse{}, err
	}

	exists, err := uc.Repo.ExistsByKey(u.Key)
	if err != nil {
		return ports.CreateUnitResponse{}, err
	}
	if exists {
		return ports.CreateUnitResponse{}, domain.ErrUnitAlreadyExists
	}

	if err := uc.Repo.SaveUnit(u); err != nil {
		return ports.CreateUnitResponse{}, err
	}

	return ports.CreateUnitResponse{
		UnitID:      u.ID,
		Key:         u.Key,
		Title:       u.Title,
		Description: u.Description,
	}, nil
}
