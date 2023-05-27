package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	localPath = "./"
	homePath  = "~/.config/pilly"
	etcPath   = "/etc/pilly"
)

type Config struct {
	Bot Bot `mapstructure:"bot"`
	DB  DB  `mapstructure:"db"`
	Log Log `mapstructure:"log"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type DB struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
}

type Bot struct {
	Token string `mapstructure:"token"`
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

func ReadConfig(configPath *string) (*Config, error) {
	godotenv.Load(".env")
	if configPath != nil {
		viper.SetConfigFile(*configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(localPath)
		viper.AddConfigPath(homePath)
		viper.AddConfigPath(etcPath)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshald config: %w", err)
	}

	return &cfg, nil
}
