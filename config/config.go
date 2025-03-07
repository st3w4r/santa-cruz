package config

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB_MANAGER_PATH string
}

func LoadConfig() (Config, error) {
	godotenv.Load()
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	homeDataDir := xdg.DataHome

	dbPath := filepath.Join(homeDataDir, "/santa-cruz/dbmanager.db")

	if cfg.DB_MANAGER_PATH == "" {
		cfg.DB_MANAGER_PATH = dbPath
	}

	fmt.Println(cfg)

	return cfg, nil
}
