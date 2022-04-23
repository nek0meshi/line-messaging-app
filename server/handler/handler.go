package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/hello", handleHello)
	r.POST("/webhook", handleWebhook)
}
