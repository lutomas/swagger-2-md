package types

type Swagger struct {
	Openapi    string      `json:"openapi"`
	Components *Components `json:"components"`
}

type Components struct {
	Schemas map[string]interface{} `json:"schemas"`
}

type Schema = map[string]interface{}
