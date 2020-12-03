package types

type Object struct {
	Type        string     `json:"type"`
	Format      *string    `json:"format,omitempty"`
	Description *string    `json:"description,omitempty"`
	Required    []string   `json:"required,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
}

type Properties = map[string]*PropertyType

type PropertyType struct {
	Type        string      `json:"type"`
	Format      *string     `json:"format,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Description *string     `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}
