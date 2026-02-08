package memory

import "digiemu-core/internal/kernel/domain"

func (r *UnitRepo) FindVersionByID(versionID string) (domain.Version, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.versionsByID[versionID]
	return v, ok, nil
}
