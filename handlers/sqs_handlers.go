package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/espher/GoLang-API-REST/sqs"
)

type SQSHAndler struct {
	db *gorm.DB
}

func SQSRouter(db *gorm.DB) *SQSHAndler {
	return &SQSHAndler{db: db}
}

const (
	SQS_URL      = "https://sqs.us-east-1.amazonaws.com/480033710745/GoQueueTest"
	SQS_MAX_MSGS = 10
)

func (uh *SQSHAndler) GetMessages(c *gin.Context) {
	result, err := sqs.ReceiveMessages(SQS_URL, SQS_MAX_MSGS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (uh *SQSHAndler) CreateMessage(c *gin.Context) {
	type MessageBody struct {
		Message string
	}

	var requestBodyMessage MessageBody
	c.BindJSON(&requestBodyMessage)

	err := sqs.SendMessage(SQS_URL, requestBodyMessage.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "message sent")
}

func (uh *SQSHAndler) CreateMessageBulk(c *gin.Context) {
	messages := []string{
		"Mensaje 1",
		"Mensaje 2",
		"Mensaje 3",
		"Mensaje 4",
		"Mensaje 5",
	}

	err := sqs.SendMessagesBulk(SQS_URL, messages)
	if err != nil {
		log.Fatal(err)
	}

}
