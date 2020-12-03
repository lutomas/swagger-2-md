package types

type ObjectType struct {
	Type        string                 `json:"type"`
	Format      *string                `json:"format,omitempty"`
	Description *string                `json:"description,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Ref         *string                `json:"$ref,omitempty"`
	Enum        []string               `json:"enum,omitempty"`
	Properties  map[string]*ObjectType `json:"properties,omitempty"`
	Example     interface{}            `json:"example,omitempty"`
	AllOf       []*ObjectType          `json:"allOf,omitempty"`
}
