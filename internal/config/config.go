package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type Config struct {
DBHost string
DBPort string
DBUser string
DBPassword string
DBName string
Port string
FraudAPIURL string
}


func Load() (*Config, error) {
cfg := &Config{
DBHost: getEnv("DB_HOST", "localhost"),
DBPort: getEnv("DB_PORT", "3306"),
DBUser: getEnv("DB_USER", "root"),
DBPassword: getEnv("DB_PASSWORD", "secret"),
DBName: getEnv("DB_NAME", "compliance"),
Port: getEnv("PORT", "8080"),
FraudAPIURL: getEnv("FRAUD_API_URL", ""),
}
return cfg, nil
}


func getEnv(key, fallback string) string {
if v := os.Getenv(key); v != "" {
return v
}
return fallback
}


func NewGormDB(cfg *Config) (*gorm.DB, error) {
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}