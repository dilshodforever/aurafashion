package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Webhook payload structure
type WebhookEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func (h Handler) WebhookHandler(ctx *gin.Context) {
	var event WebhookEvent
	// Decode the incoming JSON payload
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Process each message received from Instagram
	for _, entry := range event.Entry {
		for _, messaging := range entry.Messaging {
			// If a message is received, respond with "Please wait..."
			if messaging.Message.Text != "" {
				// Call function to respond to Instagram
				RespondToInstagram(messaging.Sender.ID, "Iltimos ozroq kuting...")
			}
		}
	}

	// Respond back with OK status
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func RespondToInstagram(userID, message string) {
	// Your access token
	accessToken := "IGAAJ9eVGFRoNBZAE1DRG5wOHptejQwRE1TczV3SWdaMjhZAM1FQMFk3QWdyNldjVXVCV3BydmF5eU40TnBrc19PTGVrMjBsNnhKc25vS2dBTlZAFTElQRUt1M1hLS2xPUGEyWjdMM0E0bkY3cTdwajVxRnZAEN0dpZAVBUYlFGbG1DYwZDZD"

	// Instagram API URL for sending messages
	url := fmt.Sprintf("https://graph.facebook.com/v12.0/me/messages?access_token=%s", accessToken)

	// Prepare message payload
	payload := fmt.Sprintf(`{
		"recipient": {"id": "%s"},
		"message": {"text": "%s"}
	}`, userID, message)

	// Create a new HTTP client and send the POST request to the Instagram API
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
	} else {
		log.Printf("Failed to send message. Status Code: %d", resp.StatusCode)
	}
}
