package ports

import "digiemu-core/internal/kernel/domain"

type UnitRepository interface {
	ExistsByKey(key string) (bool, error)
	SaveUnit(u domain.Unit) error
	FindUnitByKey(key string) (domain.Unit, bool, error)
	FindUnitByID(id string) (domain.Unit, bool, error)

	SaveVersion(v domain.Version) error
	ListVersionsByUnitID(unitID string) ([]domain.Version, error)
}
