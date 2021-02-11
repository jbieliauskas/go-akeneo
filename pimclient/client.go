package pimclient

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is a Pim client that users of this library interact with.
type Client struct {
	client     *http.Client
	reqFactory requestFactory
	token      token
}

// Credentials is what's used to authenticate connector and get a token.
type Credentials struct {
	Cid, Secret, User, Pass string
}

type token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	TTL     int    `json:"expires_in"`
}

// New authenticates and returns PIM client
func New(url string, creds Credentials) (Client, error) {
	client := new(http.Client)

	treq := newTokenRequest(url, creds)
	var t token
	err := sendAkeneoRequest(client, treq, &t)
	if err != nil {
		return Client{}, err
	}

	return Client{client, requestFactory{url}, t}, nil
}

func (c *Client) do(req *http.Request, result interface{}) error {
	req.Header.Add("Authorization", "Bearer "+c.token.Access)

	return sendAkeneoRequest(c.client, req, result)
}

func sendAkeneoRequest(client *http.Client, req *http.Request, result interface{}) error {
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func newTokenRequest(pimurl string, creds Credentials) *http.Request {
	form := url.Values{}
	form.Add("username", creds.User)
	form.Add("password", creds.Pass)
	form.Add("grant_type", "password")
	body := strings.NewReader(form.Encode())

	token := base64.StdEncoding.EncodeToString([]byte(creds.Cid + ":" + creds.Secret))

	req, _ := http.NewRequest("POST", pimurl+"/api/oauth/v1/token", body)

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+token)

	return req
}
