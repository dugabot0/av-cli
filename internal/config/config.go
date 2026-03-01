package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	DMM  DMMConfig  `mapstructure:"dmm"`
	DUGA DUGAConfig `mapstructure:"duga"`
}

type DMMConfig struct {
	APIID       string `mapstructure:"api_id"`
	AffiliateID string `mapstructure:"affiliate_id"`
}

type DUGAConfig struct {
	AppID    string `mapstructure:"app_id"`
	AgentID  string `mapstructure:"agent_id"`
	BannerID string `mapstructure:"banner_id"`
}

// Load reads config from file and environment variables via viper.
// Priority: env vars > config file > defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("duga.banner_id", "01")

	// Config file
	home, err := os.UserHomeDir()
	if err == nil {
		cfgPath := filepath.Join(home, ".config", "av-cli")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(cfgPath)
		// Ignore file-not-found errors; env vars may suffice
		_ = v.ReadInConfig()
	}

	// Environment variable bindings
	v.SetEnvPrefix("")
	_ = v.BindEnv("dmm.api_id", "DMM_API_ID")
	_ = v.BindEnv("dmm.affiliate_id", "DMM_AFFILIATE_ID")
	_ = v.BindEnv("duga.app_id", "DUGA_APP_ID")
	_ = v.BindEnv("duga.agent_id", "DUGA_AGENT_ID")
	_ = v.BindEnv("duga.banner_id", "DUGA_BANNER_ID")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
