package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

func run(fs *flag.FlagSet, args []string, _, stderr io.Writer) int {
	cfg, err := config.NewConfig(fs, args)
	if err != nil {
		// If the error is caused due to no source being provided, print help
		if errors.Is(err, config.ErrNoSourceProvided) {
			printHelp()
			return 1
		}
		// Otherwise, print the error
		fmt.Fprintln(stderr, "error loading config:", err)
		return 2
	}

	logger := logger.NewConsoleLogger(cfg.Verbose, nil)
	logger.Info("dish run: started")

	// Run tests on sockets
	res, err := runTests(cfg, logger)
	if err != nil {
		logger.Error(err)
		return 3
	}

	// Submit results and alerts
	alerter := alert.NewAlerter(logger)
	alerter.HandleAlerts(res.messengerText, res.results, res.failedCount, cfg)

	if res.failedCount > 0 {
		logger.Warn("dish run: some tests failed:\n", res.messengerText)
		return 4
	}

	logger.Info("dish run: all tests ok")
	return 0
}

func main() {
	os.Exit(run(flag.CommandLine, os.Args[1:], os.Stdout, os.Stderr))
}
