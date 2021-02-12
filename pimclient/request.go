package pimclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

type requestFactory struct {
	url string
}

func (factory *requestFactory) newGetRequest(path string, filter interface{}) *http.Request {
	values, _ := query.Values(filter)
	q := values.Encode()
	if q != "" {
		path += "?" + q
	}

	return factory.newRequest("GET", path, nil)
}

func (factory *requestFactory) newPostRequest(path string, payload interface{}) *http.Request {
	return factory.newJSONRequest("POST", path, payload)
}

func (factory *requestFactory) newJSONRequest(method, path string, payload interface{}) *http.Request {
	body, _ := json.Marshal(payload)

	req := factory.newRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (factory *requestFactory) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, factory.url+path, body)

	return req
}
