package repositories

import (
	"ams-backend/models"
	"database/sql"
)

type GradeRepository struct{}

func NewGradeRepository() *GradeRepository { return &GradeRepository{} }

func (r *GradeRepository) Upsert(registrationID, facultyID int64, grade string) (*models.Grade, error) {
	var g models.Grade
	err := DB.QueryRow(`
		INSERT INTO grades (registration_id, grade, graded_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (registration_id) DO UPDATE
			SET grade=EXCLUDED.grade, graded_by=EXCLUDED.graded_by,
			    graded_at=NOW(), updated_at=NOW()
		RETURNING id, registration_id, grade, graded_by, graded_at, created_at, updated_at`,
		registrationID, grade, facultyID,
	).Scan(&g.ID, &g.RegistrationID, &g.Grade, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	return &g, err
}

func (r *GradeRepository) FindByRegistration(registrationID int64) (*models.Grade, error) {
	var g models.Grade
	err := DB.QueryRow(`
		SELECT id, registration_id, grade, graded_by, graded_at, created_at, updated_at
		FROM grades WHERE registration_id=$1`, registrationID,
	).Scan(&g.ID, &g.RegistrationID, &g.Grade, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &g, err
}
