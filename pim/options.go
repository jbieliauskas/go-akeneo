package pim

import "fmt"

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string            `json:"code"`
	Ord    int               `json:"-"`
	Labels map[string]string `json:"labels,omitempty"`
}

// ListAttributeOptions lists options of an attribute.
func (c *PIMClient) ListAttributeOptions(attr string) (Page, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.list(path, nil, func(d *pageItemDecoder) interface{} {
		opts := []AttributeOption{}

		for d.more() {
			opt := decodeAttributeOption(d)
			opts = append(opts, opt)
		}

		return opts
	})
}

// GetAttributeOption gets an attribute option.
func (c *PIMClient) GetAttributeOption(attr, code string) (AttributeOption, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options/%s", attr, code)
	res := c.get(path)
	opt := decodeAttributeOption(res)
	return opt, res.err
}

// CreateAttributeOption creates option.
func (c *PIMClient) CreateAttributeOption(attr string, opt AttributeOption) (string, error) {
	var ord *int
	if opt.Ord != 0 {
		ord = &opt.Ord
	}

	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.create(path, struct {
		AttributeOption
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		opt,
		ord,
	})
}

func (c *PIMClient) UpsertAttributeOption(attr string, opt AttributeOption) (upsertAction, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	var ord *int
	if opt.Ord != 0 {
		ord = &opt.Ord
	}

	return c.upsert(path, struct {
		AttributeOption
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		opt,
		ord,
	})
}

func decodeAttributeOption(d pimDecoder) AttributeOption {
	var opt struct {
		AttributeOption
		SortOrder int `json:"sort_order"`
	}

	d.decode(&opt)
	opt.Ord = opt.SortOrder + 1

	return opt.AttributeOption
}
