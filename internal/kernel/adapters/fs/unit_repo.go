package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"digiemu-core/internal/kernel/domain"
)

type UnitRepo struct {
	basePath string
	unitsDir string
	mu       sync.RWMutex

	// v0.2.7: persistent index (key -> unitID)
	index *indexStore
}

func NewUnitRepo(basePath string) *UnitRepo {
	units := filepath.Join(basePath, "units")
	_ = os.MkdirAll(units, 0o755)

	r := &UnitRepo{basePath: basePath, unitsDir: units}
	r.index = newIndexStore(basePath)
	return r
}

func (r *UnitRepo) unitPath(id string) string {
	return filepath.Join(r.unitsDir, id+".json")
}

func (r *UnitRepo) unitTempPath(id string) string {
	return filepath.Join(r.unitsDir, id+".json.tmp")
}

func (r *UnitRepo) ExistsByKey(key string) (bool, error) {
	// ensure index is loaded (or rebuilt) best-effort
	if r.index != nil {
		_ = r.index.ensureLoaded(func() error {
			m, err := rebuildUnitsByKey(r.basePath)
			if err != nil {
				return err
			}
			return r.index.saveUnitsByKeyAtomic(m)
		})
		_, ok := r.index.getUnitIDByKey(key)
		if ok {
			return true, nil
		}
	}

	// fallback: scan directory for unit with matching key
	files, err := ioutil.ReadDir(r.unitsDir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		b, err := ioutil.ReadFile(filepath.Join(r.unitsDir, f.Name()))
		if err != nil {
			return false, err
		}
		var ur UnitRecord
		if err := json.Unmarshal(b, &ur); err != nil {
			return false, err
		}
		if ur.Key == key {
			return true, nil
		}
	}
	return false, nil
}

func (r *UnitRepo) SaveUnit(u domain.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ur := UnitRecord{
		ID:            u.ID,
		Key:           u.Key,
		Title:         u.Title,
		Description:   u.Description,
		CreatedAt:     nowRFC3339(),
		HeadVersionID: u.HeadVersionID,
		Versions:      []VersionRecord{},
	}

	data, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}

	tmp := r.unitTempPath(u.ID)
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, r.unitPath(u.ID)); err != nil {
		return err
	}

	// update index best-effort
	if r.index != nil {
		r.index.upsertUnitKey(u.Key, u.ID)
	}
	return nil
}

func (r *UnitRepo) FindUnitByKey(key string) (domain.Unit, bool, error) {
	// ensure index is loaded (or rebuilt) best-effort
	if r.index != nil {
		_ = r.index.ensureLoaded(func() error {
			m, err := rebuildUnitsByKey(r.basePath)
			if err != nil {
				return err
			}
			return r.index.saveUnitsByKeyAtomic(m)
		})
		if id, ok := r.index.getUnitIDByKey(key); ok {
			if id != "" {
				return r.FindUnitByID(id)
			}
			// treat empty id as not found and fallthrough to scan
		}
	}

	// fallback: scan directory for unit with matching key
	files, err := ioutil.ReadDir(r.unitsDir)
	if err != nil {
		return domain.Unit{}, false, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		p := filepath.Join(r.unitsDir, f.Name())
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return domain.Unit{}, false, err
		}
		var ur UnitRecord
		if err := json.Unmarshal(b, &ur); err != nil {
			return domain.Unit{}, false, err
		}
		if ur.Key == key {
			return domain.Unit{
				ID:            ur.ID,
				Key:           ur.Key,
				Title:         ur.Title,
				Description:   ur.Description,
				HeadVersionID: ur.HeadVersionID,
			}, true, nil
		}
	}
	return domain.Unit{}, false, nil
}

func (r *UnitRepo) FindUnitByID(id string) (domain.Unit, bool, error) {
	p := r.unitPath(id)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.Unit{}, false, nil
	}
	if err != nil {
		return domain.Unit{}, false, err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return domain.Unit{}, false, err
	}
	return domain.Unit{
		ID:            ur.ID,
		Key:           ur.Key,
		Title:         ur.Title,
		Description:   ur.Description,
		HeadVersionID: ur.HeadVersionID,
	}, true, nil
}

