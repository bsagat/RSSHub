package cli

import (
	"RSSHub/config"
	"RSSHub/internal/domain/ports"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/utils"
	"context"
	"flag"
	"fmt"
	"os"
)

type CLIHandler struct {
	aggregator ports.Aggregator
	args       []string

	cfg config.CLI_APP
	log logger.Logger
}

func NewCLIHandler(aggregator ports.Aggregator, cfg config.CLI_APP, log logger.Logger) *CLIHandler {
	return &CLIHandler{
		aggregator: aggregator,
		args:       os.Args[1:],

		cfg: cfg,
		log: log,
	}
}

func (h *CLIHandler) ParseFlags() error {
	var err error
	flag.Parse()

	if *helpFlag {
		utils.PrintHelp()
		return nil
	}

	if len(h.args) < 1 {
		utils.PrintHelp()
		return nil
	}

	if h.args[0] == "--race" {
		h.args = h.args[1:]
	}

	switch h.args[0] {
	case fetchFlag:
		err = h.handleFetch()
	case addFlag:
		err = h.handleAdd()
	case setIntervalFlag:
		err = h.handleInterval()
	case setWorkerFlag:
		err = h.handleWorkers()
	case listFlag:
		err = h.handleList()
	case deleteFlag:
		err = h.handleDelete()
	case articlesFlag:
		err = h.handleArticle()
	default:
		return fmt.Errorf("flag is undefined: %v", h.args[0])
	}
	if err != nil {
		return err
	}

	return nil
}

func (h *CLIHandler) Close() error {
	if err := h.aggregator.Stop(); err != nil {
		h.log.Error(context.Background(), "Failed to close aggregator service", "error", err)
		return err
	}
	return nil
}
