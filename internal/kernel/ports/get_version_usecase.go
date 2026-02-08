package ports

type GetVersionRequest struct {
	VersionID string
}

type GetVersionResponse struct {
	Version VersionDTO
}

type GetVersionUsecase interface {
	GetVersion(in GetVersionRequest) (GetVersionResponse, error)
}
