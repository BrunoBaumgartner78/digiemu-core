package fs

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"digiemu-core/internal/kernel/domain"
)

// FindVersionByID scans unit records to locate a version by ID.
// v0.2.3: simple implementation; can be optimized later by indexing or splitting storage.
func (r *UnitRepo) FindVersionByID(versionID string) (domain.Version, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	files, err := ioutil.ReadDir(r.unitsDir)
	if err != nil {
		return domain.Version{}, false, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		p := filepath.Join(r.unitsDir, f.Name())
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return domain.Version{}, false, err
		}
		var ur UnitRecord
		if err := json.Unmarshal(b, &ur); err != nil {
			return domain.Version{}, false, err
		}
		for _, vr := range ur.Versions {
			if vr.ID == versionID {
				return domain.Version{
					ID:            vr.ID,
					UnitID:        ur.ID,
					Label:         vr.Label,
					Content:       vr.Content,
					PrevVersionID: vr.PrevVersionID,
					ContentHash:   vr.ContentHash,
					CreatedAtUnix: 0,
					ActorID:       vr.ActorID,
				}, true, nil
			}
		}
	}
	return domain.Version{}, false, nil
}
