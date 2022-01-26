package entrypoints

import (
	"bytes"
	"crowdin-grazie/grazie"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
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

func TranslateHandler(grazieInstance grazie.Grazie, clientSecret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		token := r.URL.Query().Get("jwt[token]")
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(clientSecret), nil
		})

		if err != nil {
			httpErrorAndLog(w, fmt.Errorf("failed to parse JWT: %v", err), http.StatusBadRequest)
			return
		}
		if !parsedToken.Valid {
			httpErrorAndLog(w, fmt.Errorf("invalid JWT token"), http.StatusBadRequest)
			return
		}

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
