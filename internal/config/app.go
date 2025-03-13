package config

import (
	"gorm.io/gorm"
)

type Application struct {
	Env      *Env
	Postgres *gorm.DB
	Deps     *Dependencies
}

func App() *Application {
	env := NewEnv()
	db := NewPostgresDatabase(env)
	deps := NewDependencies(db, env)
	return &Application{
		Env:      env,
		Postgres: db,
		Deps:     deps,
	}
}
