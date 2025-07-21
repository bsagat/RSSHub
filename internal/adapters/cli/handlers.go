package cli

import (
	"RSSHub/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

func (h *CLIHandler) handleFetch() error {
	const op = "CLIHandler.handleFetch"
	log := h.log.GetSlogLogger().With(slog.String("op", op))

	if len(h.args) != 1 {
		log.Error("Invalid fetch command usage", "expected", "rsshub fetch", "got", h.args)
		return ErrInvFetchFlag
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Fetching feeds from all sources...")
	if err := h.aggregator.Start(ctx); err != nil {
		log.Error("Data fetch failed", "error", err)
		return err
	}
	return nil
}

func (h *CLIHandler) handleAdd() error {
	const op = "CLIHandler.handleAdd"
	log := h.log.GetSlogLogger().With(slog.String("op", op))

	if len(h.args) != 7 {
		log.Error("Invalid add command usage", "expected", "rsshub add --name <name> --url <url> --desc <description>", "got", h.args)
		return ErrInvAddFlag
	}

	if h.args[1] != nameSubFlag {
		log.Error("Missing required --name flag", "got", h.args[1])
		return ErrMissingNameFlag
	}
	if h.args[3] != urlSubFlag {
		log.Error("Missing required --url flag", "got", h.args[3])
		return ErrMissingUrlFlag
	}

	if h.args[5] != descriptionFlag {
		log.Error("Missing required --desc flag", "got", h.args[5])
		return ErrMissingDescFlag
	}

	name := h.args[2]
	url := h.args[4]
	desc := h.args[6]

	if len(name) == 0 {
		log.Error("Feed name cannot be empty")
		return ErrEmptyName
	}
	if len(url) == 0 {
		log.Error("Feed URL cannot be empty")
		return ErrEmptyUrl
	}
	if len(desc) == 0 {
		log.Error("Feed description cannot be empty")
		return ErrEmptyDesc
	}

	log.Info("Adding new feed", "name", name, "url", url, "description", desc)
	if err := h.aggregator.AddFeed(name, desc, url); err != nil {
		log.Error("Failed to add feed", "error", err)
		return err
	}

	msg := fmt.Sprintf("Feed %s added succesfully with URL %s", name, url)
	h.log.Notify(msg)
	return nil
}

func (h *CLIHandler) handleInterval() error {
	const op = "CLIHandler.handleInterval"
	log := h.log.GetSlogLogger().With(slog.String("op", op))

	if len(h.args) != 2 {
		log.Error("Invalid interval command usage", "expected", "rsshub set-interval <duration>", "got", h.args)
		return ErrInvIntervalFlag
	}

	interval, err := time.ParseDuration(h.args[1])
	if err != nil {
		log.Error("Invalid duration format", "input", h.args[1], "error", err)
		return err
	}

	if interval < time.Minute*2 {
		log.Error("Interval must be at least 2 min")
		return errors.New("invalid interval, must be at least 2 min")
	}

	if err := h.aggregator.SetInterval(interval); err != nil {
		log.Error("Failed to set fetch interval", "error", err)
		return err
	}

	return nil
}

func (h *CLIHandler) handleWorkers() error {
	const op = "CLIHandler.handleWorkers"
	log := h.log.GetSlogLogger().With(slog.String("op", op))

	if len(h.args) != 2 {
		log.Error("Invalid workers command usage", "expected", "rsshub set-workers <count>", "got", h.args)
		return ErrInvWorkersFlag
	}

	workerCount, err := strconv.Atoi(h.args[1])
	if err != nil {
		log.Error("Invalid worker count, must be an integer", "input", h.args[1], "error", err)
		return err
	}

	if workerCount < 1 {
		return errors.New("worker count must be greater than 0")
	}

	log.Info("Setting worker count", "count", workerCount)
	if err := h.aggregator.Resize(workerCount); err != nil {
		log.Error("Failed to set worker count", "error", err)
		return err
	}

	return nil
}

func (h *CLIHandler) handleDelete() error {
	const op = "CLIHandler.handleDelete"
	log := h.log.GetSlogLogger().With(slog.String("op", op))

	if len(h.args) != 3 {
		log.Error("Invalid delete command usage", "expected", "rsshub delete --name <name>", "got", h.args)
		return ErrInvDeleteFlag
	}

	if h.args[1] != nameSubFlag {
		log.Error("Missing required --name flag", "got", h.args[1])
		return ErrMissingNameFlag
	}

	name := h.args[2]
	if len(name) == 0 {
		log.Error("Feed name for deletion cannot be empty")
		return ErrEmptyName
	}

	log.Info("Deleting feed", "name", name)
	if err := h.aggregator.DeleteFeed(name); err != nil {
		log.Error("Failed to delete feed", "name", name, "error", err)
		return err
	}

	log.Info("Feed deleted successfully", "name", name)
	return nil
}

func (h *CLIHandler) handleList() error {
	const op = "CLIHandler.handleList"
	log := h.log.GetSlogLogger().With(
		slog.String("op", op),
	)

	var (
		feedCount int
		err       error
	)
	switch len(h.args) {
	case 1:
		feedCount = 0
	case 3:
		if h.args[1] != numSubFlag {
			log.Error("Missing --num flag", "got", h.args)
			return ErrMissingNumFlag
		}
		feedCount, err = strconv.Atoi(h.args[2])
		if err != nil {
			log.Error("Invalid feed count, must be an integer", "input", h.args[2], "error", err)
			return ErrMissingNumFlag
		}

		if feedCount < 1 {
			return ErrInvNumFlag
		}
	default:
		log.Error("Invalid list command usage", "expected", "rsshub list | rsshub list --num <num>", "got", h.args)
		return ErrInvListFlag
	}

	log.Info("Getting feeds list", "feed count", feedCount)
	feeds, err := h.aggregator.ListFeeds(feedCount)
	if err != nil {
		log.Error("Failed to get feeds list", "error", err)
		return err
	}

	utils.PrintFeedsList(feeds)
	return nil
}

func (h *CLIHandler) handleArticle() error {
	const op = "CLIHandler.handleArticle"
	log := h.log.GetSlogLogger().With(
		slog.String("op", op),
	)

	feedName, num := "", 0

	switch len(h.args) {
	case 3:
		if h.args[1] != feednameSubFlag {
			log.Error("Missing required --feed-name flag", "got", h.args)
			return ErrMissingFeedNameSubFlag
		}
		feedName = h.args[2]

		if len(feedName) == 0 {
			log.Error("Feed name flag cannot be empty")
			return ErrEmptyFeedName
		}
	case 5:
		var err error
		if h.args[1] != feednameSubFlag {
			log.Error("Missing required --feed-name flag", "got", h.args)
			return ErrMissingFeedNameSubFlag
		}
		feedName = h.args[2]

		if len(feedName) == 0 {
			log.Error("Feed name flag cannot be empty")
			return ErrEmptyFeedName
		}

		if h.args[3] != numSubFlag {
			log.Error("Missing --num flag", "got", h.args)
			return ErrMissingNumFlag
		}

		num, err = strconv.Atoi(h.args[4])
		if err != nil {
			log.Error("Invalid article count, must be an integer", "input", h.args[4], "error", err)
			return ErrMissingNumFlag
		}

		if num < 1 {
			return ErrInvNumFlag
		}

	default:
		log.Error(ErrArticleFlagExpected.Error(), "got", h.args)
		return ErrArticleFlagExpected
	}

	log.Info("Getting articles list", "feedName", feedName, "num", num)
	articles, err := h.aggregator.GetArticles(feedName, num)
	if err != nil {
		log.Error("Failes to get articles", "error", err)
		return err
	}

	utils.PrintArticleList(articles, feedName)
	return nil
}

func (h *CLIHandler) handleStatus() error {
	const op = "CLIHandler.handleStatus"
	log := h.log.GetSlogLogger().With(
		slog.String("op", op),
	)

	rssConfig, err := h.aggregator.GetConfig(context.Background())
	if err != nil {
		log.Error("failed to load config", "error", err)
		return errors.New("failed to load config")
	}

	msg := utils.PrettyRssConfig(rssConfig)

	h.log.Notify(msg)
	return nil
}
