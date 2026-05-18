package config

import (
	"os"
	"testing"
)

func TestNewConfigManager(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
database:
  host: testhost
  port: "5432"
  user: testuser
  password: testpass
  dbname: testdb

server:
  port: "9000"

jwt:
  secret: test-secret
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cm, err := NewConfigManager(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	cfg := cm.GetConfig()
	if cfg.Database.Host != "testhost" {
		t.Errorf("Expected database host 'testhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Server.Port != "9000" {
		t.Errorf("Expected server port '9000', got '%s'", cfg.Server.Port)
	}
	if cfg.JWT.Secret != "test-secret" {
		t.Errorf("Expected JWT secret 'test-secret', got '%s'", cfg.JWT.Secret)
	}
}

func TestConfigManagerDefaults(t *testing.T) {
	cm, err := NewConfigManager("nonexistent.yaml")
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	cfg := cm.GetConfig()
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected default database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("Expected default database port '5432', got '%s'", cfg.Database.Port)
	}
	if cfg.JWT.Secret == "" {
		t.Error("Expected auto-generated JWT secret, got empty string")
	}
}

func TestConfigManagerEnvOverride(t *testing.T) {
	os.Setenv("BOOK_DATABASE_HOST", "envhost")
	os.Setenv("BOOK_SERVER_PORT", "7000")
	defer func() {
		os.Unsetenv("BOOK_DATABASE_HOST")
		os.Unsetenv("BOOK_SERVER_PORT")
	}()

	cm, err := NewConfigManager("nonexistent.yaml")
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	cfg := cm.GetConfig()
	if cfg.Database.Host != "envhost" {
		t.Errorf("Expected env override 'envhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Server.Port != "7000" {
		t.Errorf("Expected env override '7000', got '%s'", cfg.Server.Port)
	}
}