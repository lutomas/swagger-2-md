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

type Schema = map[string]*ObjectType

func writeComponentsSchema(stdout *os.File, s Schema) error {
	// required, properties (type, format, description),
	fmt.Fprintln(stdout, "TODO WRITE")
	return nil
}
