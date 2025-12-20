package config

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddr    = "127.0.0.1:8020"
	syncIntervalDefault  = "5m"
	logLevelDefault      = "INFO"
	secretKeyDefault     = "SecretKey"
	tokenDurationDefault = "2h"
)

// NewServerConfig возвращает конфиг для сервера.
func NewServerConfig() (*Config, error) {
	newConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	logLevel := serverFlagSet.String("w", logLevelDefault, "log level")
	databaseDsn := serverFlagSet.String("d", "", "Database DSN")
	key := serverFlagSet.String("k", secretKeyDefault, "sha key")
	cyptoKey := serverFlagSet.String("crypto-key", "", "CRYPTO KEY")
	configFilePath := serverFlagSet.String("c", "", "config file")
	cert := serverFlagSet.String("cert", "", "certifacate")
	err = serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if newConfig.ServerAddr == nil {
		newConfig.ServerAddr = serverAddr
	}
	if newConfig.LogLevel == nil {
		newConfig.LogLevel = logLevel
	}

	if newConfig.DatabaseDsn == nil {
		newConfig.DatabaseDsn = databaseDsn
	}
	if newConfig.Key == nil {
		newConfig.Key = key
	}

	if newConfig.CryptoKey == nil {
		newConfig.CryptoKey = cyptoKey
	}
	if newConfig.ConfigFilePath == nil {
		newConfig.ConfigFilePath = configFilePath
	}
	if newConfig.Cert == nil {
		newConfig.Cert = cert
	}

	if *newConfig.ConfigFilePath != "" {
		err = newConfig.UpdateFromConfig()
		if err != nil {
			return newConfig, err
		}
	}
	return newConfig, nil
}

// NewAgentConfig возвращает конфиг для агента.
func NewAgentConfig() (*Config, error) {
	newConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	syncInterval := agentFlagSet.String("p", syncIntervalDefault, "input pollInterval")
	key := agentFlagSet.String("k", secretKeyDefault, "sha key")
	cyptoKey := agentFlagSet.String("crypto-key", "", "CRYPTO KEY")
	configFilePath := agentFlagSet.String("c", "", "config file")
	tokenDurationString := agentFlagSet.String("tt", tokenDurationDefault, "jwn token duration")
	err = agentFlagSet.Parse(os.Args[1:])
	if err != nil {
		return newConfig, err
	}
	if newConfig.ServerAddr == nil {
		newConfig.ServerAddr = serverAddr
	}

	if newConfig.SyncInterval == nil {
		pollIntervalDuration, err := time.ParseDuration(*syncInterval)
		if err != nil {
			return newConfig, err
		}
		newConfig.SyncInterval = &timeConfig{Duration: pollIntervalDuration}
	}
	if newConfig.Key == nil {
		newConfig.Key = key
	}

	if newConfig.CryptoKey == nil {
		newConfig.CryptoKey = cyptoKey
	}
	if newConfig.ConfigFilePath == nil {
		newConfig.ConfigFilePath = configFilePath
	}
	if newConfig.TokenDuration == nil {
		tokenDuration, err := time.ParseDuration(*tokenDurationString)
		if err != nil {
			return newConfig, err
		}
		newConfig.TokenDuration = &timeConfig{Duration: tokenDuration}
	}

	if *newConfig.ConfigFilePath != "" {
		err = newConfig.UpdateFromConfig()
		if err != nil {
			return newConfig, err
		}
	}

	return newConfig, nil
}

// InitConfig иницилизирует конфиг
func InitConfig() (*Config, error) {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		return nil, err
	}
	return &newConfig, err
}

// Config хранит конфиг
type Config struct {
	ServerAddr     *string     `env:"ADDRESS" json:"address"`
	SyncInterval   *timeConfig `env:"SYNC_INTERVAL" json:"sync_interval"`
	LogLevel       *string     `env:"LOG_LEVEL"`
	DatabaseDsn    *string     `env:"DATABASE_DSN" json:"database_dsn"`
	Key            *string     `env:"KEY"`
	Wait           *timeConfig `env:"WAIT"`
	CryptoKey      *string     `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFilePath *string     `env:"CONFIG"`
	Cert           *string     `env:"CERT" json:"cert"`
	TokenDuration  *timeConfig `env:"TOKEN_DURATION" json:"token_duration"`
}

type timeConfig struct {
	time.Duration
}

func (c Config) GetServeAddress() string {
	return *c.ServerAddr
}

func (c *Config) UpdateFromConfig() error {
	fileBytes, err := os.ReadFile(*c.ConfigFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileBytes, &c)
}

func (t *timeConfig) UnmarshalJSON(b []byte) error {
	return t.parseDuration(b)
}

func (t *timeConfig) UnmarshalText(text []byte) error {
	return t.parseDuration(text)
}

func (t *timeConfig) parseDuration(data []byte) error {
	s := strings.Trim(string(data), "\"")
	duration, err := time.ParseDuration(string(s))
	if err != nil {
		return err
	}
	t.Duration = duration
	return nil
}

func InitDefaultEnv() error {
	envDefaults := map[string]string{
		"ADDRESS":        defaultServerAddr,
		"LOG_LEVEL":      logLevelDefault,
		"KEY":            secretKeyDefault,
		"TOKEN_DURATION": tokenDurationDefault,
	}
	for k, v := range envDefaults {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
