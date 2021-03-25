package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeOptions lists options of an attribute.
func (c *PIMClient) ListAttributeOptions(attr string) (Page, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.list(path, nil, func(d *pageItemDecoder) interface{} {
		opts := []pim.AttributeOption{}

		for d.more() {
			opt := decodeAttributeOption(d)
			opts = append(opts, opt)
		}

		return opts
	})
}

// GetAttributeOption gets an attribute option.
func (c *PIMClient) GetAttributeOption(attr, code string) (pim.AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	res := c.get(path)
	opt := decodeAttributeOption(res)
	return opt, res.err
}

// CreateAttributeOption creates option.
func (c *PIMClient) CreateAttributeOption(attr string, opt pim.AttributeOption) (string, error) {
	var ord *int
	if opt.Ord != 0 {
		ord = &opt.Ord
	}

	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.create(path, struct {
		pim.AttributeOption
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		opt,
		ord,
	})
}

func decodeAttributeOption(d pimDecoder) pim.AttributeOption {
	var opt struct {
		pim.AttributeOption
		SortOrder int `json:"sort_order"`
	}

	d.decode(&opt)
	opt.Ord = opt.SortOrder + 1

	return opt.AttributeOption
}
