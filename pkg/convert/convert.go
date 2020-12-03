package convert

import (
	"fmt"

	"go.uber.org/zap"
)

type Opts struct {
	InPath string
	Logger *zap.Logger
}
type Converter struct {
	opts *Opts
}

func New(opts *Opts) (*Converter, error) {
	return &Converter{
		opts: opts,
	}, nil
}

func (c *Converter) Convert() error {
	return fmt.Errorf("not implemented")
}
