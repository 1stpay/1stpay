package app

import (
	"time"

	"github.com/gin-contrib/cors"
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
	gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	route.SetupRoutes(env, db, gin, deps)
	gin.Run(":" + env.HttpPort)
}
