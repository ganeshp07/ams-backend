package models

import "time"

type Student struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	RollNumber      string    `json:"roll_number"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone"`
	DateOfBirth     *string   `json:"date_of_birth"`
	Gender          string    `json:"gender"`
	Address         string    `json:"address"`
	ProgramID       int64     `json:"program_id"`
	AdmissionYear   int       `json:"admission_year"`
	CurrentSemester int       `json:"current_semester"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// StudentProfile is the flat view returned by API responses.
type StudentProfile struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	Email           string  `json:"email"`
	RollNumber      string  `json:"roll_number"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	FullName        string  `json:"full_name"`
	Phone           string  `json:"phone"`
	DateOfBirth     *string `json:"date_of_birth"`
	Gender          string  `json:"gender"`
	Address         string  `json:"address"`
	ProgramID       int64   `json:"program_id"`
	ProgramName     string  `json:"program_name"`
	ProgramCode     string  `json:"program_code"`
	AdmissionYear   int     `json:"admission_year"`
	CurrentSemester int     `json:"current_semester"`
	IsActive        bool    `json:"is_active"`
}

// StudentCourse is returned for course history and grade views.
type StudentCourse struct {
	RegistrationID    int64   `json:"registration_id"`
	CourseOfferingID  int64   `json:"course_offering_id"`
	CourseCode        string  `json:"course_code"`
	CourseName        string  `json:"course_name"`
	Credits           int     `json:"credits"`
	CourseType        string  `json:"course_type"`
	SemesterName      string  `json:"semester_name"`
	AcademicYear      string  `json:"academic_year"`
	FacultyName       string  `json:"faculty_name"`
	Grade             string  `json:"grade"`
}
