package config

import (
	"encoding/json"
	"os"
)

const ConfigPath = "./.gatorconfig.json"

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFile, err := os.ReadFile(ConfigPath)
	if err != nil {
		return Config{}, err
	}

	var configuration Config
	if err := json.Unmarshal(configFile, &configuration); err != nil {
		return Config{}, err
	}
	return configuration, nil
}

func SetUser(config *Config, username string) error {
	config.CurrentUser = username

	encoded, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(ConfigPath, encoded, 0644); err != nil {
		return err
	}
	return nil
}