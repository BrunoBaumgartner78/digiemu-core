package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	mem "digiemu-core/internal/kernel/adapters/memory"
	usecases "digiemu-core/internal/kernel/usecases"
)

func TestAPI_CreateUnit_And_Version(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := mem.RealClock{}
	api := API{Units: usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}, Vers: usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}}
	srv := httptest.NewServer(NewRouter(api))
	defer srv.Close()

	// create unit
	reqBody := map[string]string{"name": "my-unit", "description": "desc"}
	b, _ := json.Marshal(reqBody)
	res, err := http.Post(srv.URL+"/v1/units", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post unit err: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", res.StatusCode)
	}
	var cru struct {
		UnitID string `json:"unitId"`
		Key    string `json:"key"`
	}
	_ = json.NewDecoder(res.Body).Decode(&cru)
	res.Body.Close()

	// post version
	vreq := map[string]string{"content": "hello"}
	vb, _ := json.Marshal(vreq)
	res, err = http.Post(srv.URL+"/v1/units/"+cru.Key+"/versions", "application/json", bytes.NewReader(vb))
	if err != nil {
		t.Fatalf("post version err: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status version: %d", res.StatusCode)
	}
	var vr struct {
		VersionID string `json:"versionId"`
	}
	_ = json.NewDecoder(res.Body).Decode(&vr)
	res.Body.Close()

	// validate file exists and contains version
	p := filepath.Join(dir, "units", cru.UnitID+".json")
	dat, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("read file err: %v", err)
	}
	if !bytes.Contains(dat, []byte(vr.VersionID)) {
		t.Fatalf("version id not in file")
	}
}

func TestAPI_Version_UnknownUnit_404(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := mem.RealClock{}
	api := API{Units: usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}, Vers: usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}}
	srv := httptest.NewServer(NewRouter(api))
	defer srv.Close()

	vreq := map[string]string{"content": "hello"}
	vb, _ := json.Marshal(vreq)
	res, err := http.Post(srv.URL+"/v1/units/does-not-exist/versions", "application/json", bytes.NewReader(vb))
	if err != nil {
		t.Fatalf("post version err: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestAPI_ValidationAndErrors(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := mem.RealClock{}
	api := API{Units: usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}, Vers: usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}}
	srv := httptest.NewServer(NewRouter(api))
	defer srv.Close()

	// 1) POST /v1/units with empty title => 400 VALIDATION_ERROR
	b, _ := json.Marshal(map[string]string{"name": ""})
	res, err := http.Post(srv.URL+"/v1/units", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post unit err: %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	var errBody map[string]map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&errBody)
	res.Body.Close()
	if errBody["error"]["code"] != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %v", errBody["error"]["code"])
	}

	// 2) POST /v1/units valid => 201 + includes key
	b, _ = json.Marshal(map[string]string{"name": "ok-unit", "description": "desc"})
	res, err = http.Post(srv.URL+"/v1/units", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post unit err: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
	var cru struct {
		Key string `json:"key"`
	}
	_ = json.NewDecoder(res.Body).Decode(&cru)
	res.Body.Close()
	if cru.Key == "" {
		t.Fatalf("expected key in response")
	}

	// 3) POST /v1/units/{key}/versions with empty content => 400 VALIDATION_ERROR
	b, _ = json.Marshal(map[string]string{"content": ""})
	res, err = http.Post(srv.URL+"/v1/units/"+cru.Key+"/versions", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post version err: %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	var eb map[string]map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&eb)
	res.Body.Close()
	if eb["error"]["code"] != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %v", eb["error"]["code"])
	}

	// 4) POST /v1/units/{missing}/versions => 404 UNIT_NOT_FOUND
	b, _ = json.Marshal(map[string]string{"content": "x"})
	res, err = http.Post(srv.URL+"/v1/units/does-not-exist/versions", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post version err: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
	var eb2 map[string]map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&eb2)
	res.Body.Close()
	if eb2["error"]["code"] != "UNIT_NOT_FOUND" {
		t.Fatalf("expected UNIT_NOT_FOUND, got %v", eb2["error"]["code"])
	}
	_ = io.EOF
}
