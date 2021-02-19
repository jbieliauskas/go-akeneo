package pimclient

import (
	"encoding/base64"
	"net/http"
)

// PIMClient is a Pim client that users of this library interact with.
type PIMClient struct {
	client *http.Client
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
func New(url string, creds Credentials) (PIMClient, error) {
	client := new(http.Client)

	t, err := getAccessToken(client, url, creds)
	if err != nil {
		return PIMClient{}, wrapFailedError()
	}

	return PIMClient{client, url, t}, nil
}

func getAccessToken(client *http.Client, url string, creds Credentials) (token, error) {
	var t token

	req := newJSONRequest("POST", url+"/api/oauth/v1/token", struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
	}{
		creds.User,
		creds.Pass,
		"password",
	})

	token := base64.StdEncoding.EncodeToString([]byte(creds.Cid + ":" + creds.Secret))
	req.Header.Set("Authorization", "Basic "+token)

	err := sendAkeneoRequest(client, req, &t)
	if err != nil {
		return t, wrapFailedError()
	}

	return t, nil
}
