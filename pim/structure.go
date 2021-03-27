package pim

import "fmt"

// AttributeGroup is an attribute group.
type AttributeGroup struct {
	Code   string            `json:"code"`
	Ord    int               `json:"-"`
	Attrs  []string          `json:"attributes,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string            `json:"code"`
	Ord    int               `json:"-"`
	Labels map[string]string `json:"labels,omitempty"`
}

// ListAttributeGroups returns a list of all attribute groups in PIM.
func (c *PIMClient) ListAttributeGroups() (Page, error) {
	return c.list("/api/rest/v1/attribute-groups", nil, func(d *pageItemDecoder) interface{} {
		gs := []AttributeGroup{}

		for d.more() {
			g := decodeAttributeGroup(d)
			gs = append(gs, g)
		}

		return gs
	})
}

// GetAttributeGroup gets attribute group by code.
func (c *PIMClient) GetAttributeGroup(code string) (AttributeGroup, error) {
	path := fmt.Sprintf("/api/rest/v1/attribute-groups/%s", code)
	res := c.get(path)
	g := decodeAttributeGroup(res)
	return g, res.err
}

// CreateAttributeGroup creates a group.
func (c *PIMClient) CreateAttributeGroup(g AttributeGroup) (string, error) {
	var ord *int
	if g.Ord != 0 {
		ord = &g.Ord
	}

	return c.create("/api/rest/v1/attribute-groups", struct {
		AttributeGroup
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		g,
		ord,
	})
}

func (c *PIMClient) UpsertAttributeGroup(g AttributeGroup) (upsertAction, error) {
	var ord *int
	if g.Ord != 0 {
		ord = &g.Ord
	}

	return c.upsert("/api/rest/v1/attribute-groups", struct {
		AttributeGroup
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		g,
		ord,
	})
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

func decodeAttributeGroup(d pimDecoder) AttributeGroup {
	var g struct {
		AttributeGroup
		SortOrder int `json:"sort_order"`
	}

	d.decode(&g)
	g.Ord = g.SortOrder + 1

	return g.AttributeGroup
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
