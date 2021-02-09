package pimclient

import (
	"fmt"
	"net/http"
)

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Attr   string `json:"attribute"`
	Code   string `json:"code"`
	Order  int    `json:"sort_order"`
	Labels Labels `json:"labels"`
}

// Labels is a JSON object that some Akeneo entities have.
type Labels map[string]string

// GetAttributeOption gets an attribute option.
func (c *Client) GetAttributeOption(attr, code string) (AttributeOption, error) {
	var opt AttributeOption

	req, err := composeOptionRequest(c.url, attr, code, c.token.Access)
	if err != nil {
		return opt, err
	}

	err = c.do(req, &opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}

func composeOptionRequest(pimurl, attr, opt, token string) (*http.Request, error) {
	query := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, opt)

	req, err := http.NewRequest("GET", pimurl+query, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}
