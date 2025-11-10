package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	dbConfig := Config{}

	filePath, err := getConfigFilePath()
	if err != nil {
		return dbConfig, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return dbConfig, err
	}

	if err := json.Unmarshal(data, &dbConfig); err != nil {
		return dbConfig, err
	}

	return dbConfig, nil
}

func getConfigFilePath() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error in getting home directory: %s", err)
	}

	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}

func write(cfg Config) error {

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0666); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return write(*cfg)
}
