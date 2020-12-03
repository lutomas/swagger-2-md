package types

type Object struct {
	Type        string     `json:"type"`
	Description *string    `json:"description"`
	Required    []string   `json:"required"`
	Properties  Properties `json:"properties"`
}

type Properties = map[string]*PropertyType

type PropertyType struct {
	Type        string      `json:"type"`
	Description *string     `json:"description"`
	Example     interface{} `json:"example"`
}
