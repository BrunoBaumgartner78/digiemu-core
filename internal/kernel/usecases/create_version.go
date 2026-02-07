package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type CreateVersionInput struct {
	UnitKey string
	Label   string
	Content string
}

type CreateVersionOutput struct {
	Version domain.Version
}

type CreateVersion struct {
	Repo ports.UnitRepository
}

func (uc CreateVersion) Execute(in CreateVersionInput) (CreateVersionOutput, error) {
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return CreateVersionOutput{}, err
	}
	if !ok {
		return CreateVersionOutput{}, domain.ErrUnitNotFound
	}

	v, err := domain.NewVersion(unit.ID, in.Label, in.Content)
	if err != nil {
		return CreateVersionOutput{}, err
	}

	if err := uc.Repo.SaveVersion(v); err != nil {
		return CreateVersionOutput{}, err
	}

	return CreateVersionOutput{Version: v}, nil
}
