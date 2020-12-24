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
	CustomCSS   string
}
type Writer struct {
	outFile *os.File
	opts    *Opts
	refsMap map[string]*types.OpenApiType
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

func (w *Writer) Write(v *types.OpenApiFileWrapper) error {
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

func (w *Writer) writeSchemas(schemas types.OpenApiSchema) (err error) {
	if schemas == nil {
		w.opts.Logger.Warn("No schemas to write.")
	}

	if w.opts.CustomCSS != "" {
		_, err = fmt.Fprintf(w.outFile, "<link rel=\"stylesheet\" type=\"text/css\" media=\"all\" href=\"%s\" /> \n\n", w.opts.CustomCSS)
		if err != nil {
			return err
		}
	}

	w.refsMap = make(map[string]*types.OpenApiType)
	// types := make([]string, 0)
	for k, v := range schemas {
		// types = append(types, k)
		w.refsMap["#/components/schemas/"+k] = v
	}

	// *********************
	// PREPARE MD FORMATS
	// *********************
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
			_, err = fmt.Fprintf(w.outFile, "## Max field depth\n%d\n\n", depth)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(w.outFile, "## Details\n%s\n\n", v.Description)
			if err != nil {
				return err
			}

			// Write props
			_, err = fmt.Fprintf(w.outFile, "| Field | Type | Mandatory | Description |\n")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w.outFile, "|:------|:------|:------|:------|\n")
			if err != nil {
				return err
			}

			subPropertiesFn := func(prefix string, in *types.MDProperty) error {
				if len(in.Properties) <= 0 {
					return nil
				}

				sort.Slice(in.Properties, func(i, j int) bool {
					return in.Properties[i].Name < in.Properties[j].Name
				})
				for _, p := range in.Properties {
					name := "*" + prefix + "*"
					if p.Name != "" {
						name = name + ".**" + p.Name + "**"
					}
					_, err = fmt.Fprintf(w.outFile, "|%s|%s|%s|%s|\n", name, p.Type, p.Mandatory, p.Description)
					if err != nil {
						return err
					}
				}
				return nil
			}

			if strings.HasPrefix(v.Type, "array") == false {
				for _, p := range v.Properties {
					_, err = fmt.Fprintf(w.outFile, "|**%s**|%s|%s|%s|\n", p.Name, p.Type, p.Mandatory, p.Description)
					if err != nil {
						return err
					}
					err = subPropertiesFn(makeSubPropertyPrefix(p), p)
					if err != nil {
						return err
					}
				}
			} else {
				// Array is special case and needs to be formatted differently.
				for _, p := range v.Properties {
					_, err = fmt.Fprintf(w.outFile, "|*%s*.**%s**|%s|%s|%s|\n", "[i]", p.Name, p.Type, p.Mandatory, p.Description)
					if err != nil {
						return err
					}

					err = subPropertiesFn(makeSubPropertyPrefix(p), p)
					if err != nil {
						return err
					}
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

func makeSubPropertyPrefix(p *types.MDProperty) string {
	prefix := p.Name
	if strings.HasPrefix(p.Type, "array") {
		prefix = prefix + "[i]"
	}
	return prefix
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

func (w *Writer) MDSchemasType(v *types.OpenApiType) *types.MDSchemasType {
	r := &types.MDSchemasType{
		O:                    v,
		Description:          prepareDescription(v.Description),
		AllOff:               len(v.AllOf) > 0,
		AdditionalProperties: v.AdditionalProperties != nil,
	}

	// Type
	r.Type = w.getMDType(v)

	w.prepareMDProperties(v, r)

	return r
}

func (w *Writer) getDescription(v *types.OpenApiType) string {
	if v == nil {
		return "--unspecified-description---"
	}
	if v.Ref != nil {
		return prepareDescription(w.refsMap[*v.Ref].Description)
	}

	return prepareDescription(v.Description)
}
func (w *Writer) getMDType(v *types.OpenApiType) string {
	if v == nil {
		return "--unknown-type---"
	}
	if len(v.AllOf) > 0 {
		return "--AllOff--"
	} else if v.Type == "" && v.AdditionalProperties != nil {
		return "--AdditionalProperties--"
	} else if v.Ref != nil {
		return w.getMDType(w.refsMap[*v.Ref])
	}

	if v.Type == "" {
		return "--unspecified-type---"
	}

	t := v.Type

	// AdditionalProperties - special case
	if v.AdditionalProperties != nil {
		t = t + "\n- format: --AdditionalProperties--"
	}

	// Check format
	if v.Format != nil {
		t = t + "\n- format: " + *v.Format
	}

	switch t {
	case "string":
		// Check enums
		if v.Enum != nil {
			t = t + "\n- one of: `" + strings.Join(v.Enum, "`, `") + "`"
		}
		// Check minlength
		if v.MinLength != nil {
			t = fmt.Sprintf("%s\n- min-length: %d", t, *v.MinLength)
		}
		// Check maxlength
		if v.MaxLength != nil {
			t = fmt.Sprintf("%s\n- max-length: %d", t, *v.MaxLength)
		}
	case "array":
		if v.Items != nil {
			t = t + "\n- items: " + w.getMDType(v.Items)
		}
	}

	return strings.ReplaceAll(t, "\n", "<br/>")
}

func (w *Writer) prepareMDProperties(o *types.OpenApiType, r types.MDProperties) {
	// AllOff
	if len(o.AllOf) > 0 {
		for _, v := range o.AllOf {
			w.prepareMDProperties(v, r)
		}
	}
	// AdditionalProperties
	if o.AdditionalProperties != nil {
		r.AddMDProperty(w.makeAdditionalMDProperty("additionalProp1", o.AdditionalProperties))
		r.AddMDProperty(w.makeAdditionalMDProperty("additionalProp2", o.AdditionalProperties))
		r.AddMDProperty(w.makeAdditionalMDProperty("additionalProp3", o.AdditionalProperties))
	}

	// Ref
	if o.Ref != nil {
		w.prepareMDProperties(w.refsMap[*o.Ref], r)
	}

	// Object
	for propName, propType := range o.Properties {
		r.AddMDProperty(w.makeMDProperty(o.Required, propName, propType))
	}

	// Array
	if o.Type == "array" && o.Items != nil {
		w.prepareMDProperties(o.Items, r)
	}
}

func (w *Writer) makeMDProperty(requiredProps []string, name string, o *types.OpenApiType) (p *types.MDProperty) {
	defer func() {
		if r := recover(); r != nil {
			w.opts.Logger.Error("failed to makeProperty", zap.String("name", name), zap.Any("obj", o), zap.Any("error", r))
		}
	}()

	p = &types.MDProperty{
		P:           o,
		Name:        name,
		Type:        w.getMDType(o),
		Mandatory:   isRequired(requiredProps, name),
		Description: w.getDescription(o),
	}

	switch o.Type {
	case "array":
		if o.Items == nil {
			panic("array without items")
		}

		if o.Items.Ref != nil {
			w.prepareMDProperties(w.refsMap[*o.Items.Ref], p)
		} else {
			switch t := o.Items.Type; t {
			case "object":
				w.prepareMDProperties(o.Items, p)
			case "string":
				p.AddMDProperty(&types.MDProperty{
					P:           o.Items,
					Name:        "",
					Type:        w.getMDType(o.Items),
					Mandatory:   "-",
					Description: w.getDescription(o.Items),
				})
			default:
				panic(fmt.Sprintf("item type '%s' not supported", t))
			}

		}
	case "object":
		if o.Ref != nil {
			w.prepareMDProperties(w.refsMap[*o.Ref], p)
		}
	}

	if o.AdditionalProperties != nil {
		// This is very special case to handle additionalProperties
		w.prepareMDProperties(o, p)
	}

	if len(o.AllOf) > 0 {
		for _, a := range o.AllOf {
			w.prepareMDProperties(a, p)
		}
	}

	return p
}

func (w *Writer) makeAdditionalMDProperty(name string, o *types.OpenApiType) *types.MDProperty {
	return &types.MDProperty{
		P:           nil,
		Name:        name,
		Type:        w.getMDType(o),
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
