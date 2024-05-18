package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
}

type DBConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"db_name"`
	MigrationsPath string `yaml:"migrations_path"`
}

func (cfg DBConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func (cfg ServerConfig) String() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Get() (Config, error) {
	fileName := "config.yaml"
	cfg := Config{}
	f, err := os.Open("config.yaml")
	if err != nil {
		return cfg, fmt.Errorf("error open \"%s\": %w", fileName, err)
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("error decode: %w", err)
	}
	return cfg, nil
}
