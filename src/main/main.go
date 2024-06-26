package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"crowdin-grazie/config"
	"crowdin-grazie/entrypoints"
	"crowdin-grazie/slack"

	graziego "git.jetbrains.team/mau/grazie-ml-go-client.git"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		logrus.WithError(err).Fatal("cannot get config")
	}

	grazieMlClient, err := graziego.New(
		cfg.GrazieHost,
		graziego.AuthTypeApplication,
		cfg.GrazieToken,
		graziego.GrazieAgent{
			Name:    "mau-crowdin-mt",
			Version: cfg.GrazieVersion,
		},
	)
	if err != nil {
		logrus.WithError(err).Fatal("cannot create Grazie client")
	}

	slackClient := slack.New(cfg)
	hc := entrypoints.NewHandlerCreator(slackClient)

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", func(_ http.ResponseWriter, _ *http.Request) {}).Methods(http.MethodGet)
	r.HandleFunc("/manifest.json", hc.ManifestHandler(cfg.ClientID)).Methods(http.MethodGet)
	r.HandleFunc("/installed", hc.InstalledHandler).Methods(http.MethodPost)
	r.HandleFunc("/translate/", hc.TranslateHandler(grazieMlClient, cfg.ClientSecret)).Methods(http.MethodPost)
	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("static")))

	logrus.Info("service starting..")
	log.Fatal(http.ListenAndServe(":8080", logRequest(r)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
