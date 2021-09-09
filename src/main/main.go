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
	var appUrl = environment.MustGetEnv(environment.AppUrl)

	r := mux.NewRouter()
	r.HandleFunc("/manifest.json", entrypoints.ManifestHandler(appUrl)).Methods(http.MethodGet)
	r.HandleFunc("/installed", entrypoints.InstalledHandler).Methods(http.MethodPost)
	r.HandleFunc("/translate/", entrypoints.TranslateHandler(grazie)).Methods(http.MethodPost)
	r.PathPrefix("/logo.svg").Handler(http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":8080", logRequest(r)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}