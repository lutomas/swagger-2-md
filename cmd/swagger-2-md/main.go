package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/lutomas/swagger-2-md/pkg/config"
	"github.com/lutomas/swagger-2-md/pkg/convert"
	"github.com/lutomas/swagger-2-md/pkg/zap_logger"
	"github.com/lutomas/swagger-2-md/types"
)

func main() {
	version := types.NewVersion("swagger-2-md")
	fmt.Printf("Version: %+v\n", *version)

	cfg, err := config.LoadMainAppConfig()
	if err != nil {
		fmt.Println("FAILED LOAD CONFIG:", err.Error())
		os.Exit(1)
	}

	logger := zap_logger.GetInstanceFromConfig(&cfg.Logger)

	converter, err := convert.New(&convert.Opts{
		InPath: "",
		Logger: logger,
	})

	if err != nil {
		logger.Error("FAILED CREATE CONVERTER", zap.Error(err))
		os.Exit(1)
	}

	err = converter.Convert()
	if err != nil {
		logger.Error("FAILED TO CONVERT", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("DONE!")
}
