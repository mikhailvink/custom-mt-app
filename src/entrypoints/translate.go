package entrypoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	graziego "git.jetbrains.team/mau/grazie-ml-go-client.git"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
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

		// ignore requests from the "regular" Crowdin instance
		projectID := query.Get("project_id")
		if len(projectID) == 6 {
			http.Error(w, "requests only from jetbrains.crowdin.com are supported", http.StatusBadRequest)
			return
		}

		token := query.Get("jwtToken")
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
			//hc.httpErrorAndLog(w, fmt.Errorf("failed to parse JWT: %v", err), http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		//target := query.Get("target")

		//results := make(map[int]string, len(requestBody.Strings))
		//var reqError error
		//wg := sync.WaitGroup{}
		//m := sync.Mutex{}
		//
		//for i := 0; i < len(requestBody.Strings); i += 1 {
		//	wg.Add(1)
		//	go func(number int, stringToTranslate string) {
		//		defer wg.Done()
		//
		//		translateResponse, err := grazieMlClient.Translate(
		//			r.Context(), graziego.CrowdinTranslateTag, target, stringToTranslate,
		//		)
		//		if err != nil {
		//			logEntry.WithError(err).Error("error translating")
		//
		//			m.Lock()
		//			defer m.Unlock()
		//			reqError = fmt.Errorf("error translating: %v", err)
		//			return
		//		}
		//
		//		m.Lock()
		//		defer m.Unlock()
		//		results[number] = translateResponse
		//	}(i, requestBody.Strings[i])
		//}
		//
		//wg.Wait()
		//
		//if reqError != nil {
		//	hc.httpErrorAndLog(w, reqError, http.StatusInternalServerError)
		//	return
		//}
		//
		//translations := make([]string, 0, len(requestBody.Strings))
		//for i := 0; i < len(requestBody.Strings); i += 1 {
		//	translations = append(translations, results[i])
		//}

		// TODO" remove mock after setting translation workflow
		translations := make([]string, len(requestBody.Strings))

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
