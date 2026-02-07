package ports

// Usecase interfaces - explicit contract for kernel operations.

type CreateUnitUsecase interface {
	CreateUnit(req CreateUnitRequest) (CreateUnitResponse, error)
}

type CreateVersionUsecase interface {
	CreateVersion(req CreateVersionRequest) (CreateVersionResponse, error)
}
