package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeOptions lists options of an attribute.
func (c *PIMClient) ListAttributeOptions(attr string) (Page, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)
	return c.list(path, nil)
}

// GetAttributeOption gets an attribute option.
func (c *PIMClient) GetAttributeOption(attr, code string) (pim.AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)

	var opt pim.AttributeOption
	err := c.get(path, &opt)

	return opt, err
}

// CreateAttributeOption creates option.
func (c *PIMClient) CreateAttributeOption(attr string, opt pim.AttributeOption) (string, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)
	return c.create(path, opt)
}
