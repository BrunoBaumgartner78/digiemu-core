package memory

import (
	"digiemu-core/internal/kernel/domain"
	"sync"
)

type UnitRepo struct {
	mu sync.RWMutex

	unitsByID  map[string]domain.Unit
	unitsByKey map[string]string // key -> unitID

	versionsByUnitID map[string][]domain.Version
}

func NewUnitRepo() *UnitRepo {
	return &UnitRepo{
		unitsByID:        map[string]domain.Unit{},
		unitsByKey:       map[string]string{},
		versionsByUnitID: map[string][]domain.Version{},
	}
}

func (r *UnitRepo) ExistsByKey(key string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.unitsByKey[key]
	return ok, nil
}

func (r *UnitRepo) SaveUnit(u domain.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.unitsByID[u.ID] = u
	r.unitsByKey[u.Key] = u.ID
	return nil
}

func (r *UnitRepo) FindUnitByKey(key string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.unitsByKey[key]
	if !ok {
		return domain.Unit{}, false, nil
	}
	u, ok := r.unitsByID[id]
	if !ok {
		return domain.Unit{}, false, nil
	}
	return u, true, nil
}

func (r *UnitRepo) FindUnitByID(id string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.unitsByID[id]
	return u, ok, nil
}

func (r *UnitRepo) SaveVersion(v domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.versionsByUnitID[v.UnitID] = append(r.versionsByUnitID[v.UnitID], v)
	return nil
}

func (r *UnitRepo) ListVersionsByUnitID(unitID string) ([]domain.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vs := r.versionsByUnitID[unitID]
	out := make([]domain.Version, len(vs))
	copy(out, vs)
	return out, nil
}
