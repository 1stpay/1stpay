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
	Repos       *config.Repos
	Usecases    *config.Usecases
	Deps        *config.Dependencies
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
	deps := config.NewDependencies(database.GormDB, env)
	route.SetupRoutes(env, database.GormDB, ginEngine, deps)
	t.Cleanup(func() {
		if err := database.Cleanup(ctx); err != nil {
			t.Fatalf("Error during test DB cleanup: %v", err)
		}
	})
	testFactory := factory.NewTestFactory(database.GormDB, deps)
	return &IntegrationTest{
		Database:    database,
		Context:     ctx,
		Env:         env,
		GinEngine:   ginEngine,
		TestFactory: testFactory,
		Repos:       deps.Repos,
		Deps:        deps,
	}
}
