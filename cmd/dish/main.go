package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		// If the error is caused due to no source being provided, print help
		if errors.Is(err, config.ErrNoSourceProvided) {
			printHelp()
			os.Exit(1)
		}
		// Otherwise, print the error
		log.Print("error loading config: ", err)
		return
	}

	logger := logger.NewConsoleLogger(cfg.Verbose)

	logger.Info("dish run: started")

	// Run tests on sockets
	res, err := runTests(cfg, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	// Submit results and alerts
	alerter := alert.NewAlerter(logger)
	alerter.HandleAlerts(res.messengerText, res.results, res.failedCount, cfg)

	if res.failedCount > 0 {
		logger.Warn("dish run: some tests failed:\n", res.messengerText)
		return
	}

	logger.Info("dish run: all tests ok")
}
