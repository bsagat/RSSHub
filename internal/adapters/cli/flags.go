package cli

import "flag"

var (
	fetchFlag       = "fetch"
	addFlag         = "add"
	setIntervalFlag = "set-interval"
	setWorkerFlag   = "set-workers"
	listFlag        = "list"
	deleteFlag      = "delete"
	articlesFlag    = "articles"
)

var (
	nameSubFlag     = "--name"
	feednameSubFlag = "--feed-name"
	numSubFlag      = "--num"
	urlSubFlag      = "--url"
	descriptionFlag = "--desc"
)

var helpFlag = flag.Bool("help", false, "Prints help message")
