package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// HTTPServerConfig represents the HTTP server settings.
type HTTPServerConfig struct {
	Addr string `json:"addr"`
}

// RedisConfig represents the redis settings.
type RedisConfig struct {
	Addr            string `json:"addr"`
	Password        string `json:"password"`
	PoolMaxActive   int    `json:"pool_max_active"`
	PoolMaxIdle     int    `json:"pool_max_idle"`
	PoolIdleTimeout int    `json:"pool_max_timeout"`
	PoolWait        bool   `json:"pool_wait"`
}

// PostgreSQLConfig represents the PostgreSQL settings.
type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// JWTConfig contains the signing / verifying key for JWT.
type JWTConfig struct {
	KID              string `json:"kid"`
	Alg              string `json:"alg"`
	SigningKeyFile   string `json:"signing_key_file"`
	VerifyingKeyFile string `json:"verifying_key_file"`
}

// Config represents the app settings.
type Config struct {
	HTTPServer HTTPServerConfig `json:"http_server"`
	Redis      RedisConfig      `json:"redis"`
	PostgreSQL PostgreSQLConfig `json:"postgresql"`
	JWT        JWTConfig        `json:"jwt"`
}

// loadConfig loads app config.
func loadConfig(configFile string, config *Config) error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}
