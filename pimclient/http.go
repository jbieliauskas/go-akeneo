package pimclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *PIMClient) get(path string, query url.Values, result interface{}) error {
	url := c.url + path
	if len(query) > 0 {
		url += "?" + query.Encode()
	}

	return c.getViaFullURL(url, result)
}

func (c *PIMClient) getViaFullURL(url string, result interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)

	token := fmt.Sprintf("Bearer %s", c.token.Access)
	req.Header.Add("Authorization", token)

	return sendAkeneoRequest(c.client, req, result)
}

func (c *PIMClient) post(path string, payload interface{}) (string, error) {
	req := newJSONRequest("POST", c.url+path, payload)
	req.Header.Set("Authorization", "Bearer "+c.token.Access)

	res, err := c.client.Do(req)
	if err != nil {
		return "", wrapFailedError()
	}

	res.Body.Close()

	return res.Header.Get("Location"), nil
}

func newJSONRequest(method, url string, payload interface{}) *http.Request {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func sendAkeneoRequest(client *http.Client, req *http.Request, result interface{}) error {
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
