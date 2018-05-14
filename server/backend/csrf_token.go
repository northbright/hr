package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/northbright/uuid"
)

var (
	csrfTokenTimeout int = 60
)

func createCSRFToken(pool *redis.Pool) (string, error) {
	token, err := uuid.New()
	if err != nil {
		return "", err
	}

	conn := redisPool.Get()
	defer conn.Close()

	conn.Do("MULTI")

	k := fmt.Sprintf("csrf_token:%v", token)
	conn.Send("SET", k, "")
	conn.Send("EXPIRE", k, csrfTokenTimeout)

	if _, err = conn.Do("EXEC"); err != nil {
		return "", err
	}

	return token, nil
}

func validCSRFToken(pool *redis.Pool, token string) (bool, error) {
	conn := redisPool.Get()
	defer conn.Close()

	k := fmt.Sprintf("csrf_token:%v", token)
	exists, err := redis.Bool(conn.Do("EXISTS", k))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func getCSRFToken(c *gin.Context) {
	type Reply struct {
		ErrMsg    string `json:"err_msg"`
		CSRFToken string `json:"csrf_token"`
	}

	var (
		err        error
		statusCode int
		reply      Reply
	)

	defer func() {
		if err != nil {
			log.Printf("getCSRFToken() error: %v", err)
		}

		c.JSON(statusCode, reply)
	}()

	if reply.CSRFToken, err = createCSRFToken(redisPool); err != nil {
		statusCode = 500
		reply.ErrMsg = "failed to create CSRF token"
		return
	}

	statusCode = 200
}
