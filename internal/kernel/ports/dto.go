package ports

// DTOs for kernel usecases - keep them primitive-friendly and independent
// from domain or adapter types.

type CreateUnitRequest struct {
	Key         string
	Title       string
	Description string
}

type CreateUnitResponse struct {
	UnitID      string
	Key         string
	Title       string
	Description string
}

type CreateVersionRequest struct {
	UnitKey string
	Label   string
	Content string
}

type CreateVersionResponse struct {
	VersionID string
	UnitID    string
	Label     string
	Content   string
}
