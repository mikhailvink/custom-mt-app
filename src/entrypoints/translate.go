package entrypoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	graziego "git.jetbrains.team/mau/grazie-ml-go-client.git"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

const (
	batchSize = 10
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

		target := query.Get("target")

		results := make(map[int64][]string, len(requestBody.Strings)/batchSize)
		wg := sync.WaitGroup{}
		m := sync.Mutex{}
		number := int64(0)

		for i := 0; i < len(requestBody.Strings); i += batchSize {
			end := i + batchSize
			if end > len(requestBody.Strings) {
				end = len(requestBody.Strings)
			}

			batch := requestBody.Strings[i:end]

			wg.Add(1)
			go func(number int64, batch []string) {
				defer wg.Done()

				data, err := json.Marshal(batch)
				if err != nil {
					logEntry.WithField("strings", batch).Error("failed to marshal strings")
					hc.httpErrorAndLog(w, fmt.Errorf("failed to marshal strings: %v", err), http.StatusInternalServerError)
					return
				}

				prompt := fmt.Sprintf(
					"Translate JSON array of strings to the target language. "+
						"Answer with only JSON array of translated strings in the same order. Do not use any wrapping for response.\n"+
						"Target language code: %s\n"+
						"Array:\n%s",
					target, string(data),
				)

				chatgptResp, err := grazieMlClient.Chat(r.Context(), "gpt-4-1106-preview", []graziego.ChatMessage{
					{
						Role: graziego.RoleUser,
						Text: prompt,
					},
				})
				if err != nil {
					logEntry.WithError(err).Error("error translating")
					hc.httpErrorAndLog(w, fmt.Errorf("error translating: %v", err), http.StatusInternalServerError)
					return
				}

				answer := chatgptResp.Text
				logEntry = logEntry.WithField("answer", answer)

				batchTranslations := make([]string, 0)
				err = json.Unmarshal([]byte(answer), &batchTranslations)
				if err != nil {
					logEntry.WithError(err).Error("failed to unmarshal response")
					hc.httpErrorAndLog(w, fmt.Errorf("failed to unmarshal response: %v", err), http.StatusInternalServerError)
					return
				}

				m.Lock()
				defer m.Unlock()
				results[number] = batchTranslations
			}(number, batch)

			number += 1
		}

		wg.Wait()

		translations := make([]string, 0, len(requestBody.Strings))
		number = 0
		for i := 0; i < len(requestBody.Strings); i += batchSize {
			translations = append(translations, results[number]...)
			number += 1
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
