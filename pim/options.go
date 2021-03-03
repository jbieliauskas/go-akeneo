package pim

// AttributeOption is an attribute option response structure.
type AttributeOption struct {
	Code   string `json:"code"`
	Ord    int    `json:"sort_order"`
	Labels Labels `json:"labels"`
}

// Labels is a JSON object that some Akeneo entities have.
type Labels map[string]string
