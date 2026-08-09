package repositories

import (
	"ams-backend/models"
	"database/sql"
)

type SemesterRepository struct{}

func NewSemesterRepository() *SemesterRepository { return &SemesterRepository{} }

func (r *SemesterRepository) List() ([]models.Semester, error) {
	rows, err := DB.Query(`
		SELECT id, academic_year, semester_number, name,
			TO_CHAR(start_date,'YYYY-MM-DD'), TO_CHAR(end_date,'YYYY-MM-DD'),
			registration_start, registration_end, is_active, created_at, updated_at
		FROM semesters ORDER BY academic_year DESC, semester_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Semester
	for rows.Next() {
		s, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SemesterRepository) FindByID(id int64) (*models.Semester, error) {
	s, err := r.scan(DB.QueryRow(`
		SELECT id, academic_year, semester_number, name,
			TO_CHAR(start_date,'YYYY-MM-DD'), TO_CHAR(end_date,'YYYY-MM-DD'),
			registration_start, registration_end, is_active, created_at, updated_at
		FROM semesters WHERE id=$1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SemesterRepository) GetActive() (*models.Semester, error) {
	s, err := r.scan(DB.QueryRow(`
		SELECT id, academic_year, semester_number, name,
			TO_CHAR(start_date,'YYYY-MM-DD'), TO_CHAR(end_date,'YYYY-MM-DD'),
			registration_start, registration_end, is_active, created_at, updated_at
		FROM semesters WHERE is_active = true LIMIT 1`))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SemesterRepository) Create(s *models.Semester) (int64, error) {
	var id int64
	err := DB.QueryRow(`
		INSERT INTO semesters (academic_year, semester_number, name, start_date, end_date,
			registration_start, registration_end, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		s.AcademicYear, s.SemesterNumber, s.Name, s.StartDate, s.EndDate,
		s.RegistrationStart, s.RegistrationEnd, s.IsActive,
	).Scan(&id)
	return id, err
}

func (r *SemesterRepository) Update(id int64, s *models.Semester) error {
	_, err := DB.Exec(`
		UPDATE semesters SET academic_year=$1, semester_number=$2, name=$3,
			start_date=$4, end_date=$5, registration_start=$6, registration_end=$7,
			is_active=$8, updated_at=NOW() WHERE id=$9`,
		s.AcademicYear, s.SemesterNumber, s.Name, s.StartDate, s.EndDate,
		s.RegistrationStart, s.RegistrationEnd, s.IsActive, id)
	return err
}

func (r *SemesterRepository) scan(row interface{ Scan(...interface{}) error }) (models.Semester, error) {
	var s models.Semester
	err := row.Scan(&s.ID, &s.AcademicYear, &s.SemesterNumber, &s.Name,
		&s.StartDate, &s.EndDate, &s.RegistrationStart, &s.RegistrationEnd,
		&s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}
