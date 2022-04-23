package handler

import (
	"github.com/gin-gonic/gin"

)
func handleHello(c *gin.Context) (int, string, gin.H) {
	return 200, "", gin.H{
		"message": "Hello World!",
	}
}
