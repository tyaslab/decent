package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type ConfigManager struct {
	viper  *viper.Viper
	config Config
}

func NewConfigManager(configPath string) (*ConfigManager, error) {
	v := viper.New()

	v.SetConfigFile(configPath)

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.dbname", "bookdb")
	v.SetDefault("server.port", "8080")

	v.SetEnvPrefix("BOOK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.JWT.Secret == "" {
		config.JWT.Secret = generateRandomSecret()
	}

	return &ConfigManager{
		viper:  v,
		config: config,
	}, nil
}

func (cm *ConfigManager) GetConfig() Config {
	return cm.config
}

func (cm *ConfigManager) GetString(key string) string {
	return cm.viper.GetString(key)
}

func (cm *ConfigManager) GetBool(key string) bool {
	return cm.viper.GetBool(key)
}

func (cm *ConfigManager) GetInt(key string) int {
	return cm.viper.GetInt(key)
}

func (cm *ConfigManager) GetInt32(key string) int32 {
	return cm.viper.GetInt32(key)
}

func generateRandomSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "default-secret-key-change-in-production"
	}
	return hex.EncodeToString(b)
}