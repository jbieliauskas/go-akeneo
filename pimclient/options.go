package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeOptions lists options of an attribute.
func (c *PIMClient) ListAttributeOptions(attr string) (Page, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	var opts pageResponse
	err := c.get(path, nil, &opts)
	if err != nil {
		return Page{}, err
	}

	return Page{opts}, nil
}

// GetAttributeOption gets an attribute option.
func (c *PIMClient) GetAttributeOption(attr, code string) (pim.AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	var opt pim.AttributeOption

	err := c.get(path, nil, &opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}

// CreateAttributeOption creates option.
func (c *PIMClient) CreateAttributeOption(attr string, opt pim.AttributeOption) (string, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.post(path, opt)
}
