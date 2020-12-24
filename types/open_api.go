package types

type OpenApiFileWrapper struct {
	Openapi    string      `json:"openapi"`
	Components *Components `json:"components"`
}

type Components struct {
	Schemas Schema `json:"schemas"`
}

type Schema = map[string]*OpenApiType

type OpenApiType struct {
	Type                 string                  `json:"type"`
	Format               *string                 `json:"format,omitempty"`
	Description          *string                 `json:"description,omitempty"`
	Required             []string                `json:"required,omitempty"`
	Ref                  *string                 `json:"$ref,omitempty"`
	Enum                 []string                `json:"enum,omitempty"`
	Properties           map[string]*OpenApiType `json:"properties,omitempty"`
	Example              interface{}             `json:"example,omitempty"`
	AllOf                []*OpenApiType          `json:"allOf,omitempty"`
	Items                *OpenApiType            `json:"items,omitempty"`
	AdditionalProperties *OpenApiType            `json:"additionalProperties,omitempty"`
	MaxLength            *int64                  `json:"maxLength,omitempty"`
	MinLength            *int64                  `json:"minLength,omitempty"`
}

func (v *OpenApiType) IsObject() bool {
	return v.Type == "object"
}
