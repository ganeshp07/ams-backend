package services

import (
	"ams-backend/models"
	"ams-backend/repositories"
	"ams-backend/utils"
	"fmt"
	"io"
)

type GradeService struct {
	gradeRepo    *repositories.GradeRepository
	courseRepo   *repositories.CourseRepository
	studentRepo  *repositories.StudentRepository
	facultyRepo  *repositories.FacultyRepository
	emailService *EmailService
}

func NewGradeService(
	gr *repositories.GradeRepository,
	cr *repositories.CourseRepository,
	sr *repositories.StudentRepository,
	fr *repositories.FacultyRepository,
	es *EmailService,
) *GradeService {
	return &GradeService{gr, cr, sr, fr, es}
}

func (s *GradeService) UploadGrade(facultyUserID, offeringID, studentID int64, grade string) (*models.Grade, error) {
	// Validate grade value
	if !models.IsValidGrade(grade) {
		return nil, fmt.Errorf("invalid grade: %s", grade)
	}

	// Verify faculty owns this course offering
	faculty, err := s.facultyRepo.FindByUserID(facultyUserID)
	if err != nil || faculty == nil {
		return nil, fmt.Errorf("faculty not found")
	}
	ok, err := s.facultyRepo.TeachesCourseOffering(faculty.ID, offeringID)
	if err != nil || !ok {
		return nil, fmt.Errorf("you are not assigned to this course")
	}

	// Find the registration
	regID, err := s.courseRepo.GetRegistrationForStudent(studentID, offeringID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if regID == 0 {
		return nil, fmt.Errorf("student is not registered for this course")
	}

	g, err := s.gradeRepo.Upsert(regID, faculty.ID, grade)
	if err != nil {
		return nil, fmt.Errorf("failed to save grade")
	}

	// Send async email notification
	s.notifyStudent(studentID, offeringID, grade)
	return g, nil
}

func (s *GradeService) BulkUpload(facultyUserID, offeringID int64, file io.Reader) (*models.BulkUploadResult, error) {
	// Verify faculty owns this offering
	faculty, err := s.facultyRepo.FindByUserID(facultyUserID)
	if err != nil || faculty == nil {
		return nil, fmt.Errorf("faculty not found")
	}
	ok, err := s.facultyRepo.TeachesCourseOffering(faculty.ID, offeringID)
	if err != nil || !ok {
		return nil, fmt.Errorf("you are not assigned to this course")
	}

	rows, err := utils.ParseGradeCSV(file)
	if err != nil {
		return nil, err
	}

	result := &models.BulkUploadResult{Errors: []models.BulkUploadError{}}
	for i, row := range rows {
		rowNum := i + 2 // 1-indexed + header
		if !models.IsValidGrade(row.Grade) {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkUploadError{
				Row: rowNum, RollNumber: row.RollNumber,
				Message: "invalid grade: " + row.Grade,
			})
			continue
		}

		student, err := s.studentRepo.FindByRollNumber(row.RollNumber)
		if err != nil || student == nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkUploadError{
				Row: rowNum, RollNumber: row.RollNumber,
				Message: "student not found",
			})
			continue
		}

		regID, err := s.courseRepo.GetRegistrationForStudent(student.ID, offeringID)
		if err != nil || regID == 0 {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkUploadError{
				Row: rowNum, RollNumber: row.RollNumber,
				Message: "student not registered for this course",
			})
			continue
		}

		if _, err := s.gradeRepo.Upsert(regID, faculty.ID, row.Grade); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkUploadError{
				Row: rowNum, RollNumber: row.RollNumber,
				Message: "database error",
			})
			continue
		}

		result.Successful++
		s.notifyStudent(student.ID, offeringID, row.Grade)
	}
	return result, nil
}

func (s *GradeService) notifyStudent(studentID, offeringID int64, grade string) {
	sp, err := s.studentRepo.FindByID(studentID)
	if err != nil || sp == nil {
		return
	}
	offering, err := s.courseRepo.GetOfferingByID(offeringID)
	if err != nil || offering == nil {
		return
	}
	s.emailService.SendGradeNotification(
		sp.Email, sp.FullName,
		offering.CourseName, offering.CourseCode,
		offering.SemesterName, grade,
	)
}
