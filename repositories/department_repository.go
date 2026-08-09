package repositories

import "ams-backend/models"

type DepartmentRepository struct{}

func NewDepartmentRepository() *DepartmentRepository { return &DepartmentRepository{} }

func (r *DepartmentRepository) List() ([]models.Department, error) {
	rows, err := DB.Query(`SELECT id, name, code, created_at, updated_at FROM departments ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
