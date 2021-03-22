package pimclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *PIMClient) list(path string, query url.Values, decodeItems func(d pageItemDecoder) interface{}) (Page, error) {
	url := c.url + path
	if len(query) > 0 {
		url += "?" + query.Encode()
	}

	return c.getPage(url, decodeItems)
}

func (c *PIMClient) get(path string, result interface{}) error {
	req := c.authenticate(newGetRequest(c.url + path))
	return sendRequest(c.client, req, result)
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

func (c *PIMClient) getPage(url string, decodeItems func(d pageItemDecoder) interface{}) (Page, error) {
	req := c.authenticate(newGetRequest(url))

	res, err := c.client.Do(req)
	if err != nil {
		return Page{}, wrapFailedError()
	}

	p, err := newPage(res.Body, decodeItems)
	if err != nil {
		return Page{}, wrapFailedError()
	}

	return p, nil
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

func sendRequest(client *http.Client, req *http.Request, result interface{}) error {
	res, err := client.Do(req)
	if err != nil {
		return wrapFailedError()
	}

	defer res.Body.Close()

	if result != nil {
		err := json.NewDecoder(res.Body).Decode(result)
		if err != nil {
			return wrapFailedError()
		}
	}

	return nil
}
