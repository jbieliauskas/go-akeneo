package pim

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string `json:"code"`
	Ord    int    `json:"-"`
	Labels Labels `json:"labels,omitempty"`
}

// AttributeGroup is an attribute group.
type AttributeGroup struct {
	Code   string   `json:"code"`
	Ord    int      `json:"-"`
	Attrs  []string `json:"attributes,omitempty"`
	Labels Labels   `json:"labels,omitempty"`
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
