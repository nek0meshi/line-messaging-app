package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/hello", wrap(handleHello))
	r.POST("/webhook", wrap(handleWebhook))
}

type HandleFunc func(c *gin.Context) (int, string, gin.H)

func wrap(handleFunc HandleFunc) (func(c *gin.Context)) {
	return func(c *gin.Context) {
		statusCode, errorMessage, body := handleFunc(c)

		if statusCode >= 400 && errorMessage != "" {
			body = gin.H{
				"message": errorMessage,
			}
		}

		c.JSON(statusCode, body)
	}
}
