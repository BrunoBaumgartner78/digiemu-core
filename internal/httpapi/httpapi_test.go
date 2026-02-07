package httpapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	usecases "digiemu-core/internal/kernel/usecases"
)

func TestAPI_CreateUnit_And_Version(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	api := API{Units: usecases.CreateUnit{Repo: repo}, Vers: usecases.CreateVersion{Repo: repo}}
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
	api := API{Units: usecases.CreateUnit{Repo: repo}, Vers: usecases.CreateVersion{Repo: repo}}
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
