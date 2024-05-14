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
	AuthTypeService     = "service"
	AuthTypeApplication = "application"
	AuthTypeUser        = "user"

	eventTypeData  = "data"
	eventTypeError = "error"

	streamTypeContent       = "Content"
	streamTypeQuotaMetadata = "QuotaMetadata"
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

	current string
	maximum string
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

func (c *client) GetQuota() (string, string) {
	return c.current, c.maximum
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

func (c *client) requestStream(ctx context.Context, method string, url string, body io.Reader) (string, error) {
	resp, err := c.request(ctx, method, url, body)
	if err != nil {
		return "", err
	}

	data, err := parseResponse(resp)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse response")
	}

	result := make([]string, 0, len(data))
	for _, d := range data {
		switch value := d.(type) {
		case contentLine:
			result = append(result, value.Content)
		case quotaMetadataLine:
			c.current = value.Updated.Current.Amount
			c.maximum = value.Updated.Maximum.Amount
			data = data[:len(data)-1]
		default:
			return "", errors.Errorf("unknown line type: %t", d)
		}
	}

	return strings.Join(result, ""), nil
}

func (c *client) buildUrl(path string) string {
	return fmt.Sprintf("https://%s/%s%s", c.domain, c.authType, path)
}

func parseResponse(response string) ([]line, error) {
	lines := strings.Split(response, "\n\n")

	result := make([]line, 0, len(lines))
	for _, l := range lines {
		if !strings.HasPrefix(l, "data: ") {
			return nil, errors.Errorf("unexpected line prefix: %s", l)
		}

		l = l[6:]
		if l == "end" {
			break
		}

		basicL := &basicLine{}
		err := json.Unmarshal([]byte(l), basicL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal basic line")
		}

		switch {
		case basicL.Type == streamTypeContent && basicL.EventType == eventTypeData:
			contentL := contentLine{}
			err = json.Unmarshal([]byte(l), &contentL)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal content line")
			}

			result = append(result, contentL)
		case basicL.Type == streamTypeQuotaMetadata && basicL.EventType == eventTypeData:
			contentL := quotaMetadataLine{}
			err = json.Unmarshal([]byte(l), &contentL)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal quota metadata line")
			}

			result = append(result, contentL)
		case basicL.EventType == eventTypeError:
			return nil, errors.Errorf("server error: %s", basicL.ErrorMessage)
		}
	}

	return result, nil
}

type line interface {
	GetType() string
}

type basicLine struct {
	Type         string `json:"type"`
	EventType    string `json:"event_type"`
	ErrorMessage string `json:"error_message"`
}

func (l basicLine) GetType() string {
	return l.Type
}

type quotaMetadataLine struct {
	basicLine
	Updated quotaMetadata `json:"updated"`
	Spent   amount        `json:"spent"`
}

type quotaMetadata struct {
	License string  `json:"license"`
	Current amount  `json:"current"`
	Maximum amount  `json:"maximum"`
	Until   int64   `json:"until"`
	QuotaID quotaID `json:"quotaID"`
}

type amount struct {
	Amount string `json:"amount"`
}

type quotaID struct {
	QuotaID string `json:"quotaId"`
}

type contentLine struct {
	basicLine
	Content string `json:"content"`
}
