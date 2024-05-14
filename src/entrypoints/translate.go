package entrypoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	graziego "git.jetbrains.team/mau/grazie-ml-go-client.git"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var (
	supportedLanguages = map[string]bool{
		graziego.LangEN: true,
		graziego.LangDE: true,
		graziego.LangFR: true,
		graziego.LangES: true,
		graziego.LangRU: true,
		graziego.LangKO: true,
		graziego.LangZH: true,
		graziego.LangJA: true,
	}
)

type request struct {
	Strings []string `json:"strings"`
}

type response struct {
	Data responseData `json:"data"`
}

type responseData struct {
	Translations []string `json:"translations"`
}

func (hc *HandlerCreator) TranslateHandler(grazieMlClient graziego.Client, clientSecret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		query := r.URL.Query()
		logEntry := logrus.WithField("query", query)

		token := r.URL.Query().Get("jwtToken")
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(clientSecret), nil
		})

		if err != nil {
			logEntry.WithError(err).Error("failed to parse JWT")
			hc.httpErrorAndLog(w, fmt.Errorf("failed to parse JWT: %v", err), http.StatusBadRequest)
			return
		}
		if !parsedToken.Valid {
			logEntry.WithField("token", parsedToken).Error("invalid JWT token")
			hc.httpErrorAndLog(w, fmt.Errorf("invalid JWT token"), http.StatusBadRequest)
			return
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logEntry.Error("error reading body")
			hc.httpErrorAndLog(w, fmt.Errorf("error reading body: %v", err), http.StatusBadRequest)
			return
		}

		logEntry = logEntry.WithField("request_body", string(reqBody))

		var requestBody = request{}
		err = json.NewDecoder(bytes.NewReader(reqBody)).Decode(&requestBody)
		if err != nil {
			logEntry.WithField("body", string(reqBody)).Error("error parsing request body")
			hc.httpErrorAndLog(w, fmt.Errorf("error parsing request body: %v", err), http.StatusBadRequest)
			return
		}

		logEntry = logEntry.WithField("request", requestBody)

		var translations []string
		target := getLang(query.Get("target"))
		source := getLang(query.Get("source"))
		if !supportedLanguages[target] || !supportedLanguages[source] {
			translations = requestBody.Strings
		} else {
			translateResponse, err := grazieMlClient.TranslateWithoutAI(r.Context(), source, target, requestBody.Strings)
			if err != nil || translateResponse == nil {
				logEntry.WithError(err).Error("error translating")
				hc.httpErrorAndLog(w, fmt.Errorf("error translating: %v", err), http.StatusInternalServerError)
				return
			}

			translations = make([]string, 0, len(translateResponse.Translations))
			for _, t := range translateResponse.Translations {
				translations = append(translations, t.Translation)
			}
		}

		resp := response{
			Data: responseData{
				Translations: translations,
			},
		}
		marshalledResponse, err := json.Marshal(resp)
		if err != nil {
			logEntry.WithField("response", resp).WithError(err).Error("error marshalling response")
			hc.httpErrorAndLog(w, fmt.Errorf("error marshalling response: %v", err), http.StatusInternalServerError)
			return
		}
		hc.httpSuccess(w, marshalledResponse)
	}
}

func getLang(lang string) string {
	return strings.ToLower(strings.Split(lang, "-")[0])
}
