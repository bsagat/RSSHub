package app

import "os"

func (a *App) Run() {
	code := a.cliHandler.ParseFlags()
	a.CleanUp()
	os.Exit(code)
}

func (a *App) CleanUp() {
	a.log.Info("Started cleaning up...")
	a.postrgresRepo.Db.Close()
	a.cliHandler.Close()
	a.log.Info("Cleaning up finished")
}
