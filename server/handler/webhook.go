package handler

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

func handleWebhook(c *gin.Context) (int, string, gin.H) {
	// check signature
	// https://developers.line.biz/ja/reference/messaging-api/#signature-validation
	defer c.Request.Body.Close()

	b := &bytes.Buffer{}
	_, err := io.Copy(b, c.Request.Body)
	bodyCopy := bytes.NewReader(b.Bytes())
	c.Request.Body = io.NopCloser(bodyCopy)
	bodyCopy.Seek(0, 0)
	body := b.Bytes()

	signature := c.Request.Header.Get("x-line-signature")

	// body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Print("failed to read body")
		return 400, "failed to read body", nil
	}
	log.Print("signature: " + signature)
	log.Print(body)
	log.Print("body: " + string(body))

	accessToken := os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	if accessToken == "" || channelSecret == "" {
		return 500, "Internal Server Error", nil
	}
	err = validateSignature(channelSecret, signature, body)

	if err != nil {
		log.Print(err.Error())
		return 400, "invalid signature", nil
	}

	var event WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Print("failed to parse json ... " + err.Error())
		return 400, "failed to parse json", nil
	}

	bot := createLineBot(channelSecret, accessToken)

	if bot == nil {
		return 500, "Internal Server Error", nil
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

				return 400, "Failed to Reply", nil
			}
		} else {
			log.Print("message type = " + event.Message.Type)
		}
	}

	return 200, "", gin.H{}
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
