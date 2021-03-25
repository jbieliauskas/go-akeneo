package pimclient

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// Credentials is what's used to authenticate connector and get a token.
type Credentials struct {
	Cid, Secret, User, Pass string
}

type token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	TTL     int    `json:"expires_in"`
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
	res := sendRequest(client, req)
	res.decode(&t)

	return t, res.err
}
