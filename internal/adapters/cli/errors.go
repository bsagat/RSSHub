package cli

import "errors"

var (
	ErrInvFetchFlag           = errors.New("fetch flag is invalid")
	ErrInvAddFlag             = errors.New("add flag is invalid")
	ErrInvIntervalFlag        = errors.New("set-interval flag is invalid")
	ErrInvWorkersFlag         = errors.New("set-workers count flag is invalid")
	ErrInvDeleteFlag          = errors.New("delete flag is invalid")
	ErrInvListFlag            = errors.New("list flag is invalid")
	ErrInvArticlesFlag        = errors.New("articles flag is invalid")
	ErrMissingNameFlag        = errors.New("--name flag is required")
	ErrMissingUrlFlag         = errors.New("--url flag is required")
	ErrMissingNumFlag         = errors.New("--num flag is missing")
	ErrMissingFeedNameSubFlag = errors.New("--feed-name flag is missing")
	ErrEmptyFeedName          = errors.New("--feed-name is required")
	ErrEmptyName              = errors.New("--name flag is required")
	ErrEmptyUrl               = errors.New("--url flag is required")
)

var (
	ErrStatusCode = 1
	OkStatusCode  = 0
)
