package cli

import (
	"RSSHub/config"
	"RSSHub/internal/domain/ports"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/utils"
	"context"
	"flag"
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

func (h *CLIHandler) ParseFlags() int {
	const fn = "CLIHandler.ParseFlags"
	log := h.log.GetSlogLogger().With("fn", fn)

	var err error
	flag.Parse()

	if *helpFlag {
		utils.PrintHelp()
		return OkStatusCode
	}

	if len(h.args) < 1 {
		log.Error("missing CLI arguments", "args", h.args)
		utils.PrintHelp()
		return ErrStatusCode
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
		log.Error("flag is undefined", "flag", h.args[0])
		return ErrStatusCode
	}
	if err != nil {
		return ErrStatusCode
	}
	return OkStatusCode
}

func (h *CLIHandler) Close() error {
	if err := h.aggregator.Stop(); err != nil {
		h.log.Error(context.Background(), "Failed to close aggregator service", "error", err)
		return err
	}
	return nil
}
