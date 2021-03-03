package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeGroups returns a list of all attribute groups in PIM.
func (c *PIMClient) ListAttributeGroups() (Page, error) {
	return c.list("/api/rest/v1/attribute-groups", nil)
}

// GetAttributeGroup gets attribute group by code.
func (c *PIMClient) GetAttributeGroup(code string) (pim.AttributeGroup, error) {
	path := fmt.Sprintf("/api/rest/v1/attribute-groups/%s", code)

	var group pim.AttributeGroup
	err := c.get(path, &group)

	group.Ord++

	return group, err
}

// CreateAttributeGroup creates a group.
func (c *PIMClient) CreateAttributeGroup(group pim.AttributeGroup) (string, error) {
	m := map[string]interface{}{"code": group.Code}
	if group.Ord > 0 {
		m["sort_order"] = group.Ord - 1
	}
	if group.Attrs != nil && len(group.Attrs) > 0 {
		m["attributes"] = group.Attrs
	}
	if group.Labels != nil && len(group.Labels) > 0 {
		m["labels"] = group.Labels
	}

	return c.create("/api/rest/v1/attribute-groups", m)
}
