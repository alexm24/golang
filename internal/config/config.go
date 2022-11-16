package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/alexm24/golang/internal/models"
)

func getConfigPath(path string) string {
	flag.StringVar(&path, "c", path, "set path to string")
	flag.Parse()
	return path
}

func ParseConfig(path string) (*models.Config, error) {
	var cfg models.Config
	configPath := getConfigPath(path)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config file %w", err)
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file from %s as yaml %w", configPath, err)
	}

	return &cfg, nil
}
