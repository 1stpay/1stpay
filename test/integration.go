package test

import (
	"context"
	"fmt"
	"testing"

	"log"

	"github.com/1stpay/1stpay/internal/config"
	route "github.com/1stpay/1stpay/internal/transport/rest"
	"github.com/1stpay/1stpay/test/factory"
	"github.com/gin-gonic/gin"
)

type IntegrationTest struct {
	Database    *TestDatabase
	Context     context.Context
	Env         *config.Env
	GinEngine   *gin.Engine
	TestFactory *factory.TestFactory
}

func NewIntegrationTest(t *testing.T, rootPath string) *IntegrationTest {
	ctx := context.Background()
	envPath := fmt.Sprintf("%v.env", rootPath)
	env := config.NewEnv(envPath)
	database, err := NewTestPostgresDatabase(ctx, env, rootPath)
	if err != nil {
		log.Fatalf("Eror during test DB setup, %v", err)
	}
	ginEngine := gin.Default()
	gin.SetMode(gin.TestMode)
	route.SetupRoutes(env, database.GormDB, ginEngine)
	t.Cleanup(func() {
		if err := database.Cleanup(ctx); err != nil {
			t.Fatalf("Error during test DB cleanup: %v", err)
		}
	})

	return &IntegrationTest{
		Database:    database,
		Context:     ctx,
		Env:         env,
		GinEngine:   ginEngine,
		TestFactory: factory.NewTestFactory(database.GormDB),
	}
}
