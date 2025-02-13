package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Requests struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReqId          int        `gorm:"type:integer;autoIncrement;unique;not null"`
	Summery        string     `gorm:"size:1500;not null"`
	Status         string     `gorm:"size:100;not null"`
	Date_Due       time.Time  `gorm:"type:date"`
	Priority_level string     `gorm:"size:100;not null"`
	Requestor_id   uuid.UUID  `gorm:"type:uuid"`
	Requester      Requesters `gorm:"foreignKey:Requestor_id"`
	Assignee_id    *uuid.UUID `gorm:"type:uuid;null"`
	Assignee       Assignees  `gorm:"foreignKey:Assignee_id"`
	Created_Date   time.Time  `gorm:"type:date"`
}

type NewRequest struct {
	Summery      string    `json:"summery" binding:"required"`
	Priority     string    `json:"priority" binding:"required"`
	DateDue      time.Time `json:"dateDue" binding:"required"`
	Requestor_id uuid.UUID `json:"requestor_id" binding:"required"`
}

//func GetRequestByID(DB *gorm.DB, Id uuid.UUID) (*Requests, error) {
//	var request *Requests
//	result := DB.First(&request, "id = ?", Id)
//	return request, result.Error
//}

func GetRequestByID(DB *gorm.DB, Id uuid.UUID) (*Requests, error) {
	var request *Requests
	result := DB.Preload("Requester").First(&request, "id = ?", Id)
	return request, result.Error
}

func GetRequestByReqID(DB *gorm.DB, Id int) (*Requests, error) {
	var request *Requests
	result := DB.Preload("Requester").First(&request, "req_id = ?", Id)
	return request, result.Error
}

func GetRequests(DB *gorm.DB) ([]Requests, error) {
	var requests []Requests
	result := DB.Preload("Requester").Find(&requests)
	return requests, result.Error
}

func CreateRequest(DB *gorm.DB, newRequest NewRequest) (Requests, error) {
	request := Requests{
		ID:             uuid.New(), // Generate a new UUID for the request
		Summery:        newRequest.Summery,
		Status:         "pending", // Default status
		Date_Due:       newRequest.DateDue,
		Priority_level: newRequest.Priority,
		Created_Date:   time.Now(), // Set to current time
		Requestor_id:   newRequest.Requestor_id,
		Assignee_id:    nil,
	}

	if err := DB.Create(&request).Error; err != nil {
		return request, err // Return the error and a zero UUID
	}
	return request, nil // Return the ID of the created request
}

func GetAssigneeTasks(DB *gorm.DB, assignee uuid.UUID) ([]Requests, error) {
	var requests []Requests
	result := DB.Preload("Requester").Preload("Assignee").Where("assignee_id =?", assignee).Find(&requests)
	return requests, result.Error
}

func GetRequesterRequests(DB *gorm.DB, requester_id uuid.UUID) ([]Requests, error) {
	var requests []Requests
	result := DB.Preload("Requester").Preload("Assignee").Where("requestor_id =?", requester_id).Find(&requests)
	return requests, result.Error
}

func GetAssigneeTask(DB *gorm.DB, id uuid.UUID) ([]Requests, error) {
	var requests []Requests
	result := DB.Preload("Requester").Preload("Assignee").First(&requests, "ID = ?", id)
	return requests, result.Error
}

func UpdateRequestStatus(DB *gorm.DB, requestID int, newStatus string) error {
	// Find the request by ID
	var request Requests
	if err := DB.First(&request, "req_id = ?", requestID).Error; err != nil {
		return err // Return the error if the request is not found or if there is another issue
	}

	// Update the status
	request.Status = newStatus
	if err := DB.Save(&request).Error; err != nil {
		return err // Return the error if the update fails
	}

	return nil // Return nil if the update is successful
}

func AssignRequestToAnalyst(DB *gorm.DB, Id uuid.UUID, analystID uuid.UUID) (error, error) {
	// Find the request by ID
	var request Requests
	if err := DB.First(&request, "id = ?", Id).Error; err != nil {
		return err, nil // Return the error if the request is not found or if there is another issue
	}

	//Update the status
	request.Assignee_id = &analystID
	request.Status = "assigned"
	if err := DB.Save(&request).Error; err != nil {
		return err, nil // Return the error if the update fails
	}

	return nil, nil // Return nil if the update is successful
}

func GetRequesterRequestDetails(DB *gorm.DB, id uuid.UUID) (Requests, error) {
	var request Requests
	result := DB.Preload("Requester").Preload("Assignee").First(&request, "id = ?", id)
	return request, result.Error
}
