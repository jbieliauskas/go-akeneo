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

	opt.Ord++

	return opt, err
}

// CreateAttributeOption creates option.
func (c *PIMClient) CreateAttributeOption(attr string, opt pim.AttributeOption) (string, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	m := map[string]interface{}{"code": opt.Code}
	if opt.Ord > 0 {
		m["sort_order"] = opt.Ord - 1
	}
	if opt.Labels != nil && len(opt.Labels) > 0 {
		m["labels"] = opt.Labels
	}

	return c.create(path, m)
}
