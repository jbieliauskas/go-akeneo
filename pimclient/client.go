package pimclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PIMClient is a Pim client that users of this library interact with.
type PIMClient struct {
	client *http.Client
	url    string
	token  token
}

type pimDecoder interface {
	decode(interface{})
}

type pimResponse struct {
	body io.ReadCloser
	err  error
}

// New authenticates and returns PIM client
func New(url string, creds Credentials) (PIMClient, error) {
	client := new(http.Client)

	t, err := getAccessToken(client, url, creds)
	if err != nil {
		return PIMClient{}, wrapFailedError()
	}

	return PIMClient{client, url, t}, nil
}

func (c *PIMClient) list(path string, query url.Values, decodeItems decodePageItemsFunc) (Page, error) {
	url := c.url + path
	if len(query) > 0 {
		url += "?" + query.Encode()
	}

	return c.getPage(url, decodeItems)
}

func (c *PIMClient) get(path string) *pimResponse {
	req := c.authenticate(newGetRequest(c.url + path))
	return sendRequest(c.client, req)
}

func (c *PIMClient) create(path string, payload interface{}) (string, error) {
	req := c.authenticate(newJSONRequest("POST", c.url+path, payload))

	res, err := c.client.Do(req)
	if err != nil {
		return "", wrapFailedError()
	}

	res.Body.Close()

	return res.Header.Get("Location"), nil
}

func (c *PIMClient) getPage(url string, decodeItems decodePageItemsFunc) (Page, error) {
	var p Page

	req := c.authenticate(newGetRequest(url))
	res := sendRequest(c.client, req)

	err := res.err
	if err == nil {
		p, err = newPage(res.body, decodeItems)
		if err != nil {
			err = wrapFailedError()
		}
	}

	return p, err
}

func (c *PIMClient) authenticate(req *http.Request) *http.Request {
	token := fmt.Sprintf("Bearer %s", c.token.Access)
	req.Header.Add("Authorization", token)
	return req
}

func newGetRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func newJSONRequest(method, url string, payload interface{}) *http.Request {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func sendRequest(client *http.Client, req *http.Request) *pimResponse {
	res, err := client.Do(req)

	if err != nil {
		err = wrapFailedError()
	}

	return &pimResponse{res.Body, err}
}

func (res *pimResponse) decode(entity interface{}) {
	if res.err != nil {
		return
	}

	res.err = json.NewDecoder(res.body).Decode(entity)
	if res.err != nil {
		res.err = wrapFailedError()
	}

	res.body.Close()
}
