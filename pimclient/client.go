package pimclient

import (
	"encoding/base64"
	"fmt"
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
	body := struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
	}{
		creds.User,
		creds.Pass,
		"password",
	}

	req := newJSONRequest("POST", url+"/api/oauth/v1/token", body)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(creds.Cid+":"+creds.Secret)),
	))

	var t token
	err := sendRequest(client, req, &t)
	return t, err
}
