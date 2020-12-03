package types

import (
	"fmt"
	"os"
)

type Swagger struct {
	Openapi    string      `json:"openapi"`
	Components *Components `json:"components"`
}

func (v *Swagger) WriteComponentsSchema(stdout *os.File) error {
	if v.Components != nil {
		return v.Components.WriteComponentsSchema(stdout)
	}

	return nil
}

type Components struct {
	Schemas Schema `json:"schemas"`
}

func (v *Components) WriteComponentsSchema(stdout *os.File) error {
	if v.Schemas != nil {
		return writeComponentsSchema(stdout, v.Schemas)
	}
	return nil
}

type Schema = map[string]interface{}

func writeComponentsSchema(stdout *os.File, s Schema) error {
	required := map[string]bool{}
	if s["required"] != nil {
		data, err := asMap(s["required"])
		if err != nil {
			return err
		}

		required = parseRequired(data)
	}

	// required, properties (type, format, description),
	fmt.Fprintln(stdout, "TODO WRITE", required)
	return nil
}

func parseRequired(data map[string]interface{}) map[string]bool {
	return map[string]bool{}
}

func asMap(i interface{}) (map[string]interface{}, error) {
	switch v := i.(type) {
	case map[string]interface{}:
		return v, nil
	}

	return nil, fmt.Errorf("unsupported type: %T", i)
}
