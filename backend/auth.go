package main

import (
	"fmt"
	"log"
	"net/http"
	//"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/northbright/jwthelper"
)

func getLoginID(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie("jwt")
	switch err {
	case http.ErrNoCookie:
		return "", fmt.Errorf("no jwt found in cookie")
	case nil:
		tokenString := cookie.Value
		parser, err := jwthelper.NewParser(jwt.SigningMethodRS256, []byte(rsaPubPEM))
		if err != nil {
			return "", err
		}

		m, err := parser.Parse(tokenString)
		if err != nil {
			return "", fmt.Errorf("parser.Parse() error: %v", err)
		}

		// Convert interface{} to string
		ID, ok := m["id"].(string)
		if !ok {
			return "", fmt.Errorf("failed to convert interface{} to string")
		}
		return ID, nil
	default:
		return "", fmt.Errorf("get JWT cookie error: %v", err)
	}
}

func validLogin(username, password string) bool {
	return true
}

func postLoginHandler(c *gin.Context) {
	type Req struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		CSRFToken string `json:"csrf_token"`
	}

	type Reply struct {
		ErrMsg string `json:"err_msg"`
		ID     string `json:"id"`
	}

	var (
		err        error
		statusCode int
		r          Req
		reply      Reply
	)

	defer func() {
		if err != nil {
			log.Printf("login() error: %v", err)
		}

		c.JSON(statusCode, reply)
	}()

	if err = c.BindJSON(&r); err != nil {
		statusCode = 400
		reply.ErrMsg = "invalid request"
		return
	}

	valid, err := validCSRFToken(redisPool, r.CSRFToken)
	if err != nil {
		statusCode = 500
		reply.ErrMsg = "failed to validate CSRF token"
		return
	}
	if !valid {
		statusCode = 400
		reply.ErrMsg = "invalid CSRF token"
		return
	}

	if !validLogin(r.Username, r.Password) {
		statusCode = 401
		reply.ErrMsg = "incorrect password"
		return
	}

	reply.ID = "1"

	signer, err := jwthelper.NewSigner(jwt.SigningMethodRS256, []byte(rsaPrivPEM))
	if err != nil {
		statusCode = 500
		reply.ErrMsg = "failed to create token signer"
		return
	}

	tokenString, err := signer.SignedString(
		jwthelper.NewClaim("id", reply.ID),
	)
	if err != nil {
		statusCode = 500
		reply.ErrMsg = "failed to create token string"
		return
	}
	cookie := jwthelper.NewCookie(tokenString)
	http.SetCookie(c.Writer, cookie)

	statusCode = 200
}
