package entrypoints

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Manifest struct {
	Identifier     string                         `json:"identifier"`
	Name           string                         `json:"name"`
	Logo           string                         `json:"logo"`
	BaseUrl        string                         `json:"baseUrl"`
	Authentication Authentication                 `json:"authentication"`
	Events         map[string]string              `json:"events"`
	Scopes         []string                       `json:"scopes"`
	Modules        map[string][]ModuleDeclaration `json:"modules"`
}

type Authentication struct {
	Type     string `json:"type"`
	ClientId string `json:"clientId"`
}

type ModuleDeclaration struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Url         string `json:"url"`
}

func ManifestHandler(clientId string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var manifest = Manifest{
			Identifier: "grazie-mt",
			Name:       "Grazie MT",
			Logo:       "/assets/logo.svg",
			BaseUrl:    "https://" + r.Host,
			Authentication: Authentication{
				Type:     "authorization_code",
				ClientId: clientId,
			},
			Events: map[string]string{
				"installed": "/installed",
			},
			Scopes: []string{"project"},
			Modules: map[string][]ModuleDeclaration{
				"custom-mt": {
					{
						Key:         "grazie-mt",
						Name:        "Grazie MT",
						Icon:        "/assets/logo.svg",
						Description: "",
						Logo:        "/assets/logo.svg",
						Url:         "/translate",
					},
				},
			},
		}
		s, err := json.Marshal(manifest)
		if err != nil {
			httpErrorAndLog(w, fmt.Errorf("error marshalling manifest: %v", err), http.StatusInternalServerError)
			return
		}

		httpSuccess(w, s)
	}
}
