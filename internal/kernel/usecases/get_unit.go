package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type GetUnit struct {
	Repo ports.UnitRepository
}

func (uc GetUnit) GetUnit(in ports.GetUnitRequest) (ports.GetUnitResponse, error) {
	u, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.GetUnitResponse{}, err
	}
	if !ok {
		return ports.GetUnitResponse{}, domain.ErrUnitNotFound
	}
	return ports.GetUnitResponse{
		Unit: ports.UnitDTO{
			ID:            u.ID,
			Key:           u.Key,
			Title:         u.Title,
			Description:   u.Description,
			HeadVersionID: u.HeadVersionID,
		},
	}, nil
}
