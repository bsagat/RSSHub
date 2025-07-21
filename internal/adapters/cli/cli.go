package cli

import (
	"RSSHub/internal/domain/ports"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/utils"
	"context"
	"fmt"
	"os"
)

type CLIHandler struct {
	aggregator ports.Aggregator
	args       []string

	log logger.Logger
}

func NewCLIHandler(aggregator ports.Aggregator, log logger.Logger) *CLIHandler {
	return &CLIHandler{
		aggregator: aggregator,
		args:       os.Args[1:],

		log: log,
	}
}

func (h *CLIHandler) ParseFlags() error {
	if len(h.args) < 1 {
		utils.PrintHelp()
		return nil
	}

	var err error
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
	case statusFlag:
		err = h.handleStatus()
	default:
		utils.PrintHelp()
		return fmt.Errorf("flag is undefined: %v", h.args[0])
	}

	if err != nil {
		h.log.Notify(err.Error())
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