func (r *UnitRepo) ListUnits() ([]domain.Unit, error) {
	// ensure index is loaded (or rebuilt) best-effort
	if r.index != nil {
		_ = r.index.ensureLoaded(func() error {
			m, err := rebuildUnitsByKey(r.basePath)
			if err != nil {
				return err
			}
			return r.index.saveUnitsByKeyAtomic(m)
		})
		ids := r.index.listUnitIDs()
		out := make([]domain.Unit, 0, len(ids))
		for _, id := range ids {
			u, ok, err := r.FindUnitByID(id)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			out = append(out, u)
		}
		return out, nil
	}

	files, err := ioutil.ReadDir(r.unitsDir)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Unit, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		p := filepath.Join(r.unitsDir, f.Name())
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var ur UnitRecord
		if err := json.Unmarshal(b, &ur); err != nil {
			return nil, err
		}
		out = append(out, domain.Unit{
			ID:            ur.ID,
			Key:           ur.Key,
			Title:         ur.Title,
			Description:   ur.Description,
			HeadVersionID: ur.HeadVersionID,
		})
	}
	return out, nil
}

func (r *UnitRepo) SaveVersion(v domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.unitPath(v.UnitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ErrUnitNotFound
	}
	if err != nil {
		return err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return err
	}

	vr := VersionRecord{
		ID:            v.ID,
		Label:         v.Label,
		Content:       v.Content,
		CreatedAt:     nowRFC3339(),
		PrevVersionID: v.PrevVersionID,
		ContentHash:   v.ContentHash,
		ActorID:       v.ActorID,
	}

	ur.Versions = append(ur.Versions, vr)

	data, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.unitTempPath(ur.ID)
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, p); err != nil {
		return err
	}
	return nil
}

func (r *UnitRepo) UpdateUnitHead(unitID, headVersionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.unitPath(unitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ErrUnitNotFound
	}
	if err != nil {
		return err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return err
	}

	ur.HeadVersionID = headVersionID

	data, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.unitTempPath(ur.ID)
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, p); err != nil {
		return err
	}
	return nil
}

func (r *UnitRepo) ListVersionsByUnitID(unitID string) ([]domain.Version, error) {
	p := r.unitPath(unitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return nil, err
	}
	out := make([]domain.Version, 0, len(ur.Versions))
	for _, vr := range ur.Versions {
		out = append(out, domain.Version{
			ID:              vr.ID,
			UnitID:          unitID,
			Label:           vr.Label,
			Content:         vr.Content,
			PrevVersionID:   vr.PrevVersionID,
			ContentHash:     vr.ContentHash,
			ActorID:         vr.ActorID,
			MeaningHash:     vr.MeaningHash,
			ClaimSetHash:    vr.ClaimSetHash,
			UncertaintyHash: vr.UncertaintyHash,
		})
	}
	return out, nil
}

// meaningPath returns the filesystem path for a unit's meaning.json file.
func (r *UnitRepo) meaningPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".meaning.json")
}

func (r *UnitRepo) meaningTempPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".meaning.json.tmp")
}

func (r *UnitRepo) claimsetPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".claimset.json")
}

func (r *UnitRepo) claimsetTempPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".claimset.json.tmp")
}

func (r *UnitRepo) uncertaintyPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".uncertainty.json")
}

func (r *UnitRepo) uncertaintyTempPath(unitID, versionID string) string {
	return filepath.Join(r.unitsDir, unitID+"."+versionID+".uncertainty.json.tmp")
}

// SaveMeaning stores a canonicalized meaning.json for the unit atomically.
// Returns ErrUnitNotFound if the unit does not exist.
func (r *UnitRepo) SaveMeaning(unitID, versionID string, m domain.Meaning, meaningHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.unitPath(unitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ErrUnitNotFound
	}
	if err != nil {
		return err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return err
	}

	// find version
	found := false
	for i := range ur.Versions {
		if ur.Versions[i].ID == versionID {
			ur.Versions[i].MeaningHash = meaningHash
			found = true
			break
		}
	}
	if !found {
		return domain.ErrVersionNotFound
	}

	// write meaning sidecar
	sidecar := filepath.Join(r.unitsDir, unitID+"."+versionID+".meaning.json")
	tmp := sidecar + ".tmp"
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, sidecar); err != nil {
		return err
	}

	// persist updated unit record
	out, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}
	tmpUnit := r.unitTempPath(ur.ID)
	if err := ioutil.WriteFile(tmpUnit, out, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpUnit, p); err != nil {
		return err
	}
	return nil
}

