package ports

// DTOs for kernel usecases - keep them primitive-friendly and independent
// from domain or adapter types.

type CreateUnitRequest struct {
	Key         string
	Title       string
	Description string
	ActorID     string // v0.2: strict audit
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
	// v0.2
	BaseVersionID string // optional optimistic locking; "" = no check
	ActorID       string // strict audit
}

type CreateVersionResponse struct {
	VersionID string
	UnitID    string
	Label     string
	Content   string
}

type SetMeaningRequest struct {
	UnitKey     string
	VersionID   string // optional; empty means use head
	MeaningJSON []byte
	ActorID     string
}

type SetMeaningResponse struct {
	UnitID      string
	VersionID   string
	MeaningHash string
}
