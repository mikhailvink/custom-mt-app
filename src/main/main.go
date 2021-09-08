package main

import (
	"crowdin-grazie/entrypoints"
	"crowdin-grazie/environment"
	grazie2 "crowdin-grazie/grazie"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var grazie = grazie2.New(environment.MustGetEnv(environment.EnvGrazieToken))

	r := mux.NewRouter()
	r.HandleFunc("/manifest.json", entrypoints.ManifestHandler).Methods(http.MethodGet)
	r.HandleFunc("/installed", entrypoints.InstalledHandler).Methods(http.MethodPost)
	r.HandleFunc("/translate", entrypoints.TranslateHandler(grazie)).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
