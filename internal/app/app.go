package app

import (
	"github.com/gin-gonic/gin"

	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/transport/rest/route"
)

func Run() {
	app := config.App()
	env := app.Env
	db := app.Postgres
	gin := gin.Default()
	route.Setup(env, db, gin)
	gin.Run(":" + env.HttpPort)
}
