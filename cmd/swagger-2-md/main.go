package main

import (
	"fmt"
	"os"

	"github.com/lutomas/swagger-2-md/pkg/md"
	"go.uber.org/zap"

	"github.com/lutomas/swagger-2-md/pkg/config"
	"github.com/lutomas/swagger-2-md/pkg/parser"
	"github.com/lutomas/swagger-2-md/pkg/zap_logger"
	"github.com/lutomas/swagger-2-md/types"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile      = kingpin.Flag("swagger", "Path to swagger JSON file.").Default("swagger.json").Short('s').File()
	outFilePath = kingpin.Flag("md", "Path to out Markdown file.").Default("swagger.md").Short('o').String()
)

func main() {
	version := types.NewVersion("swagger-2-md")
	fmt.Printf("Version: %+v\n", *version)

	kingpin.Version("0.0.1")
	kingpin.Parse()

	fmt.Println("Input file:", (*inFile).Name())
	fmt.Println("Output file:", *outFilePath)

	cfg, err := config.LoadMainAppConfig()
	if err != nil {
		fmt.Println("FAILED LOAD CONFIG:", err.Error())
		os.Exit(1)
	}

	logger := zap_logger.GetInstanceFromConfig(&cfg.Logger)

	parser, err := parser.New(&parser.Opts{
		InFile: *inFile,
		Logger: logger,
	})
	if err != nil {
		logger.Error("FAILED CREATE PARSER", zap.Error(err))
		os.Exit(1)
	}

	swagger, err := parser.Parse()
	if err != nil {
		logger.Error("FAILED TO PARSE", zap.Error(err))
		os.Exit(1)
	}

	mdWriter, err := md.New(&md.Opts{
		OutFilePath: outFilePath,
		Logger:      logger,
		CustomCSS:   cfg.CustomCSS,
	})
	if err != nil {
		logger.Error("FAILED CREATE WRITER", zap.Error(err))
		os.Exit(1)
	}

	err = mdWriter.Write(swagger)
	if err != nil {
		logger.Error("FAILED TO WRITE", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("DONE!")
}
