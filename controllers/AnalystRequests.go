package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/palladiumkenya/individual-data-request-backend/internal/db"
	"github.com/palladiumkenya/individual-data-request-backend/internal/models"
	"github.com/palladiumkenya/individual-data-request-backend/services"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetApprovedTasks(c *gin.Context) {
	DB, err := db.Connect()
	assigneeUuidStr := c.Query("assignee")
	assigneeUuid, err := uuid.Parse(assigneeUuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee UUID"})
		log.Fatalf("Error invalid UUID: %v\n", err)
		return
	}

	// Retrieve a requests
	requests, err := models.GetAssigneeTasks(DB, assigneeUuid)
	if err != nil {
		log.Fatalf("Error retrieving requests: %v\n", err)
	}

	// Set the Content-Type header and write the JSON response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   requests,
	})

}

func GetApprovedTask(c *gin.Context) {
	DB, err := db.Connect()

	idUuidStr := c.Query("id")
	idUuid, err := uuid.Parse(idUuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request UUID"})
		log.Fatalf("Error invalid UUID: %v\n", err)
		return
	}

	// Retrieve a requests
	request, err := models.GetAssigneeTask(DB, idUuid)
	if err != nil {
		log.Fatalf("Error retrieving requests: %v\n", err)
	}

	// Set the Content-Type header and write the JSON response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   request,
	})

}

func UpdateAnalystRequest(c *gin.Context) {
	DB, err := db.Connect()

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request id"})
		log.Fatalf("Error invalid int: %v\n", err)
		return
	}

	var requestStatus models.UpdateStatusRequest
	if err := c.BindJSON(&requestStatus); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	// Updated Status a requests
	err = models.UpdateRequestStatus(DB, idInt, requestStatus.Status)
	if err != nil {
		log.Fatalf("Error retrieving requests: %v\n", err)
	}

	// Set the Content-Type header and write the JSON response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   "Updated request",
	})

	// Launch background job to send email alert
	go func() {

		// Get Request Details
		request, _ := models.GetRequestByReqID(DB, idInt)

		// send update email to requester
		template := "email_templates/requester_request_status_notification.html"
		frontendUrl := os.Getenv("FRONTEND_URL")
		body := map[string]interface{}{
			"request_id":   request.ID,
			"request_url":  frontendUrl + "/requester/request-details?id=" + request.ID.String(),
			"frontend_url": frontendUrl,
		}

		requester, _ := models.GetRequesterByID(DB, request.Requestor_id)
		email := requester.Email
		subject := "Update on your request"

		emailId, err := services.SendEmailAlerts(subject, body, email, template, c)
		if err != nil {
			log.Fatalf("Error sending email: %v\n", err)
		} else {
			fmt.Printf("Email sent successfully. Email ID: %s\n", emailId)
		}

	}()

}

func GetAssignedAnalyst(c *gin.Context) {
	requestId := c.Param("request_id")

	approvals, err := models.GetAssignedAnalyst(DB, uuid.MustParse(requestId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Printf("Return assignee results")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   approvals,
	})

}
