package models

import "time"

var ValidGrades = map[string]bool{
	"A+": true, "A": true, "A-": true,
	"B+": true, "B": true, "B-": true,
	"C+": true, "C": true, "C-": true,
	"D": true, "F": true, "I": true, "W": true,
}

func IsValidGrade(g string) bool {
	return ValidGrades[g]
}

type Grade struct {
	ID             int64     `json:"id"`
	RegistrationID int64     `json:"registration_id"`
	Grade          string    `json:"grade"`
	GradedBy       *int64    `json:"graded_by"`
	GradedAt       time.Time `json:"graded_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type BulkUploadResult struct {
	Successful int                `json:"successful"`
	Failed     int                `json:"failed"`
	Errors     []BulkUploadError  `json:"errors"`
}

type BulkUploadError struct {
	Row        int    `json:"row"`
	RollNumber string `json:"roll_number"`
	Message    string `json:"message"`
}
