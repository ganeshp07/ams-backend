package services

import (
	"ams-backend/models"
	"ams-backend/repositories"
	"database/sql"
	"fmt"
)

type AdminService struct {
	db          *sql.DB
	userRepo    *repositories.UserRepository
	studentRepo *repositories.StudentRepository
	facultyRepo *repositories.FacultyRepository
}

func NewAdminService(
	db *sql.DB,
	ur *repositories.UserRepository,
	sr *repositories.StudentRepository,
	fr *repositories.FacultyRepository,
) *AdminService {
	return &AdminService{db, ur, sr, fr}
}

type CreateStudentInput struct {
	RollNumber      string  `json:"roll_number"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	Phone           string  `json:"phone"`
	DateOfBirth     *string `json:"date_of_birth"`
	Gender          string  `json:"gender"`
	Address         string  `json:"address"`
	ProgramID       int64   `json:"program_id"`
	AdmissionYear   int     `json:"admission_year"`
	CurrentSemester int     `json:"current_semester"`
}

func (s *AdminService) CreateStudent(input *CreateStudentInput, createdBy int64) (*models.StudentProfile, error) {
	if exists, _ := s.userRepo.EmailExists(input.Email); exists {
		return nil, fmt.Errorf("email already in use")
	}
	if exists, _ := s.studentRepo.RollExists(input.RollNumber); exists {
		return nil, fmt.Errorf("roll number already in use")
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userID, err := s.userRepo.CreateWithTx(tx, input.Email, hash, models.RoleStudent)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account")
	}

	if input.CurrentSemester == 0 {
		input.CurrentSemester = 1
	}

	student := &models.Student{
		UserID:          userID,
		RollNumber:      input.RollNumber,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Phone:           input.Phone,
		DateOfBirth:     input.DateOfBirth,
		Gender:          input.Gender,
		Address:         input.Address,
		ProgramID:       input.ProgramID,
		AdmissionYear:   input.AdmissionYear,
		CurrentSemester: input.CurrentSemester,
	}
	_, err = s.studentRepo.CreateWithTx(tx, student, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create student profile")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.studentRepo.FindByUserID(userID)
}

type CreateFacultyInput struct {
	EmployeeID   string  `json:"employee_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	Phone        string  `json:"phone"`
	DepartmentID *int64  `json:"department_id"`
	Designation  string  `json:"designation"`
}

func (s *AdminService) CreateFaculty(input *CreateFacultyInput, createdBy int64) (*models.FacultyProfile, error) {
	if exists, _ := s.userRepo.EmailExists(input.Email); exists {
		return nil, fmt.Errorf("email already in use")
	}
	if exists, _ := s.facultyRepo.EmployeeIDExists(input.EmployeeID); exists {
		return nil, fmt.Errorf("employee ID already in use")
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userID, err := s.userRepo.CreateWithTx(tx, input.Email, hash, models.RoleFaculty)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account")
	}

	faculty := &models.Faculty{
		UserID:       userID,
		EmployeeID:   input.EmployeeID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		DepartmentID: input.DepartmentID,
		Designation:  input.Designation,
	}
	_, err = s.facultyRepo.CreateWithTx(tx, faculty, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create faculty profile")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.facultyRepo.FindByUserID(userID)
}
