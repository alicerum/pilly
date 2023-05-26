package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Bot Bot `json:"bot"`
	DB  DB  `json:"db"`
	Log Log `json:"log"`
}

type Log struct {
	Level string `json:"level"`
}

type DB struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

type Bot struct {
	Token string `json:"token"`
}

func ReadConfig(configPath string) (*Config, error) {
	cfgBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshald config: %w", err)
	}

	return &cfg, nil
}

func (l *Log) SetLogger() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(l.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	return nil
}

func (d *DB) CreateDB(ctx context.Context) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     d.Host,
		User:     d.Username,
		Password: d.Password,
		Database: d.Database,
	})

	return db, db.Ping(ctx)
}
