package main

import (
	"fmt"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

var (
	currentDir, configFile    string
	templatesPath, staticPath string
	config                    Config
	redisPool                 *redis.Pool
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}()

	if err = loadConfig(configFile, &config); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	// Create redis pool.
	redisPool = redishelper.NewRedisPool(
		config.Redis.Addr,
		config.Redis.Password,
		config.Redis.PoolMaxActive,
		config.Redis.PoolMaxIdle,
		config.Redis.PoolIdleTimeout,
		config.Redis.PoolWait,
	)
	defer redisPool.Close()

	r := gin.Default()

	// Core APIs.
	r.POST("/api/login", login)
	r.GET("/api/csrf_token", getCSRFToken)

	r.Run(config.HTTPServer.Addr)
}

// init initializes path variables.
func init() {
	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")
	templatesPath = path.Join(currentDir, "templates")
	staticPath = path.Join(currentDir, "static")
}
