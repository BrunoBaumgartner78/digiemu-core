package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type CreateUnitInput struct {
	Key   string
	Title string
}

type CreateUnitOutput struct {
	Unit domain.Unit
}

type CreateUnit struct {
	Repo ports.UnitRepository
}

func (uc CreateUnit) Execute(in CreateUnitInput) (CreateUnitOutput, error) {
	u, err := domain.NewUnit(in.Key, in.Title)
	if err != nil {
		return CreateUnitOutput{}, err
	}

	exists, err := uc.Repo.ExistsByKey(u.Key)
	if err != nil {
		return CreateUnitOutput{}, err
	}
	if exists {
		return CreateUnitOutput{}, domain.ErrUnitAlreadyExists
	}

	if err := uc.Repo.SaveUnit(u); err != nil {
		return CreateUnitOutput{}, err
	}

	return CreateUnitOutput{Unit: u}, nil
}
