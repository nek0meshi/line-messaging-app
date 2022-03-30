package main

import (
	"os"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	return r
}

func main() {
	if len(os.Args) < 2 {
		r := setupRouter()

		r.Run(":80")
	}
}
