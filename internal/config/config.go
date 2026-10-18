package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPAddr   string `yaml:"http_addr"`
	SSHAddr    string `yaml:"ssh_addr"`
	DataDir    string `yaml:"data_dir"`
	DBPath     string `yaml:"db_path"`
	ReposRoot  string `yaml:"repos_root"`
	AvatarsDir string `yaml:"avatars_dir"`
	SSHDir     string `yaml:"ssh_dir"`
}

func Default() Config {
	return Config{
		HTTPAddr:   ":8080",
		SSHAddr:    ":2222",
		DataDir:    "data",
		DBPath:     "data/db.sqlite",
		ReposRoot:  "data/repos",
		AvatarsDir: "data/avatars",
		SSHDir:     "data/ssh",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		path = envOr("GITWOLD_CONFIG", "config.yaml")
	}

	if _, err := os.Stat(path); err == nil {
		raw, err := os.ReadFile(path)
		if err != nil {
			return cfg, fmt.Errorf("read config: %w", err)
		}

		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return cfg, fmt.Errorf("parse config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return cfg, err
	}

	cfg.HTTPAddr = envOr("GITWOLD_HTTP_ADDR", cfg.HTTPAddr)
	cfg.SSHAddr = envOr("GITWOLD_SSH_ADDR", cfg.SSHAddr)
	cfg.DataDir = envOr("GITWOLD_DATA_DIR", cfg.DataDir)
	cfg.DBPath = envOr("GITWOLD_DB_PATH", cfg.DBPath)
	cfg.ReposRoot = envOr("GITWOLD_REPOS_ROOT", cfg.ReposRoot)
	cfg.AvatarsDir = envOr("GITWOLD_AVATARS_DIR", cfg.AvatarsDir)
	cfg.SSHDir = envOr("GITWOLD_SSH_DIR", cfg.SSHDir)

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}

	if cfg.SSHAddr == "" {
		cfg.SSHAddr = ":2222"
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}

	return fallback
}
