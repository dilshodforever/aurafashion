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

// WebhookHandler - Webhookni tasdiqlash va xabarni qaytarish
func (h Handler) WebhookHandler(c *gin.Context) {
	// Webhook verification uchun GET so'rovini tekshirish
	mode := c.DefaultQuery("hub.mode", "")
	// challenge := c.DefaultQuery("hub.challenge", "")
	verifyToken := c.DefaultQuery("hub.verify_token", "")

	// Verification tokenni tekshirish
	if verifyToken != "12345" { // 'your-verify-token'ni o'z tokeningiz bilan almashtiring
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid verify token"})
		return
	}

	// Agar mode 'subscribe' bo'lsa, challenge qaytarish
	if mode == "subscribe" {
		fmt.Println("resp in subscribe")
		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	}

	// POST so'rovini olish va xabarlarni qayta ishlash
	var event WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Instagramdan kelgan xabarlarni qayta ishlash
	for _, entry := range event.Entry {
		for _, messaging := range entry.Messaging {
			// Agar xabar kelgan bo'lsa, "Iltimos ozroq kuting..." deb javob berish
			if messaging.Message.Text != "" {
				// Instagramga javob yuborish
				RespondToInstagram(messaging.Sender.ID, "Iltimos ozroq kuting...")
			}
		}
	}

	// OK javobini qaytarish
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RespondToInstagram - Instagramga xabar yuborish
func RespondToInstagram(userID, message string) {
	// O'zingizning access tokeningizni qo'shing
	accessToken := "IGAAJ9eVGFRoNBZAE1DRG5wOHptejQwRE1TczV3SWdaMjhZAM1FQMFk3QWdyNldjVXVCV3BydmF5eU40TnBrc19PTGVrMjBsNnhKc25vS2dBTlZAFTElQRUt1M1hLS2xPUGEyWjdMM0E0bkY3cTdwajVxRnZAEN0dpZAVBUYlFGbG1DYwZDZD"

	// Instagram API URL
	url := fmt.Sprintf("https://graph.facebook.com/v12.0/me/messages?access_token=%s", accessToken)

	// Xabarni yuborish uchun payload
	payload := fmt.Sprintf(`{
		"recipient": {"id": "%s"},
		"message": {"text": "%s"}
	}`, userID, message)

	// HTTP client yaratish va so'rov yuborish
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
