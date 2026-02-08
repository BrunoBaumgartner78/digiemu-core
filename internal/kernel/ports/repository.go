package ports

import "digiemu-core/internal/kernel/domain"

type UnitRepository interface {
	ExistsByKey(key string) (bool, error)
	SaveUnit(u domain.Unit) error
	FindUnitByKey(key string) (domain.Unit, bool, error)
	FindUnitByID(id string) (domain.Unit, bool, error)

	// v0.2.1: needed for audit verification / exports
	ListUnits() ([]domain.Unit, error)

	SaveVersion(v domain.Version) error
	// ListVersionsByUnitID returns versions in creation order (oldest -> newest).
	// Adapters MUST provide stable ordering to keep history operations deterministic.
	ListVersionsByUnitID(unitID string) ([]domain.Version, error)

	// v0.2: head tracking for optimistic locking
	UpdateUnitHead(unitID, headVersionID string) error

	// FindVersionByID returns a version by its ID.
	// ok=false if not found.
	FindVersionByID(versionID string) (domain.Version, bool, error)

	// SaveMeaning persists a meaning sidecar for a specific version and updates
	// the version record's meaning_hash. Implementations MUST perform only
	// persistence operations and MUST NOT emit audit events.
	SaveMeaning(unitID, versionID string, meaning domain.Meaning, meaningHash string) error
	LoadMeaning(unitID, versionID string) (domain.Meaning, bool, error)
}
