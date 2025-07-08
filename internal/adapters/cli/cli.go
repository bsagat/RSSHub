package cli

import (
	"RSSHub/internal/domain/ports"
	"RSSHub/internal/pkg/utils"
	"flag"
	"log/slog"
	"os"
)

type CLIHandler struct {
	aggregator ports.Aggregator
	args       []string

	log *slog.Logger
}

func NewCLIHandler(aggregator ports.Aggregator, log *slog.Logger) *CLIHandler {
	return &CLIHandler{
		aggregator: aggregator,
		args:       os.Args[1:],
		log:        log,
	}
}

func (h *CLIHandler) ParseFlags() int {
	var err error
	flag.Parse()

	if *helpFlag {
		utils.PrintHelp()
		return OkStatusCode
	}

	if len(h.args) < 1 {
		h.log.Error("missing CLI arguments", "args", h.args)
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
		h.log.Error("flag is undefined", "flag", h.args[0])
		return ErrStatusCode
	}
	if err != nil {
		return ErrStatusCode
	}
	return OkStatusCode
}

func (h *CLIHandler) Close() error {
	if err := h.aggregator.Stop(); err != nil {
		h.log.Error("Failed to close aggregator service", "error", err)
		return err
	}
	return nil
}
