package pimclient

import "fmt"

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string `json:"code"`
	Ord    *int   `json:"sort_order,omitempty"`
	Labels Labels `json:"labels,omitempty"`
}

// Labels is a JSON object that some Akeneo entities have.
type Labels map[string]string

// GetAttributeOption gets an attribute option.
func (c *Client) GetAttributeOption(attr, code string) (AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	var opt AttributeOption

	err := c.get(path, nil, &opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}

// CreateAttributeOption creates option.
func (c *Client) CreateAttributeOption(attr string, opt AttributeOption) (string, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.post(path, opt)
}

// SetOrder changes sort_order property (initializes it if nil).
func (opt *AttributeOption) SetOrder(order int) {
	if opt.Ord == nil {
		opt.Ord = new(int)
	}

	*opt.Ord = order
}
