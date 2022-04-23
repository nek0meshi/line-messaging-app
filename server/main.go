package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Message struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Text string `json:"text"`
}

type Event struct {
	Type       string  `json:"type"`
	Message    Message `json:"message"`
	ReplyToken string  `json:"replyToken"`
}

type WebhookEvent struct {
	Events []Event `json:"events"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	if readDotEnv() != nil {
		log.Print("failed to read dotenv file")
	}

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.POST("/webhook", func(c *gin.Context) {
		// check signature
		// https://developers.line.biz/ja/reference/messaging-api/#signature-validation
		defer c.Request.Body.Close()

		b := &bytes.Buffer{}
		_, err := io.Copy(b, c.Request.Body)
		bodyCopy := bytes.NewReader(b.Bytes())
		c.Request.Body = io.NopCloser(bodyCopy)
		bodyCopy.Seek(0, 0)
		body := b.Bytes()

		// body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Print("failed to read body")
			c.JSON(400, gin.H{
				"message": "failed to read body",
			})
			return
		}
		log.Print("signature: " + c.Request.Header.Get("x-line-signature"))
		log.Print(body)
		log.Print("body: " + string(body))

		channelSecret := os.Getenv("CHANNEL_SECRET")
		if channelSecret == "" {
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
		}

		err = validateSignature(channelSecret, c.Request.Header.Get("x-line-signature"), body)

		if err != nil {
			log.Print(err.Error())
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		var event WebhookEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Print("failed to parse json ... " + err.Error())
			c.JSON(400, gin.H{
				"message": "failed to parse json",
			})
			return
		}


		channelSecret := os.Getenv("CHANNEL_SECRET")
		if channelSecret == "" {
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
		}

		bot := createLineBot(channelSecret, accessToken)

		if bot == nil {
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
		}

		for _, event := range event.Events {
			log.Print("reply token = " + event.ReplyToken)
			if event.Message.Type == "text" {
				log.Print("message text = " + event.Message.Text)

				// 届いたメッセージをそのまま返却する
				replyMessage := linebot.NewTextMessage(event.Message.Text)

				_, err = bot.ReplyMessage(event.ReplyToken, replyMessage).Do()

				if err != nil {
					log.Print("failed to reply message ... " + err.Error())

					c.JSON(400, gin.H{
						"message": "Failed to Reply",
					})
				}
			} else {
				log.Print("message type = " + event.Message.Type)
			}
		}



		// response
		c.JSON(200, gin.H{})
	})

	return r
}

func readDotEnv() error {
	return godotenv.Load()
}

func validateSignature(channelSecret string, signature string, body []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("failed to read signature")
	}

	hash := hmac.New(sha256.New, []byte(channelSecret))
	hash.Write(body)

	if !hmac.Equal(decoded, hash.Sum(nil)) {
		return errors.New("illeagal signature")
	}

	return nil
}

func createLineBot(channelSecret string, accessToken string) (*linebot.Client) {
	bot, err := linebot.New(channelSecret, accessToken)

	if err != nil {
		log.Print("failed to create line bot ... " + err.Error())

		return nil
	}

	return bot
}

func main() {
	if len(os.Args) < 2 {
		r := setupRouter()

		r.Run(":80")
	}
}
