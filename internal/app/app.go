package app

import (
	"github.com/gin-gonic/gin"

	"github.com/1stpay/1stpay/internal/config"
	route "github.com/1stpay/1stpay/internal/transport/rest"
)

func Run() {
	app := config.App()
	env := app.Env
	db := app.Postgres
	deps := app.Deps
	gin := gin.Default()
	route.SetupRoutes(env, db, gin, deps)
	gin.Run(":" + env.HttpPort)
}
