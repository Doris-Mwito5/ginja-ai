package configs

import (
	"log"

	"github.com/spf13/viper"
)

var Config *config

type config struct {
	Debug       bool   `mapstructure:"DEBUG"`
	Port        string `mapstructure:"HTTP_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func InitializeEnvironment() {
	// Load .env from current directory so DATABASE_URL, etc. are read
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/configs")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: .env file not found (%v), using system environment variables", err)
	}

	var cnf config

	err = viper.Unmarshal(&cnf)
	if err != nil {
		log.Fatal(err.Error())
	}

	Config = &cnf
}

func IsEnvProduction() bool {
	return Config.Environment == "production"
}
