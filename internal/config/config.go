package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return Config{}, err
	}

	gatorPath := homeDir + "/.gatorconfig.json"
	file, err := os.ReadFile(gatorPath)

	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(file, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUser = user
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	gatorPath := homeDir + "/.gatorconfig.json"
	config, err := json.Marshal(c)

	if err != nil {
		return err
	}

	err = os.WriteFile(gatorPath, config, 0644)
	if err != nil {
		return err
	}

	return nil
}
