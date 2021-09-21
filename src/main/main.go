package main

import (
	"crowdin-grazie/entrypoints"
	"crowdin-grazie/environment"
	"crowdin-grazie/grazie"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	var grazieInstance = grazie.New(environment.MustGetEnv(environment.EnvGrazieToken), grazie.Config{
		Host: os.Getenv(environment.GrazieHost),
	})
	var clientId = environment.MustGetEnv(environment.ClientId)
	var clientSecret = environment.MustGetEnv(environment.ClientSecret)

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", func(_ http.ResponseWriter, _ *http.Request) {}).Methods(http.MethodGet)
	r.HandleFunc("/manifest.json", entrypoints.ManifestHandler(clientId)).Methods(http.MethodGet)
	r.HandleFunc("/installed", entrypoints.InstalledHandler).Methods(http.MethodPost)
	r.HandleFunc("/translate/", entrypoints.TranslateHandler(grazieInstance, clientSecret)).Methods(http.MethodPost)
	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":8080", logRequest(r)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
