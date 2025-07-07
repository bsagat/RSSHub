package app

import "RSSHub/internal/adapters/cli"

func (a *App) Run() {
	cli.ParseFlags()
}
