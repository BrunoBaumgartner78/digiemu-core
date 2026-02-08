package ports

type GetUnitRequest struct {
	UnitKey string
}

type UnitDTO struct {
	ID            string
	Key           string
	Title         string
	Description   string
	HeadVersionID string
}

type GetUnitResponse struct {
	Unit UnitDTO
}

type GetUnitUsecase interface {
	GetUnit(in GetUnitRequest) (GetUnitResponse, error)
}

type ListUnitsRequest struct {
	// Optional: filter by prefix on key (simple and stable).
	KeyPrefix string
}

type ListUnitsResponse struct {
	Units []UnitDTO
}

type ListUnitsUsecase interface {
	ListUnits(in ListUnitsRequest) (ListUnitsResponse, error)
}

type ListVersionsRequest struct {
	UnitKey string
	// If true, newest first (default false: oldest first).
	NewestFirst bool
}

type VersionDTO struct {
	ID            string
	UnitID        string
	Label         string
	Content       string
	PrevVersionID string
	ContentHash   string
	CreatedAtUnix int64
	ActorID       string
}

type ListVersionsResponse struct {
	UnitID   string
	Versions []VersionDTO
}

type ListVersionsUsecase interface {
	ListVersions(in ListVersionsRequest) (ListVersionsResponse, error)
}

type GetHeadVersionRequest struct {
	UnitKey string
}

type GetHeadVersionResponse struct {
	UnitID  string
	Version VersionDTO
}

type GetHeadVersionUsecase interface {
	GetHeadVersion(in GetHeadVersionRequest) (GetHeadVersionResponse, error)
}
