package pimclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *PIMClient) list(path string, query url.Values, decodeItems func(d pageItemDecoder) interface{}) (Page, error) {
	url := c.url + path
	if len(query) > 0 {
		url += "?" + query.Encode()
	}
	req := c.newGetRequest(url)

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

func (c *PIMClient) get(path string, result interface{}) error {
	req := c.newGetRequest(c.url + path)
	return c.send(req, result)
}

func (c *PIMClient) create(path string, payload interface{}) (string, error) {
	req := newJSONRequest("POST", c.url+path, payload)
	req.Header.Set("Authorization", "Bearer "+c.token.Access)

	res, err := c.client.Do(req)
	if err != nil {
		return "", wrapFailedError()
	}

	res.Body.Close()

	return res.Header.Get("Location"), nil
}

func (c *PIMClient) newGetRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	return c.addAuth(req)
}

func (c *PIMClient) newJSONRequest(method, url string, payload interface{}) *http.Request {
	req := newJSONRequest(method, url, payload)
	return c.addAuth(req)
}

func newJSONRequest(method, url string, payload interface{}) *http.Request {
	b, _ := json.Marshal(payload)
	body := bytes.NewReader(b)

	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (c *PIMClient) addAuth(req *http.Request) *http.Request {
	token := fmt.Sprintf("Bearer %s", c.token.Access)
	req.Header.Add("Authorization", token)

	return req
}

func (c *PIMClient) send(req *http.Request, result interface{}) error {
	return sendRequest(c.client, req, result)
}

func sendRequest(client *http.Client, req *http.Request, result interface{}) error {
	res, err := client.Do(req)
	if err != nil {
		return wrapFailedError()
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return wrapFailedError()
	}

	if result == nil {
		return nil
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return wrapFailedError()
	}

	return nil
}
