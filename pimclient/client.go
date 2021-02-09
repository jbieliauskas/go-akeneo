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
	client http.Client
	url    string
	token  token
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
	c := Client{}
	c.client = http.Client{}

	treq, err := composeTokenRequest(url, creds)
	if err != nil {
		return c, err
	}

	var t token
	err = c.do(treq, &t)
	if err != nil {
		return c, err
	}

	c.url = url
	c.token = t

	return c, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func composeTokenRequest(pimurl string, creds Credentials) (*http.Request, error) {
	form := url.Values{}
	form.Add("username", creds.User)
	form.Add("password", creds.Pass)
	form.Add("grant_type", "password")
	body := strings.NewReader(form.Encode())

	token := base64.StdEncoding.EncodeToString([]byte(creds.Cid + ":" + creds.Secret))

	req, err := http.NewRequest("POST", pimurl+"/api/oauth/v1/token", body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+token)

	return req, nil
}
