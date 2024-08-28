package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/gorm"

	//"github.com/palladiumkenya/individual-data-request-backend/internal/db"
	"time"
)

type Requests struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReqId          int        `gorm:"type:integer;unique;not null"`
	Summery        string     `gorm:"size:500;not null"`
	Status         string     `gorm:"size:100;not null"`
	Date_Due       time.Time  `gorm:"type:date"`
	Priority_level string     `gorm:"size:100;unique;not null"`
	Requestor_id   uuid.UUID  `gorm:"type:uuid"`
	Requester      Requesters `gorm:"foreignKey:Requestor_id"`
	Assignee_id    uuid.UUID  `gorm:"type:uuid;null"`
	Assignee       Assignees  `gorm:"foreignKey:Assignee_id"`
	Created_Date   time.Time  `gorm:"type:date"`
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

func GetRequests(DB *gorm.DB) ([]Requests, error) {
	var requests []Requests
	result := DB.Find(&requests)
	return requests, result.Error
}

func CreateRequest(ctx context.Context, pool *pgxpool.Pool, request *Requests) error {
	_, err := pool.Exec(ctx, "INSERT INTO Requests (summery, status, due_date, priority_level, requestor_id) VALUES ($1, $2, $3)",
		request.Summery, request.Status, request.Date_Due, request.Priority_level)
	if err != nil {
		return err
	}
	return nil
}

func GetAssigneeTasks(DB *gorm.DB, assignee uuid.UUID) ([]Requests, error) {
	var requests []Requests
	result := DB.Preload("Requester").Preload("Assignee").Where("assignee_id =?", assignee).Find(&requests)
	return requests, result.Error
}
