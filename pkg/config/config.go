package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	EnvironmentDev   = "dev"
	EnvironmentStage = "stage"
	EnvironmentProd  = "prod"
	EnvironmentTest  = "test"
)

// CI trigger

type Config struct {
	Environment      string `envconfig:"ENVIRONMENT" default:"dev"`
	GRPCPort         int    `envconfig:"GRPC_PORT" default:"50051"`
	DataDir          string `envconfig:"ASSETS_DIR" default:"data"`
	GoogleMapsAPIKey string `envconfig:"GOOGLE_MAPS_API_KEY" required:"true"`
	Log              struct {
		Formater               string `envconfig:"LOG_FORMATER" default:"text"`
		Level                  string `envconfig:"LOG_LEVEL" default:"info"`
		EnableDetailedResponse bool   `envconfig:"LOG_ENABLE_DETAILED_RESPONSE" default:"true"`
		EnableDetailedRequest  bool   `envconfig:"LOG_ENABLE_DETAILED_REQUEST" default:"true"`
	}
}

func NewParsedConfig() (Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	_ = godotenv.Load(".env." + env + ".local")

	if env != EnvironmentTest {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load() // The Original .env

	cnf := Config{}
	if err := envconfig.Process("", &cnf); err != nil {
		// example: the environment variable "MODULR_WEBHOOK_SECRET" is missing
		return Config{}, err
	}
	return cnf, nil
}
