package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/palladiumkenya/individual-data-request-backend/internal/db"
	"github.com/palladiumkenya/individual-data-request-backend/internal/models"
	"net/http"
)

// CreateChatMessage handles the creation of a new chat message
func CreateChatMessage(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	var input struct {
		RequestID uuid.UUID `json:"request_id" binding:"required"`
		SenderID  uuid.UUID `json:"sender_id" binding:"required"`
		Message   string    `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatMessage, err := models.NewMessage(DB, input.RequestID, input.SenderID, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   chatMessage,
	})
}

// GetChatMessagesByRequestID handles fetching all chat messages for a specific request
func GetChatMessagesByRequestID(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	requestIDStr := c.Param("request_id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request UUID"})
		return
	}

	chatMessages, err := models.GetChatMessagesByRequestID(DB, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   chatMessages,
	})
}

// GetChatMessagesBetweenUsers handles fetching all chat messages between two users for a specific request
func GetChatMessagesBetweenUsers(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	requestIDStr := c.Param("request_id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request UUID"})
		return
	}

	user1IDStr := c.Query("user1_id")
	user1ID, err := uuid.Parse(user1IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1 UUID"})
		return
	}

	user2IDStr := c.Query("user2_id")
	user2ID, err := uuid.Parse(user2IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user2 UUID"})
		return
	}

	chatMessages, err := models.GetChatMessagesBetweenUsers(DB, requestID, user1ID, user2ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   chatMessages,
	})
}
