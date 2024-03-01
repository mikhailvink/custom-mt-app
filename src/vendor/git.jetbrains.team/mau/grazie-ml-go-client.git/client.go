package zendeskgo_sell

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	AuthTypeService = "service"
	AuthTypeUser    = "user"

	eventTypeData  = "data"
	eventTypeError = "error"
)

type GrazieAgent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type client struct {
	domain   string
	authType string
	token    string
	agent    string
}

func New(domain string, authType string, token string, agent GrazieAgent) (Client, error) {
	agentStr, err := json.Marshal(agent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal agent")
	}

	return &client{
		domain:   domain,
		authType: authType,
		token:    token,
		agent:    string(agentStr),
	}, nil
}

func (c *client) rawRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request")
	}

	req.Header.Set("Grazie-Authenticate-JWT", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Grazie-Agent", c.agent)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do request")
	}

	if resp.StatusCode >= 400 {
		return resp, errors.Errorf("response status is #%d: %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}
func (c *client) request(ctx context.Context, method string, url string, body io.Reader) (string, error) {
	resp, err := c.rawRequest(ctx, method, url, body)

	var respBody []byte
	var respErr error
	if resp != nil && resp.Body != nil {
		defer func() { _ = resp.Body.Close() }()

		respBody, respErr = ioutil.ReadAll(resp.Body)
		if respErr != nil {
			return "", errors.Wrap(respErr, "cannot read response body")
		}
	}

	if err != nil {
		return "", errors.Wrap(err, string(respBody))
	}

	return string(respBody), nil
}

func (c *client) buildUrl(path string) string {
	return fmt.Sprintf("https://%s/%s%s", c.domain, c.authType, path)
}

type responseLine struct {
	EventType    string `json:"event_type"`
	Current      string `json:"current"`
	ErrorMessage string `json:"error_message"`
}

func parseChatResponse(response string) ([]string, error) {
	lines, err := parseResponse(response)
	if err != nil {
		return []string{""}, errors.Wrap(err, "cannot parse response")
	}
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		rLine := &responseLine{}
		err := json.Unmarshal([]byte(line), rLine)
		if err != nil {
			return []string{""}, errors.Wrap(err, "cannot unmarshal response line")
		}

		switch rLine.EventType {
		case eventTypeData:
			result = append(result, rLine.Current)
		case eventTypeError:
			return []string{""}, errors.Errorf("server error: %s", rLine.ErrorMessage)
		default:
			return []string{""}, errors.Errorf("unknown event type: %s", rLine.EventType)
		}
	}
	return result, nil
}

func parseResponse(response string) ([]string, error) {
	lines := strings.Split(response, "\n\n")

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "data: ") {
			return []string{""}, errors.Errorf("unexpected line prefix: %s", line)
		}

		line = line[6:]
		if line == "end" {
			break
		}
		result = append(result, line)
	}

	return result, nil
}
