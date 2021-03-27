package pim

import "fmt"

// AttributeGroup is an attribute group.
type AttributeGroup struct {
	Code   string            `json:"code"`
	Ord    int               `json:"-"`
	Attrs  []string          `json:"attributes,omitempty"`
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

func decodeAttributeGroup(d pimDecoder) AttributeGroup {
	var g struct {
		AttributeGroup
		SortOrder int `json:"sort_order"`
	}

	d.decode(&g)
	g.Ord = g.SortOrder + 1

	return g.AttributeGroup
}
