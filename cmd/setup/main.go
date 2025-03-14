package main

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/setup"
)

func main() {
	app := config.App()
	setup := setup.NewSetup(app.Postgres, app.Deps)
	setup.Init()
}
