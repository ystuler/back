package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	IP   string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	SigningKey string
}

var JWTSigningKey string

func NewConfig() *Config {
	viper.SetDefault("server.ip", "localhost")
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "root")
	viper.SetDefault("database.dbname", "memoryCards")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "Europe/Moscow")
	viper.SetDefault("JWT.signingkey", "SecretKey")

	// Для установки переменных окружения необходимо использовать такой стиль, как `APP_SERVER_PORT="8000"`
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if config.JWT.SigningKey == "" {
		config.JWT.SigningKey = generateRandomKey(32) // 32 байта = 256 бит
		log.Printf("Generated new signing key: %s", config.JWT.SigningKey)
	}

	JWTSigningKey = config.JWT.SigningKey

	return &config
}

func generateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Unable to generate random key: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func (dbConfig *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode, dbConfig.Timezone)
}

func (srvConfig *ServerConfig) GetADDR() string {
	return fmt.Sprintf("%s:%s", srvConfig.IP, srvConfig.Port)
}
