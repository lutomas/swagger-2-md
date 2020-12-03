package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/lutomas/swagger-2-md/pkg/config"
	"github.com/lutomas/swagger-2-md/pkg/convert"
	"github.com/lutomas/swagger-2-md/pkg/zap_logger"
	"github.com/lutomas/swagger-2-md/types"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile = kingpin.Flag("swagger", "Path to swagger JSON file.").Default("swagger.json").Short('s').File()
)

func main() {
	version := types.NewVersion("swagger-2-md")
	fmt.Printf("Version: %+v\n", *version)

	kingpin.Version("0.0.1")
	kingpin.Parse()

	cfg, err := config.LoadMainAppConfig()
	if err != nil {
		fmt.Println("FAILED LOAD CONFIG:", err.Error())
		os.Exit(1)
	}

	logger := zap_logger.GetInstanceFromConfig(&cfg.Logger)

	converter, err := convert.New(&convert.Opts{
		InFile: *inFile,
		Logger: logger,
	})

	if err != nil {
		logger.Error("FAILED CREATE CONVERTER", zap.Error(err))
		os.Exit(1)
	}

	swager, err := converter.Convert()
	if err != nil {
		logger.Error("FAILED TO CONVERT", zap.Error(err))
		os.Exit(1)
	}

	swager.WriteComponentsSchema(os.Stdout)

	logger.Info("DONE!")
}
