package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	appPort              string
	accessTokenKey       string
	accessTokenLifetime  string
	refreshTokenKey      string
	refreshTokenLifetime string
	cacheSize            string
	dbHost               string
	dbPort               string
	dbName               string
	dbUser               string
	dbPassword           string
	sslMode              string
}

func New() *Config {
	return &Config{
		appPort:              mustGetEnv("PORT"),
		accessTokenKey:       mustGetEnv("ACCESS_TOKEN_KEY"),
		accessTokenLifetime:  mustGetEnv("ACCESS_TOKEN_LIFETIME"),
		cacheSize:            mustGetEnv("CACHE_SIZE"),
		refreshTokenKey:      mustGetEnv("REFRESH_TOKEN_KEY"),
		refreshTokenLifetime: mustGetEnv("REFRESH_TOKEN_LIFETIME"),
		dbHost:               mustGetEnv("DB_HOST"),
		dbPort:               mustGetEnv("DB_PORT"),
		dbName:               mustGetEnv("DB_NAME"),
		dbUser:               mustGetEnv("DB_USERNAME"),
		dbPassword:           mustGetEnv("DB_PASSWORD"),
		sslMode:              mustGetEnv("SSL_MODE"),
	}
}

func NewTest() *Config {
	os.Setenv("PORT", "8080")
	os.Setenv("ACCESS_TOKEN_KEY", "test")
	os.Setenv("ACCESS_TOKEN_LIFETIME", "1200")
	os.Setenv("REFRESH_TOKEN_KEY", "test")
	os.Setenv("REFRESH_TOKEN_LIFETIME", "604800")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "recovy")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "admin123")
	os.Setenv("CACHE_SIZE", "1") // MB
	os.Setenv("SSL_MODE", "disable")

	return New()
}

func (c *Config) AppPort() string {
	return fmt.Sprintf(":%s", c.appPort)
}

func (c *Config) AccessTokenKey() string {
	return c.accessTokenKey
}

func (c *Config) AccessTokenLifetime() int {
	res, err := strconv.Atoi(c.accessTokenLifetime)
	if err != nil {
		panic("Access token lifetime must be filled with a number greater than 0")
	}

	return res
}

func (c *Config) CacheSize() int {
	res, err := strconv.Atoi(c.cacheSize)
	if err != nil {
		panic("Cache-size must be filled with a number greater than 0")
	}

	return res
}

func (c *Config) RefreshTokenKey() string {
	return c.refreshTokenKey
}

func (c *Config) RefreshTokenLifetime() int {
	res, err := strconv.Atoi(c.refreshTokenLifetime)
	if err != nil {
		panic("Refresh token lifetime must be filled with a number greater than 0")
	}

	return res
}

func (c *Config) DatabaseConnection() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.dbHost, c.dbPort, c.dbUser, c.dbPassword, c.dbName, c.sslMode,
	)
}

func mustGetEnv(key string) string {
	res := os.Getenv(key)
	if res == "" {
		panic("key " + key + " is empty")
	}

	return res
}
