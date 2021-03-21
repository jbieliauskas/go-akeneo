package pim

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string `json:"code"`
	Ord    int    `json:"sort_order"`
	Labels Labels `json:"labels"`
}

// AttributeGroup is an attribute group.
type AttributeGroup struct {
	Code   string   `json:"code"`
	Ord    int      `json:"sort_order"`
	Attrs  []string `json:"attributes"`
	Labels Labels   `json:"labels"`
}

// Labels is a JSON object that some Akeneo entities have.
type Labels map[string]string

// AddAttributes appends one or more attributes to group.
func (group *AttributeGroup) AddAttributes(attrs ...string) {
	if group.Attrs == nil {
		group.Attrs = make([]string, 3)
	}

	group.Attrs = append(group.Attrs, attrs...)
}