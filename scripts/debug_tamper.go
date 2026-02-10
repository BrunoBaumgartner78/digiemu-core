package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func main() {
	dir, err := ioutil.TempDir("", "digiemu-debug-")
	if err != nil {
		panic(err)
	}
	fmt.Println("tempdir:", dir)
	// do not remove dir so we can inspect

	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := memory.RealClock{}

	cu := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	var outU ports.CreateUnitResponse
	outU, err = cu.CreateUnit(ports.CreateUnitRequest{Key: "kfs", Title: "tfs", Description: "d", ActorID: "debug"})
	if err != nil {
		panic(err)
	}
	cv := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	var outV ports.CreateVersionResponse
	outV, err = cv.CreateVersion(ports.CreateVersionRequest{UnitKey: "kfs", Label: "lbl", Content: "c", ActorID: "debug"})
	if err != nil {
		panic(err)
	}

	su := usecases.SetUncertainty{Repo: repo, Audit: audit, Clock: clock}
	uJSON := []byte(`{"schema_version":"uncertainty/v0","id":"u1","type":"empirical","level":"low","applies_to":{"scope":"version"}}`)
	var resp ports.SetUncertaintyResponse
	resp, err = su.SetUncertainty(ports.SetUncertaintyRequest{UnitKey: "kfs", VersionID: outV.VersionID, BodyBytes: uJSON, ActorID: "debug"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("SetUncertainty stored: %+v\n", resp)

	auditPath := filepath.Join(dir, "audit.ndjson")
	fmt.Println("Audit file path:", auditPath)
	if b, err := ioutil.ReadFile(auditPath); err == nil {
		fmt.Println("Audit contents:")
		fmt.Println(string(b))
	} else {
		fmt.Println("read audit err:", err)
	}

	side := filepath.Join(dir, "units", outU.UnitID+"."+outV.VersionID+".uncertainty.json")
	if b, err := ioutil.ReadFile(side); err == nil {
		fmt.Println("Sidecar before tamper:")
		var pretty any
		_ = json.Unmarshal(b, &pretty)
		pb, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Println(string(pb))
	}

	// tamper
	if err := ioutil.WriteFile(side, []byte("{}"), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("Tampered sidecar written")
	if b, err := ioutil.ReadFile(side); err == nil {
		fmt.Println("Sidecar after tamper:")
		fmt.Println(string(b))
	}

	// direct repo check
	u, ok, err := repo.LoadUncertainty(outU.UnitID, outV.VersionID)
	fmt.Println("LoadUncertainty ok:", ok, "err:", err)
	fmt.Printf("Loaded struct: %+v\n", u)
	// compute hash from loaded struct
	h, herr := usecases.ComputeUncertaintyHashFromStruct(u)
	fmt.Println("Computed hash from loaded sidecar:", h, "err:", herr)

	// run verify audit
	ver := usecases.VerifyAudit{Repo: repo, Audit: fsrepo.NewAuditReader(dir)}
	var res ports.VerifyAuditResponse
	res, err = ver.VerifyAudit(ports.VerifyAuditRequest{UnitKey: "kfs", StrictHash: true})
	if err != nil {
		fmt.Println("VerifyAudit error:", err)
		os.Exit(1)
	}
	fmt.Printf("Verify result: %+v\n", res)
}
