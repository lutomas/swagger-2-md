package convert

import (
	"io/ioutil"
	"os"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/lutomas/swagger-2-md/types"
)

type Opts struct {
	InFile *os.File
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

func (c *Converter) Convert() (*types.Swagger, error) {
	jsonFile, err := os.Open(c.opts.InFile.Name())
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	data := &types.Swagger{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = easyjson.Unmarshal(byteValue, data)
	if err != nil {
		return nil, err
	}

	c.opts.Logger.Info("data", zap.Any("swagger", data))

	return data, nil
}
