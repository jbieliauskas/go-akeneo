package pimclient

import (
	"fmt"
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
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	req := c.reqFactory.newGetRequest(path, nil)

	var opt AttributeOption

	err := c.do(req, &opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}

// CreateAttributeOption creates option.
func (c *Client) CreateAttributeOption(opt AttributeOption) error {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", opt.Attr)
	req := c.reqFactory.newPostRequest(path, opt)

	return c.do(req, nil)
}
