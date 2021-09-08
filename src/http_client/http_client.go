package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type PreprocessRequestFunc func(r *http.Request)

type requestDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type HttpClient struct {
	client     requestDoer
	preprocess []PreprocessRequestFunc
	IsDebug    bool
}

func New(client requestDoer) HttpClient {
	return HttpClient{
		client: client,
	}
}

func (c *HttpClient) AddPreprocessFunc(fn PreprocessRequestFunc) {
	if fn != nil {
		c.preprocess = append(c.preprocess, fn)
	}
}

func (c HttpClient) Request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	c.debug("Making %q request to url: %s", method, url)
	if err != nil {
		return nil, err
	}

	if len(c.preprocess) > 0 {
		for _, fn := range c.preprocess {
			c.debug("Calling preprocess func...")
			fn(req)
		}
	}

	c.debug("request headers: %v", req.Header)

	return c.client.Do(req)
}

func (c HttpClient) UnmarshallJSON(method string, url string, body io.Reader, data interface{}) error {
	resp, err := c.Request(method, url, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := ioutil.ReadAll(resp.Body)
	c.debug("recieved body: %s", respBody)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("response code %d: %s", resp.StatusCode, resp.Status)
	}

	return json.NewDecoder(bytes.NewReader(respBody)).Decode(data)
}

func CreateHeaderSetterPreprocessor(key, value string) PreprocessRequestFunc {
	return func(r *http.Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}

		r.Header.Set(key, value)
	}
}

func (c HttpClient) debug(pattern string, args ...interface{}) {
	if c.IsDebug {
		logrus.Debugf(pattern, args...)
	}
}
