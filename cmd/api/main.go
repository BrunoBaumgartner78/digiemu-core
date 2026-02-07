package main

import (
	"flag"
	"log"
	"net/http"

	"digiemu-core/internal/httpapi"
	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/usecases"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	data := flag.String("data", "./data", "data directory")
	flag.Parse()

	repo := fsrepo.NewUnitRepo(*data)
	createUnit := usecases.CreateUnit{Repo: repo}
	createVersion := usecases.CreateVersion{Repo: repo}

	r := httpapi.NewRouter(httpapi.API{Units: createUnit, Vers: createVersion})

	srv := &http.Server{Addr: *addr, Handler: r}
	log.Printf("api listening on %s (data=%s)", *addr, *data)
	log.Fatal(srv.ListenAndServe())
}
