package httpapi

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	j "digiemu-core/internal/httpapi/json"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type API struct {
	Units   ports.CreateUnitUsecase
	Vers    ports.CreateVersionUsecase
	Meaning ports.SetMeaningUsecase
	Repo    ports.UnitRepository
}

type createUnitReq struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
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
		j.Errorf(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid json: %v", err)
		return
	}
	// Support multiple client shapes for backward compatibility:
	// - Preferred: {"title":"...","description":"...","key":"..."}
	// - Legacy tests/clients: {"name":"desired-key","description":"title"}
	var titleVal string
	var keyVal string
	var descVal string
	if strings.TrimSpace(req.Title) != "" {
		titleVal = req.Title
		keyVal = req.Key
		descVal = req.Description
	} else if strings.TrimSpace(req.Name) != "" && strings.TrimSpace(req.Description) != "" {
		// legacy behavior: name==key, description==title
		keyVal = req.Name
		titleVal = req.Description
		descVal = req.Description
	} else if strings.TrimSpace(req.Name) != "" {
		// treat name as title if no description provided
		titleVal = req.Name
		keyVal = req.Key
		descVal = req.Description
	}

	if strings.TrimSpace(titleVal) == "" {
		j.Errorf(w, http.StatusBadRequest, "VALIDATION_ERROR", "title required")
		return
	}

	in := ports.CreateUnitRequest{Key: keyVal, Title: titleVal, Description: descVal}
	out, err := a.Units.CreateUnit(in)
	if err != nil {
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	_ = j.Write(w, http.StatusCreated, struct {
		UnitID      string `json:"unitId"`
		CreatedAt   string `json:"createdAt"`
		Key         string `json:"key"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
	}{UnitID: out.UnitID, CreatedAt: time.Now().UTC().Format(time.RFC3339), Key: out.Key, Title: out.Title, Description: out.Description})
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
		j.Errorf(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid json: %v", err)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		j.Errorf(w, http.StatusBadRequest, "VALIDATION_ERROR", "content required")
		return
	}

	// label: simple timestamp
	label := time.Now().UTC().Format("20060102T150405Z")
	in := ports.CreateVersionRequest{UnitKey: unitKey, Label: label, Content: req.Content}
	out, err := a.Vers.CreateVersion(in)
	if err != nil {
		// if unit not found, map to 404
		if err == domain.ErrUnitNotFound {
			j.ErrorCode(w, http.StatusNotFound, "UNIT_NOT_FOUND", "unit not found", nil)
			return
		}
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	_ = j.Write(w, http.StatusCreated, createVersionRes{VersionID: out.VersionID, CreatedAt: time.Now().UTC().Format(time.RFC3339)})
}

func (a API) handleHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func (a API) handleSetMeaning(w http.ResponseWriter, r *http.Request, unitKey string) {
	// optional version query param
	version := r.URL.Query().Get("version")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		j.Errorf(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid body: %v", err)
		return
	}
	in := ports.SetMeaningRequest{UnitKey: unitKey, VersionID: version, MeaningJSON: body, ActorID: "http"}
	out, err := a.Meaning.SetMeaning(in)
	if err != nil {
		if err == domain.ErrUnitNotFound {
			j.ErrorCode(w, http.StatusNotFound, "UNIT_NOT_FOUND", "unit not found", nil)
			return
		}
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	_ = j.Write(w, http.StatusCreated, struct {
		UnitID      string `json:"unit_id"`
		VersionID   string `json:"version_id"`
		MeaningHash string `json:"meaning_hash"`
	}{UnitID: out.UnitID, VersionID: out.VersionID, MeaningHash: out.MeaningHash})
}

func (a API) handleGetMeaning(w http.ResponseWriter, r *http.Request, unitKey string) {
	version := r.URL.Query().Get("version")
	// resolve unit by key or id
	u, ok, err := a.Repo.FindUnitByKey(unitKey)
	if err != nil {
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	if !ok {
		u2, ok2, err2 := a.Repo.FindUnitByID(unitKey)
		if err2 != nil {
			j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err2)
			return
		}
		if !ok2 {
			j.ErrorCode(w, http.StatusNotFound, "UNIT_NOT_FOUND", "unit not found", nil)
			return
		}
		u = u2
	}
	if version == "" {
		version = u.HeadVersionID
	}
	if version == "" {
		j.ErrorCode(w, http.StatusNotFound, "VERSION_NOT_FOUND", "no version specified and unit has no head", nil)
		return
	}
	v, found, err := a.Repo.FindVersionByID(version)
	if err != nil {
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	if !found {
		j.ErrorCode(w, http.StatusNotFound, "VERSION_NOT_FOUND", "version not found", nil)
		return
	}
	m, ok, err := a.Repo.LoadMeaning(u.ID, version)
	if err != nil {
		j.Errorf(w, http.StatusInternalServerError, "INTERNAL", "%v", err)
		return
	}
	if !ok {
		j.ErrorCode(w, http.StatusNotFound, "MEANING_NOT_FOUND", "meaning not found", nil)
		return
	}
	_ = j.Write(w, http.StatusOK, struct {
		Meaning     any    `json:"meaning"`
		MeaningHash string `json:"meaning_hash"`
	}{Meaning: m, MeaningHash: v.MeaningHash})
}
