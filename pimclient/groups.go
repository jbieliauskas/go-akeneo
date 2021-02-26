package pimclient

import (
	"fmt"

	"github.com/jbieliauskas/go-akeneo/pim"
)

// ListAttributeGroups returns a list of all attribute groups in PIM.
func (c *PIMClient) ListAttributeGroups() (Page, error) {
	var groups pageResponse
	err := c.get("/api/rest/v1/attribute-groups", nil, &groups)
	if err != nil {
		return Page{}, err
	}

	return Page{groups}, nil
}

// GetAttributeGroup gets attribute group by code.
func (c *PIMClient) GetAttributeGroup(code string) (pim.AttributeGroup, error) {
	path := fmt.Sprintf("/api/rest/v1/attribute-groups/%s", code)
	var group pim.AttributeGroup

	err := c.get(path, nil, &group)
	if err != nil {
		return group, err
	}

	return group, nil
}

// CreateAttributeGroup creates a group.
func (c *PIMClient) CreateAttributeGroup(group pim.AttributeGroup) (string, error) {
	return c.post("/api/rest/v1/attribute-groups", group)
}
