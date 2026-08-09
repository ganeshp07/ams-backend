package models

import "time"

type Semester struct {
	ID                int64      `json:"id"`
	AcademicYear      string     `json:"academic_year"`
	SemesterNumber    int        `json:"semester_number"`
	Name              string     `json:"name"`
	StartDate         *string    `json:"start_date"`
	EndDate           *string    `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (s *Semester) IsRegistrationOpen() bool {
	if s.RegistrationStart == nil || s.RegistrationEnd == nil {
		return false
	}
	now := time.Now()
	return now.After(*s.RegistrationStart) && now.Before(*s.RegistrationEnd)
}
