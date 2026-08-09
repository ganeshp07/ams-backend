package models

import "time"

type Faculty struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	EmployeeID   string    `json:"employee_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	DepartmentID *int64    `json:"department_id"`
	Designation  string    `json:"designation"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FacultyProfile struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	Email          string `json:"email"`
	EmployeeID     string `json:"employee_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	Phone          string `json:"phone"`
	DepartmentID   *int64 `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Designation    string `json:"designation"`
	IsActive       bool   `json:"is_active"`
}
