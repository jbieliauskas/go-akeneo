package pimclient

import (
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

func (factory *requestFactory) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, factory.url+path, body)

	return req
}
