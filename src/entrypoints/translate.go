package entrypoints

import (
	"bytes"
	"crowdin-grazie/grazie"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type request struct {
	Strings []string
}

type response struct {
	Data responseData `json:"data"`
}

type responseData struct {
	Translations []string `json:"translations"`
}

func TranslateHandler(grazieInstance grazie.Grazie) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		respBody, _ := ioutil.ReadAll(r.Body)

		var requestBody = request{}
		err = json.NewDecoder(bytes.NewReader(respBody)).Decode(&requestBody)
		if err != nil {
			httpErrorAndLog(w, fmt.Errorf("error parsing request body: %v", err), http.StatusBadRequest)
			return
		}

		var target = r.URL.Query().Get("target")
		translateResponse, err := grazieInstance.Translate(grazie.TranslateRequest{Texts: requestBody.Strings, ToLang: target})
		if err != nil {
			httpErrorAndLog(w, fmt.Errorf("error translating: %v", err), http.StatusInternalServerError)
			return
		}

		marshalledResponse, err := json.Marshal(response{
			Data: responseData{
				Translations: translateResponse.Translations,
			},
		})
		if err != nil {
			httpErrorAndLog(w, fmt.Errorf("error marshalling response: %v", err), http.StatusInternalServerError)
			return
		}
		httpSuccess(w, marshalledResponse)
	}
}
