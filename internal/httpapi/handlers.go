package httpapi

import (
	"net/http"
	"time"

	j "digiemu-core/internal/httpapi/json"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type API struct {
	Units ports.CreateUnitUsecase
	Vers  ports.CreateVersionUsecase
}

type createUnitReq struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type createUnitRes struct {
	UnitID    string `json:"unitId"`
	CreatedAt string `json:"createdAt"`
	Key       string `json:"key"`
}

func (a API) handleCreateUnit(w http.ResponseWriter, r *http.Request) {
	var req createUnitReq
	if err := j.Read(r, &req); err != nil {
		j.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		j.Error(w, http.StatusBadRequest, "name required")
		return
	}

	in := ports.CreateUnitRequest{Key: req.Name, Title: req.Description}
	out, err := a.Units.CreateUnit(in)
	if err != nil {
		j.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = j.Write(w, http.StatusCreated, createUnitRes{UnitID: out.UnitID, CreatedAt: time.Now().UTC().Format(time.RFC3339), Key: out.Key})
}

type createVersionReq struct {
	Content string `json:"content"`
	Note    string `json:"note,omitempty"`
}

type createVersionRes struct {
	VersionID string `json:"versionId"`
	CreatedAt string `json:"createdAt"`
}

func (a API) handleCreateVersion(w http.ResponseWriter, r *http.Request, unitKey string) {
	var req createVersionReq
	if err := j.Read(r, &req); err != nil {
		j.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Content == "" {
		j.Error(w, http.StatusBadRequest, "content required")
		return
	}

	// label: simple timestamp
	label := time.Now().UTC().Format("20060102T150405Z")
	in := ports.CreateVersionRequest{UnitKey: unitKey, Label: label, Content: req.Content}
	out, err := a.Vers.CreateVersion(in)
	if err != nil {
		// if unit not found, map to 404
		if err == domain.ErrUnitNotFound {
			j.Error(w, http.StatusNotFound, "unit not found")
			return
		}
		j.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = j.Write(w, http.StatusCreated, createVersionRes{VersionID: out.VersionID, CreatedAt: time.Now().UTC().Format(time.RFC3339)})
}

func (a API) handleHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
