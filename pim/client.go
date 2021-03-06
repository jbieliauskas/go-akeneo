package pim

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

type encodable interface {
	encode(*json.Encoder)
}

type pimResponse struct {
	status int
	body   io.ReadCloser
	err    error
}

type upsertAction uint8

const (
	Create = upsertAction(1)
	Update = upsertAction(2)
)

// NewClient authenticates and returns PIM client
func NewClient(url string, creds Credentials) (PIMClient, error) {
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

func (c *PIMClient) create(path string, payload interface{}) error {
	req := c.authenticate(newJSONRequest("POST", c.url+path, payload))
	res := sendRequest(c.client, req)
	res.body.Close()
	return res.err
}

func (c *PIMClient) upsert(path string, payload interface{}) (upsertAction, error) {
	req := c.authenticate(newJSONRequest("PATCH", c.url+path, payload))
	res := sendRequest(c.client, req)
	res.body.Close()

	err := res.err
	if err == nil {
		if res.status == http.StatusCreated {
			return Create, nil
		}
		if res.status == http.StatusNoContent {
			return Update, nil
		}

		err = wrapFailedError()
	}

	return upsertAction(0), err
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

func newJSONRequest(method, url string, v interface{}) *http.Request {
	body := bytes.NewBuffer([]byte{})
	e := json.NewEncoder(body)

	if entity, ok := v.(encodable); ok {
		entity.encode(e)
	} else {
		e.Encode(v)
	}

	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func sendRequest(client *http.Client, req *http.Request) *pimResponse {
	res, err := client.Do(req)

	if err != nil {
		err = wrapFailedError()
	}

	return &pimResponse{res.StatusCode, res.Body, err}
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
