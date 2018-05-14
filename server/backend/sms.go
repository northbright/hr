package main

/*
import (
	"fmt"

	"github.com/northbright/aliyun/message"
)

func sendValidationSMS(c *gin.Context) {
	type Req struct {
		MobilePhoneNum string `json:"mobile_phone_num"`
		CSRFToken      string `json:"csrf_token"`
	}

	type Reply struct {
		ErrMsg string `json:"err_msg"`
	}

	var (
		err        error
		statusCode int
		r          Req
		reply      Reply
	)

	defer func() {
		if err != nil {
			log.Printf("sendValidationSMS() error: %v", err)
		}

		c.JSON(statusCode, reply)
	}()

	if err = c.BindJSON(&r); err != nil {
		statusCode = 400
		reply.ErrMsg = "invalid request"
		return
	}

	// Validate CSRF token
	valid, err := validCSRFToken(redisPool, r.CSRFToken)
	if err != nil {
		statusCode = 500
		reply.ErrMsg = "failed to validate CSRF token"
		return
	}

	if !valid {
		statusCode = 400
		reply.ErrMsg = "invalid request"
		return
	}

	// Check if mobile phone number belongs to someone of employees.
}
*/
