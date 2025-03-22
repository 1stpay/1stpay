package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv          string `mapstructure:"APP_ENV"`
	HttpPort        string `mapstructure:"HTTP_PORT"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPass          string `mapstructure:"DB_PASS"`
	DBName          string `mapstructure:"DB_NAME"`
	DBUrl           string `mapstructure:"DB_URL"`
	JwtSecret       string `mapstructure:"JWT_SECRET"`
	RpcUrl          string `mapstructure:"RPC_URL"`           // @TODO REMOVE
	PriceServiceKey string `mapstructure:"PRICE_SERVICE_KEY"` // @TODO REMOVE
}

func NewEnv(envPath ...string) *Env {
	var configFile string
	if len(envPath) > 0 && envPath[0] != "" {
		configFile = envPath[0]
	} else {
		configFile = ".env"
	}
	env := Env{}
	dir, _ := os.Getwd()
	fmt.Println(dir)
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}
	return &env
}
