package repositories

import (
	"ams-backend/models"
	"database/sql"
)

type ProgramRepository struct{}

func NewProgramRepository() *ProgramRepository { return &ProgramRepository{} }

const programSelect = `
	SELECT p.id, p.name, p.code, COALESCE(p.description,''), p.duration_years,
		p.department_id, COALESCE(d.name,''), p.is_active, p.created_at, p.updated_at
	FROM programs p LEFT JOIN departments d ON d.id = p.department_id`

func scanProgram(row interface{ Scan(...interface{}) error }) (models.Program, error) {
	var p models.Program
	err := row.Scan(&p.ID, &p.Name, &p.Code, &p.Description, &p.DurationYears,
		&p.DepartmentID, &p.DepartmentName, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *ProgramRepository) List() ([]models.Program, error) {
	rows, err := DB.Query(programSelect + " ORDER BY p.name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Program
	for rows.Next() {
		p, err := scanProgram(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *ProgramRepository) FindByID(id int64) (*models.Program, error) {
	p, err := scanProgram(DB.QueryRow(programSelect+" WHERE p.id = $1", id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProgramRepository) Create(p *models.Program) (int64, error) {
	var id int64
	err := DB.QueryRow(`
		INSERT INTO programs (name, code, description, duration_years, department_id)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		p.Name, p.Code, p.Description, p.DurationYears, p.DepartmentID,
	).Scan(&id)
	return id, err
}

func (r *ProgramRepository) Update(id int64, p *models.Program) error {
	_, err := DB.Exec(`
		UPDATE programs SET name=$1, code=$2, description=$3, duration_years=$4,
			department_id=$5, is_active=$6, updated_at=NOW() WHERE id=$7`,
		p.Name, p.Code, p.Description, p.DurationYears, p.DepartmentID, p.IsActive, id)
	return err
}

func (r *ProgramRepository) CodeExists(code string) (bool, error) {
	var c int
	err := DB.QueryRow(`SELECT COUNT(*) FROM programs WHERE code=$1`, code).Scan(&c)
	return c > 0, err
}
