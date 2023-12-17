package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Args struct {
	configFilePath string
	verbose        bool
	json           bool
}

func ParseArgs() (Args, error) {
	var args Args

	configFilePath := flag.String("config", "", "config path")
	verbose := flag.Bool("verbose", false, "verbosity")
	json := flag.Bool("json", false, "json output")

	flag.Parse()

	if *configFilePath == "" {
		return args, fmt.Errorf("config path is required")
	}

	// expand tilde if relevant
	if strings.HasPrefix(*configFilePath, "~/") {
		home, _ := os.UserHomeDir()
		*configFilePath = filepath.Join(home, (*configFilePath)[2:])
	}

	args = Args{
		configFilePath: *configFilePath,
		verbose:        *verbose,
		json:           *json,
	}
	return args, nil
}

func InitLogger(verbose bool, json bool) *slog.Logger {
	var (
		level   slog.Level
		options slog.HandlerOptions
	)
	if verbose {
		level = slog.LevelInfo
	} else {
		level = slog.LevelError
	}
	options = slog.HandlerOptions{Level: level}

	if json {
		return slog.New(slog.NewJSONHandler(os.Stdout, &options))
	} else {
		return slog.New(slog.NewTextHandler(os.Stdout, &options))
	}
}

func main() {
	args, argsErr := ParseArgs()
	if argsErr != nil {
		slog.Error(argsErr.Error())
		os.Exit(1)
	}
	logger := InitLogger(args.verbose, args.json)

	config, configErr := ReadConfig(args.configFilePath)
	if configErr != nil {
		logger.Error(configErr.Error())
		os.Exit(1)
	}

	routeErr := CheckRoutes(config, logger)
	if routeErr != nil {
		logger.Error(routeErr.Error())
		os.Exit(1)
	}
}
