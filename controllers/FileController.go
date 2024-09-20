package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/palladiumkenya/individual-data-request-backend/services"
	"net/http"
	"path/filepath"
)

func UploadFile(c *gin.Context) {
	// unique id for folder name
	folderId := uuid.New().String()

	// Get destination folder
	destinationFolder := c.PostForm("destination") // either "files" or "supporting-documents"
	request := c.PostForm("request")               // request that this file is associated with
	var requestUuid *uuid.UUID
	if request != "" {
		parsedUuid, _ := uuid.Parse(request)
		requestUuid = &parsedUuid
	}

	// Get the file from the POST form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the file locally
	localFilePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, localFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Upload to Nextcloud and get the file URL
	remoteFilePath := "/idr/files/" + file.Filename
	fileURL, err := services.UploadFileToNextcloud(localFilePath, remoteFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to Nextcloud"})
		return
	}

	// Save the file URL to the database
	// TODO:: Save And User ID
	DB, err := db.Connect()
	requestFile := models.Files{
		FileName:  file.Filename,
		FileURL:   fileURL,
		RequestId: requestUuid,
	}

	// Save the file details to the database
	if err := models.UploadFiles(DB, &requestFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file details to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
}