// SaveClaimSet stores a canonicalized claimset JSON for the unit atomically and
// updates the embedded version record's claimset_hash. Returns ErrUnitNotFound
// or ErrVersionNotFound where applicable.
func (r *UnitRepo) SaveClaimSet(unitID, versionID string, cs domain.ClaimSet, claimSetHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.unitPath(unitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ErrUnitNotFound
	}
	if err != nil {
		return err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return err
	}

	// find version
	found := false
	for i := range ur.Versions {
		if ur.Versions[i].ID == versionID {
			ur.Versions[i].ClaimSetHash = claimSetHash
			found = true
			break
		}
	}
	if !found {
		return domain.ErrVersionNotFound
	}

	// write claimset sidecar
	sidecar := r.claimsetPath(unitID, versionID)
	tmp := r.claimsetTempPath(unitID, versionID)
	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, sidecar); err != nil {
		return err
	}

	// persist updated unit record
	out, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}
	tmpUnit := r.unitTempPath(ur.ID)
	if err := ioutil.WriteFile(tmpUnit, out, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpUnit, p); err != nil {
		return err
	}
	return nil
}

// LoadClaimSet loads a version-scoped claimset sidecar if present.
func (r *UnitRepo) LoadClaimSet(unitID, versionID string) (domain.ClaimSet, bool, error) {
	p := r.claimsetPath(unitID, versionID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ClaimSet{}, false, nil
	}
	if err != nil {
		return domain.ClaimSet{}, false, err
	}
	var cs domain.ClaimSet
	if err := json.Unmarshal(b, &cs); err != nil {
		return domain.ClaimSet{}, false, err
	}
	return cs, true, nil
}

// SaveUncertainty stores a canonicalized uncertainty sidecar for the unit atomically
// and updates the embedded version record's uncertainty_hash. Returns ErrUnitNotFound
// or ErrVersionNotFound where applicable.
func (r *UnitRepo) SaveUncertainty(unitID, versionID string, u domain.Uncertainty, uncertaintyHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.unitPath(unitID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.ErrUnitNotFound
	}
	if err != nil {
		return err
	}
	var ur UnitRecord
	if err := json.Unmarshal(b, &ur); err != nil {
		return err
	}

	// find version
	found := false
	for i := range ur.Versions {
		if ur.Versions[i].ID == versionID {
			ur.Versions[i].UncertaintyHash = uncertaintyHash
			found = true
			break
		}
	}
	if !found {
		return domain.ErrVersionNotFound
	}

	sidecar := r.uncertaintyPath(unitID, versionID)
	tmp := r.uncertaintyTempPath(unitID, versionID)
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, sidecar); err != nil {
		return err
	}

	// persist updated unit record
	out, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}
	tmpUnit := r.unitTempPath(ur.ID)
	if err := ioutil.WriteFile(tmpUnit, out, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpUnit, p); err != nil {
		return err
	}
	return nil
}

// LoadUncertainty loads a version-scoped uncertainty sidecar if present.
func (r *UnitRepo) LoadUncertainty(unitID, versionID string) (domain.Uncertainty, bool, error) {
	p := r.uncertaintyPath(unitID, versionID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.Uncertainty{}, false, nil
	}
	if err != nil {
		return domain.Uncertainty{}, false, err
	}
	var u domain.Uncertainty
	if err := json.Unmarshal(b, &u); err != nil {
		return domain.Uncertainty{}, false, err
	}
	return u, true, nil
}

// LoadMeaning loads a version-scoped meaning sidecar if present.
func (r *UnitRepo) LoadMeaning(unitID, versionID string) (domain.Meaning, bool, error) {
	p := r.meaningPath(unitID, versionID)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return domain.Meaning{}, false, nil
	}
	if err != nil {
		return domain.Meaning{}, false, err
	}
	var m domain.Meaning
	if err := json.Unmarshal(b, &m); err != nil {
		return domain.Meaning{}, false, err
	}
	return m, true, nil
}
