package md

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/lutomas/swagger-2-md/types"
)

type Opts struct {
	Logger      *zap.Logger
	OutFilePath *string
}
type Writer struct {
	outFile *os.File
	opts    *Opts
	refsMap map[string]*types.ObjectType
}

func New(opts *Opts) (*Writer, error) {
	var outFile *os.File = nil
	if opts.OutFilePath == nil {
		outFile = os.Stdout
	}

	return &Writer{
		outFile: outFile,
		opts:    opts,
	}, nil
}

func (w *Writer) Write(v *types.Swagger) error {
	if w.outFile == nil {
		f, err := os.Create(*w.opts.OutFilePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}

		defer f.Close()

		w.outFile = f

	}

	if v.Components != nil {
		return w.writeSchemas(v.Components.Schemas)
	}

	return nil
}

func (w *Writer) writeSchemas(schemas types.Schema) (err error) {
	if schemas == nil {
		w.opts.Logger.Warn("No schemas to write.")
	}

	w.refsMap = make(map[string]*types.ObjectType)
	// types := make([]string, 0)
	for k, v := range schemas {
		// types = append(types, k)
		w.refsMap["#/components/schemas/"+k] = v
	}

	res := make([]*types.MDSchemasType, 0)
	for k, v := range schemas {
		t := w.MDSchemasType(v)
		t.Name = k

		res = append(res, t)
	}

	// Sort
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	for _, v := range res {
		// Write TYPE
		_, err = fmt.Fprintf(w.outFile, "# %s \n\n", v.Name)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w.outFile, "## Type\n%s\n\n", v.Type)
		if err != nil {
			return err
		}

		if v.Description != "" {
			_, err = fmt.Fprintf(w.outFile, "## Description\n%s\n\n", v.Description)
			if err != nil {
				return err
			}
		}

		if len(v.Properties) > 0 {
			sort.Slice(v.Properties, func(i, j int) bool {
				return v.Properties[i].Name < v.Properties[j].Name
			})

			depth := getPropertiesDepth(v.Properties)
			_, err = fmt.Fprintf(w.outFile, "## Properties depth\n%d\n\n", depth)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(w.outFile, "## Details\n%s\n\n", v.Description)
			if err != nil {
				return err
			}

			// Write props
			subProps := strings.Repeat("|", depth)
			_, err = fmt.Fprintf(w.outFile, "| Field %s Type | Mandatory | Description |\n", subProps)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w.outFile, "|------|%s------|------|\n", strings.Repeat("------|", depth))
			if err != nil {
				return err
			}

			subPropertiesFn := func(in *types.MDProperty) error {
				if len(in.Properties) <= 0 {
					return nil
				}

				sort.Slice(in.Properties, func(i, j int) bool {
					return in.Properties[i].Name < in.Properties[j].Name
				})
				for _, p := range in.Properties {
					_, err = fmt.Fprintf(w.outFile, "||%s|%s|%s|%s|\n", p.Name, p.Type, p.Mandatory, p.Description)
					if err != nil {
						return err
					}
				}
				return nil
			}

			for _, p := range v.Properties {
				_, err = fmt.Fprintf(w.outFile, "|%s%s%s|%s|%s|\n", p.Name, subProps, p.Type, p.Mandatory, p.Description)
				if err != nil {
					return err
				}
				err = subPropertiesFn(p)
				if err != nil {
					return err
				}
			}

			_, err = fmt.Fprintln(w.outFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getPropertiesDepth(properties []*types.MDProperty) int {
	maxDepth := 1
	for _, p := range properties {
		depth := 1
		if len(p.Properties) > 0 {
			depth = 1 + getPropertiesDepth(p.Properties)
		}
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

func (w *Writer) MDSchemasType(v *types.ObjectType) *types.MDSchemasType {
	r := &types.MDSchemasType{
		O:                    v,
		Description:          prepareDescription(v.Description),
		AllOff:               len(v.AllOf) > 0,
		AdditionalProperties: v.AdditionalProperties != nil,
	}

	// Type
	r.Type = w.getType(v)

	w.makeProperties(v, r)

	return r
}

func (w *Writer) getType(v *types.ObjectType) string {
	if v == nil {
		return "--unknown-type---"
	}
	if len(v.AllOf) > 0 {
		return "--AllOff--"
	} else if v.AdditionalProperties != nil {
		return "--AdditionalProperties--"
	} else if v.Ref != nil {
		return w.getType(w.refsMap[*v.Ref])
	}

	if v.Type == "" {
		return "--unspecified-type---"
	}

	t := v.Type

	// Check format
	if v.Format != nil {
		t = t + "\n- format: " + *v.Format
	}

	switch t {
	case "string":
		// Check enums
		if v.Enum != nil {
			t = t + "\n- one of: " + strings.Join(v.Enum, ", ")
		}
		// Check minlength
		if v.MinLength != nil {
			t = fmt.Sprintf("%s\n- minlength: %d", t, *v.MinLength)
		}
		// Check maxlength
		if v.MaxLength != nil {
			t = fmt.Sprintf("%s\n- maxLength: %d", t, *v.MaxLength)
		}
	}

	return strings.ReplaceAll(t, "\n", "<br/>")
}

func (w *Writer) makeProperties(o *types.ObjectType, r types.Properties) {
	// AllOff
	if len(o.AllOf) > 0 {
		for _, v := range o.AllOf {
			w.makeProperties(v, r)
		}
	}
	// AdditionalProperties
	if o.AdditionalProperties != nil {
		r.AddProperty(w.makeAdditionalProperty(o.AdditionalProperties))
	}

	// Ref
	if o.Ref != nil {
		w.makeProperties(w.refsMap[*o.Ref], r)
	}

	// Object
	for propName, propType := range o.Properties {
		r.AddProperty(w.makeProperty(o.Required, propName, propType))
	}
}

func (w *Writer) makeProperty(requiredProps []string, name string, o *types.ObjectType) (p *types.MDProperty) {
	p = &types.MDProperty{
		P:           o,
		Name:        name,
		Type:        w.getType(o),
		Mandatory:   isRequired(requiredProps, name),
		Description: prepareDescription(o.Description),
	}

	switch p.Type {
	case "array":
		w.makeProperties(w.refsMap[*o.Items.Ref], p)
	case "object":
		if o.Ref != nil {
			w.makeProperties(w.refsMap[*o.Ref], p)
		}
	}

	return p
}

func (w *Writer) makeAdditionalProperty(o *types.ObjectType) *types.MDProperty {
	return &types.MDProperty{
		P:           nil,
		Name:        "--any--",
		Type:        w.getType(o),
		Mandatory:   "no",
		Description: prepareDescription(o.Description),
	}
}

func prepareDescription(description *string) string {
	if description == nil {
		return ""
	}

	return strings.ReplaceAll(*description, "\n", "<br/>")
}

func isRequired(s []string, e string) string {
	for _, a := range s {
		if a == e {
			return "yes"
		}
	}
	return "no"
}
