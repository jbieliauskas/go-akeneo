package pim

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string `json:"code"`
	Ord    *int   `json:"sort_order,omitempty"`
	Labels Labels `json:"labels,omitempty"`
}

// Labels is a JSON object that some Akeneo entities have.
type Labels map[string]string

// SetOrder changes sort_order property (initializes it if nil).
func (opt *AttributeOption) SetOrder(order int) {
	if opt.Ord == nil {
		opt.Ord = new(int)
	}

	*opt.Ord = order
}
