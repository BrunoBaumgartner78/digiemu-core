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
	versionsByID     map[string]domain.Version // v0.2.3: fast lookup
	meanings         map[string]domain.Meaning // key: unitID.versionID
	claimsets        map[string]domain.ClaimSet
	uncertainties    map[string]domain.Uncertainty
}

func NewUnitRepo() *UnitRepo {
	return &UnitRepo{
		unitsByID:        map[string]domain.Unit{},
		unitsByKey:       map[string]string{},
		versionsByUnitID: map[string][]domain.Version{},
		versionsByID:     map[string]domain.Version{},
		meanings:         map[string]domain.Meaning{},
		claimsets:        map[string]domain.ClaimSet{},
		uncertainties:    map[string]domain.Uncertainty{},
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

func (r *UnitRepo) ListUnits() ([]domain.Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Unit, 0, len(r.unitsByID))
	for _, u := range r.unitsByID {
		out = append(out, u)
	}
	return out, nil
}

func (r *UnitRepo) SaveVersion(v domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.versionsByUnitID[v.UnitID] = append(r.versionsByUnitID[v.UnitID], v)
	r.versionsByID[v.ID] = v
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

func (r *UnitRepo) UpdateUnitHead(unitID, headVersionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.unitsByID[unitID]
	if !ok {
		return domain.ErrUnitNotFound
	}
	u.HeadVersionID = headVersionID
	r.unitsByID[unitID] = u
	return nil
}

func (r *UnitRepo) SaveMeaning(unitID, versionID string, meaning domain.Meaning, meaningHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ensure version exists
	v, ok := r.versionsByID[versionID]
	if !ok {
		return domain.ErrVersionNotFound
	}
	// update version record
	v.MeaningHash = meaningHash
	r.versionsByID[versionID] = v

	// update in versionsByUnitID slice
	vs := r.versionsByUnitID[unitID]
	for i := range vs {
		if vs[i].ID == versionID {
			vs[i].MeaningHash = meaningHash
			r.versionsByUnitID[unitID] = vs
			break
		}
	}

	// store meaning in memory map
	key := unitID + "." + versionID
	r.meanings[key] = meaning
	return nil
}

func (r *UnitRepo) LoadMeaning(unitID, versionID string) (domain.Meaning, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := unitID + "." + versionID
	m, ok := r.meanings[key]
	return m, ok, nil
}

func (r *UnitRepo) SaveClaimSet(unitID, versionID string, claimSet domain.ClaimSet, claimSetHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ensure version exists
	v, ok := r.versionsByID[versionID]
	if !ok {
		return domain.ErrVersionNotFound
	}
	// update version record
	v.ClaimSetHash = claimSetHash
	r.versionsByID[versionID] = v

	// update in versionsByUnitID slice
	vs := r.versionsByUnitID[unitID]
	for i := range vs {
		if vs[i].ID == versionID {
			vs[i].ClaimSetHash = claimSetHash
			r.versionsByUnitID[unitID] = vs
			break
		}
	}

	// store claimset in memory map
	key := unitID + "." + versionID
	r.claimsets[key] = claimSet
	return nil
}

func (r *UnitRepo) SaveUncertainty(unitID, versionID string, u domain.Uncertainty, uncertaintyHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ensure version exists
	v, ok := r.versionsByID[versionID]
	if !ok {
		return domain.ErrVersionNotFound
	}
	// update version record
	v.UncertaintyHash = uncertaintyHash
	r.versionsByID[versionID] = v

	// update in versionsByUnitID slice
	vs := r.versionsByUnitID[unitID]
	for i := range vs {
		if vs[i].ID == versionID {
			vs[i].UncertaintyHash = uncertaintyHash
			r.versionsByUnitID[unitID] = vs
			break
		}
	}

	// store uncertainty in memory map
	key := unitID + "." + versionID
	r.uncertainties[key] = u
	return nil
}

func (r *UnitRepo) LoadUncertainty(unitID, versionID string) (domain.Uncertainty, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := unitID + "." + versionID
	u, ok := r.uncertainties[key]
	return u, ok, nil
}

func (r *UnitRepo) LoadClaimSet(unitID, versionID string) (domain.ClaimSet, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := unitID + "." + versionID
	cs, ok := r.claimsets[key]
	return cs, ok, nil
}
