package md

import (
	"fmt"
	"os"
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

	for k, v := range schemas {
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
	if v.Properties == nil {
		return nil
	}
	_, err = fmt.Fprintf(w.outFile, "| prop | type | mandatory | description | example |\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w.outFile, "|------|------|------|------|------|\n")
	if err != nil {
		return err
	}

	err = w.writeProperties(v.Required, v.Properties)
	if err != nil {
		return err
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

	for k, v := range properties {
		// | prop | type | mandatory | description | example |
		_, err = fmt.Fprintf(w.outFile, "|%s|%s|%s|%s|%s|\n", k, v.Type, contains(required, k), "?", "?")
	}

	return nil
}

func contains(s []string, e string) string {
	for _, a := range s {
		if a == e {
			return "yes"
		}
	}
	return "no"
}
