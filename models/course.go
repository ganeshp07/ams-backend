package models

import "time"

type Course struct {
	ID             int64     `json:"id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Credits        int       `json:"credits"`
	CourseType     string    `json:"course_type"`
	ProgramID      int64     `json:"program_id"`
	ProgramName    string    `json:"program_name"`
	SemesterNumber int       `json:"semester_number"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CourseOffering struct {
	ID           int64     `json:"id"`
	CourseID     int64     `json:"course_id"`
	CourseCode   string    `json:"course_code"`
	CourseName   string    `json:"course_name"`
	Credits      int       `json:"credits"`
	CourseType   string    `json:"course_type"`
	SemesterID   int64     `json:"semester_id"`
	SemesterName string    `json:"semester_name"`
	AcademicYear string    `json:"academic_year"`
	FacultyID    int64     `json:"faculty_id"`
	FacultyName  string    `json:"faculty_name"`
	Section      string    `json:"section"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// CourseStudent is a student row within a course offering (for faculty view).
type CourseStudent struct {
	StudentID    int64  `json:"student_id"`
	RollNumber   string `json:"roll_number"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Grade        string `json:"grade"`
	GradeID      *int64 `json:"grade_id"`
	RegistrationID int64 `json:"registration_id"`
}
