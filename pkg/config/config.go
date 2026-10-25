package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Application struct {
		Message string `yaml:"message"`
	} `yaml:"application"`
}

func Load(configDir string) (*Config, error) {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "production"
	}

	var configuration Config

	configFile := fmt.Sprintf("%s/config.%s.yaml", configDir, environment)
	if err := cleanenv.ReadConfig(configFile, &configuration); err != nil {
		cleanenv.ReadEnv(&configuration)
	}

	localFile := fmt.Sprintf("%s/config.local.yaml", configDir)
	if _, err := os.Stat(localFile); err == nil {
		if err := cleanenv.ReadConfig(localFile, &configuration); err != nil {
			return nil, err
		}
	}

	return &configuration, nil
}
