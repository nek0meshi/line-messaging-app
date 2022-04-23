package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"nek0meshi/line-messaging-app/handler"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	if readDotEnv() != nil {
		log.Print("failed to read dotenv file")
	}

	handler.SetupRouter(r)

	return r
}

func readDotEnv() error {
	return godotenv.Load()
}

func main() {
	if len(os.Args) < 2 {
		r := setupRouter()

		r.Run(":80")
	}
}
