package types

import (
	"fmt"
	"os"
)

type Swagger struct {
	Openapi    string      `json:"openapi"`
	Components *Components `json:"components"`
}

func (v *Swagger) WriteComponentsSchema(stdout *os.File) {
	if v.Components != nil {
		v.Components.WriteComponentsSchema(stdout)
	}
}

type Components struct {
	Schemas Schema `json:"schemas"`
}

func (v *Components) WriteComponentsSchema(stdout *os.File) {
	if v.Schemas != nil {
		writeComponentsSchema(stdout, v.Schemas)
	}
}

type Schema = map[string]interface{}

func writeComponentsSchema(stdout *os.File, s Schema) {
	fmt.Fprintln(stdout, "TODO WRITE")
}
