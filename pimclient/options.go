package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// GetAttributeOption gets an attribute option.
func (c *Client) GetAttributeOption(attr, code string) (pim.AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	var opt pim.AttributeOption

	err := c.get(path, nil, &opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}

// CreateAttributeOption creates option.
func (c *Client) CreateAttributeOption(attr string, opt pim.AttributeOption) (string, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.post(path, opt)
}
