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
}

func NewUnitRepo(basePath string) *UnitRepo {
	units := filepath.Join(basePath, "units")
	_ = os.MkdirAll(units, 0o755)
	return &UnitRepo{basePath: basePath, unitsDir: units}
}

func (r *UnitRepo) unitPath(id string) string {
	return filepath.Join(r.unitsDir, id+".json")
}

func (r *UnitRepo) unitTempPath(id string) string {
	return filepath.Join(r.unitsDir, id+".json.tmp")
}

func (r *UnitRepo) ExistsByKey(key string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
		ID:          u.ID,
		Key:         u.Key,
		Title:       u.Title,
		Description: u.Description,
		CreatedAt:   nowRFC3339(),
		Versions:    []VersionRecord{},
	}

	data, err := json.MarshalIndent(ur, "", "  ")
	if err != nil {
		return err
	}

	tmp := r.unitTempPath(u.ID)
	if err := ioutil.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.unitPath(u.ID))
}

func (r *UnitRepo) FindUnitByKey(key string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
			return domain.Unit{ID: ur.ID, Key: ur.Key, Title: ur.Title, Description: ur.Description}, true, nil
		}
	}
	return domain.Unit{}, false, nil
}

func (r *UnitRepo) FindUnitByID(id string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
	return domain.Unit{ID: ur.ID, Key: ur.Key, Title: ur.Title, Description: ur.Description}, true, nil
}

func (r *UnitRepo) SaveVersion(v domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// load unit record
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
		ID:        v.ID,
		Label:     v.Label,
		Content:   v.Content,
		CreatedAt: nowRFC3339(),
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
	return os.Rename(tmp, p)
}

func (r *UnitRepo) ListVersionsByUnitID(unitID string) ([]domain.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
		out = append(out, domain.Version{ID: vr.ID, UnitID: unitID, Label: vr.Label, Content: vr.Content})
	}
	return out, nil
}

// helper: ensure repo implements interface at compile time
var _ = func() interface{} {
	var _r interface{} = (*UnitRepo)(nil)
	_ = _r
	return nil
}()
