package types

type OpenApiFileWrapper struct {
	Openapi    string             `json:"openapi"`
	Paths      OpenApiPaths       `json:"paths"`
	Components *OpenApiComponents `json:"components"`
}

type OpenApiComponents struct {
	Schemas OpenApiSchema `json:"schemas"`
}

type OpenApiSchema = map[string]*OpenApiType
type OpenApiPaths = map[string]*OpenApiPath

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

type OpenApiPath struct {
	Get    *OpenApiPathDetails `json:"get,omitempty"`
	Post   *OpenApiPathDetails `json:"post,omitempty"`
	Delete *OpenApiPathDetails `json:"delete,omitempty"`
	Update *OpenApiPathDetails `json:"update,omitempty"`
}

type OpenApiPathDetails struct {
	Description *string `json:"description,omitempty"`
}