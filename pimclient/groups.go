package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeGroups returns a list of all attribute groups in PIM.
func (c *PIMClient) ListAttributeGroups() (Page, error) {
	return c.list("/api/rest/v1/attribute-groups", nil, func(d *pageItemDecoder) interface{} {
		gs := []pim.AttributeGroup{}

		for d.more() {
			g := decodeAttributeGroup(d)
			gs = append(gs, g)
		}

		return gs
	})
}

// GetAttributeGroup gets attribute group by code.
func (c *PIMClient) GetAttributeGroup(code string) (pim.AttributeGroup, error) {
	path := fmt.Sprintf("/api/rest/v1/attribute-groups/%s", code)
	res := c.get(path)
	g := decodeAttributeGroup(res)
	return g, res.err
}

// CreateAttributeGroup creates a group.
func (c *PIMClient) CreateAttributeGroup(g pim.AttributeGroup) (string, error) {
	var ord *int
	if g.Ord != 0 {
		ord = &g.Ord
	}

	return c.create("/api/rest/v1/attribute-groups", struct {
		pim.AttributeGroup
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		g,
		ord,
	})
}

func (c *PIMClient) UpsertAttributeGroup(g pim.AttributeGroup) (upsertAction, error) {
	var ord *int
	if g.Ord != 0 {
		ord = &g.Ord
	}

	return c.upsert("/api/rest/v1/attribute-groups", struct {
		pim.AttributeGroup
		SortOrder *int `json:"sort_order,omitempty"`
	}{
		g,
		ord,
	})
}

func decodeAttributeGroup(d pimDecoder) pim.AttributeGroup {
	var g struct {
		pim.AttributeGroup
		SortOrder int `json:"sort_order"`
	}

	d.decode(&g)
	g.Ord = g.SortOrder + 1

	return g.AttributeGroup
}
