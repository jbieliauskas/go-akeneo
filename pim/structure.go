package pim

import "fmt"

// AttributeGroup is an attribute group.
type AttributeGroup struct {
	Code   string            `json:"code"`
	Ord    SortOrder         `json:"-"`
	Attrs  []string          `json:"attributes,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string            `json:"code"`
	Ord    SortOrder         `json:"-"`
	Labels map[string]string `json:"labels,omitempty"`
}

type SortOrder uint16

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
	return c.create("/api/rest/v1/attribute-groups", struct {
		AttributeGroup
		OrdPtr *SortOrder `json:"sort_order,omitempty"`
	}{
		g,
		convertSortOrderToPointer(g.Ord),
	})
}

func (c *PIMClient) UpsertAttributeGroup(g AttributeGroup) (upsertAction, error) {
	return c.upsert("/api/rest/v1/attribute-groups", struct {
		AttributeGroup
		OrdPtr *SortOrder `json:"sort_order,omitempty"`
	}{
		g,
		convertSortOrderToPointer(g.Ord),
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
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.create(path, struct {
		AttributeOption
		OrdPtr *SortOrder `json:"sort_order,omitempty"`
	}{
		opt,
		convertSortOrderToPointer(opt.Ord),
	})
}

func (c *PIMClient) UpsertAttributeOption(attr string, opt AttributeOption) (upsertAction, error) {
	path := fmt.Sprintf("/api/rest/v1/attributes/%s/options", attr)

	return c.upsert(path, struct {
		AttributeOption
		OrdPtr *SortOrder `json:"sort_order,omitempty"`
	}{
		opt,
		convertSortOrderToPointer(opt.Ord),
	})
}

func decodeAttributeGroup(d pimDecoder) AttributeGroup {
	var g struct {
		AttributeGroup
		SortOrd SortOrder `json:"sort_order"`
	}
	d.decode(&g)

	g.Ord = g.SortOrd + 1

	return g.AttributeGroup
}

func decodeAttributeOption(d pimDecoder) AttributeOption {
	var opt struct {
		AttributeOption
		SortOrd SortOrder `json:"sort_order"`
	}
	d.decode(&opt)

	opt.Ord = opt.SortOrd + 1

	return opt.AttributeOption
}

func convertSortOrderToPointer(ord SortOrder) *SortOrder {
	if ord == 0 {
		return nil
	}
	return &ord
}
