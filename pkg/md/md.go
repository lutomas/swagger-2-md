package md

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lutomas/swagger-2-md/types"
	"go.uber.org/zap"
)

type Opts struct {
	Logger      *zap.Logger
	OutFilePath *string
}
type Writer struct {
	outFile *os.File
	opts    *Opts
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

	types := make([]string, 0)
	for k, _ := range schemas {
		types = append(types, k)
	}

	// Sort prop names
	sort.Strings(types)

	for _, k := range types {
		v := schemas[k]
		// Write TYPE
		_, err = fmt.Fprintf(w.outFile, "# %s \n\n", k)
		if err != nil {
			return err
		}

		if v.Description != nil {
			_, err = fmt.Fprintf(w.outFile, "%s\n\n", strings.ReplaceAll(*v.Description, "\n", "<br/>"))
			if err != nil {
				return err
			}
		}

		if err = w.writeObjectType(v); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) writeObjectType(v *types.ObjectType) (err error) {

	if v.Type != "object" {
		_, err = fmt.Fprintf(w.outFile, "| Type |\n")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w.outFile, "|------|\n")
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w.outFile, "| %s |\n\n", preparePropertyType(v))
		if err != nil {
			return err
		}
	}

	for _, inc := range v.AllOf {
		if inc.Ref != nil {
			_, err = fmt.Fprintf(w.outFile, "## INCLUDE MANUALLY: %s\n\n", *inc.Ref)
			if err != nil {
				return err
			}
		}
	}

	if v.Properties == nil && len(v.AllOf) == 0 {
		return nil
	}
	_, err = fmt.Fprintf(w.outFile, "| Field | Type | Mandatory | Description |\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w.outFile, "|------|------|------|------|\n")
	if err != nil {
		return err
	}

	return w.writeIncludedObjectType(v)
}

func (w *Writer) writeIncludedObjectType(v *types.ObjectType) (err error) {
	err = w.writeProperties(v.Required, v.Properties)
	if err != nil {
		return err
	}

	if len(v.AllOf) > 0 {
		for _, inc := range v.AllOf {
			w.writeIncludedObjectType(inc)
		}
	}

	_, err = fmt.Fprintln(w.outFile)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeProperties(required []string, properties map[string]*types.ObjectType) (err error) {
	if properties == nil {
		return nil
	}

	propNames := make([]string, 0)
	for k, _ := range properties {
		propNames = append(propNames, k)
	}

	// Sort prop names
	sort.Strings(propNames)

	for _, k := range propNames {
		v := properties[k]
		// | prop | type | mandatory | description | example |
		_, err = fmt.Fprintf(w.outFile, "|%s|%s|%s|%s|\n", k, preparePropertyType(v), isRequired(required, k), prepareDescription(v.Description))
	}

	return nil
}

func preparePropertyType(v *types.ObjectType) string {
	t := v.Type
	if t == "" {
		if v.Ref != nil {
			return *v.Ref
		}
		// Compound type
		if len(v.AllOf) > 0 {
			return "--AllOf--"
		}
		return "?"
	}

	if t == "array" {
		itemType := "?"
		// Check whats the type of array items
		if v.Items != nil {
			itemType = preparePropertyType(v.Items)
		}
		t = fmt.Sprintf("%s [%s]", t, itemType)
		return t
	}

	if v.Format != nil {
		t = fmt.Sprintf("%s (%s)", t, *v.Format)
	}

	if len(v.Enum) > 0 {
		enums := make([]string, 0)
		for _, enum := range v.Enum {
			enums = append(enums, "`"+enum+"`")
		}

		sort.Strings(enums)
		e := strings.Join(enums, "<br/> - ")

		t = t + "<br/><br/>*allowed values:<br/> - " + e + "*"
	}

	return t
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
