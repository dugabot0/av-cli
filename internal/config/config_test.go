package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any env vars that might interfere
	os.Unsetenv("DMM_API_ID")
	os.Unsetenv("DMM_AFFILIATE_ID")
	os.Unsetenv("DUGA_APP_ID")
	os.Unsetenv("DUGA_AGENT_ID")
	os.Unsetenv("DUGA_BANNER_ID")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Default banner_id should be "01"
	if cfg.DUGA.BannerID != "01" {
		t.Errorf("expected DUGA.BannerID=%q, got %q", "01", cfg.DUGA.BannerID)
	}

	// Credentials should be empty without env/file
	if cfg.DMM.APIID != "" {
		t.Errorf("expected empty DMM.APIID, got %q", cfg.DMM.APIID)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("DMM_API_ID", "test-api-id")
	t.Setenv("DMM_AFFILIATE_ID", "test-affiliate-990")
	t.Setenv("DUGA_APP_ID", "test-duga-app")
	t.Setenv("DUGA_AGENT_ID", "test-agent")
	t.Setenv("DUGA_BANNER_ID", "42")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.DMM.APIID != "test-api-id" {
		t.Errorf("DMM.APIID: want %q, got %q", "test-api-id", cfg.DMM.APIID)
	}
	if cfg.DMM.AffiliateID != "test-affiliate-990" {
		t.Errorf("DMM.AffiliateID: want %q, got %q", "test-affiliate-990", cfg.DMM.AffiliateID)
	}
	if cfg.DUGA.AppID != "test-duga-app" {
		t.Errorf("DUGA.AppID: want %q, got %q", "test-duga-app", cfg.DUGA.AppID)
	}
	if cfg.DUGA.AgentID != "test-agent" {
		t.Errorf("DUGA.AgentID: want %q, got %q", "test-agent", cfg.DUGA.AgentID)
	}
	if cfg.DUGA.BannerID != "42" {
		t.Errorf("DUGA.BannerID: want %q, got %q", "42", cfg.DUGA.BannerID)
	}
}
