package repositories

import (
	"ams-backend/models"
	"database/sql"
	"fmt"
)

type FacultyRepository struct{}

func NewFacultyRepository() *FacultyRepository { return &FacultyRepository{} }

const facultySelect = `
	SELECT f.id, f.user_id, u.email, f.employee_id,
		f.first_name, f.last_name, f.first_name||' '||f.last_name,
		COALESCE(f.phone,''), f.department_id, COALESCE(d.name,''),
		COALESCE(f.designation,''), f.is_active
	FROM faculty f
	JOIN users u ON u.id = f.user_id
	LEFT JOIN departments d ON d.id = f.department_id`

func scanFacultyRow(rows interface{ Scan(...interface{}) error }) (models.FacultyProfile, error) {
	var fp models.FacultyProfile
	err := rows.Scan(&fp.ID, &fp.UserID, &fp.Email, &fp.EmployeeID,
		&fp.FirstName, &fp.LastName, &fp.FullName,
		&fp.Phone, &fp.DepartmentID, &fp.DepartmentName,
		&fp.Designation, &fp.IsActive)
	return fp, err
}

func (r *FacultyRepository) FindByUserID(userID int64) (*models.FacultyProfile, error) {
	fp, err := scanFacultyRow(DB.QueryRow(facultySelect+" WHERE f.user_id = $1", userID))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fp, nil
}

func (r *FacultyRepository) FindByID(id int64) (*models.FacultyProfile, error) {
	fp, err := scanFacultyRow(DB.QueryRow(facultySelect+" WHERE f.id = $1", id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fp, nil
}

func (r *FacultyRepository) List(search string) ([]models.FacultyProfile, error) {
	args := []interface{}{}
	where := ""
	if search != "" {
		where = fmt.Sprintf(" WHERE (f.employee_id ILIKE $1 OR f.first_name ILIKE $1 OR f.last_name ILIKE $1 OR u.email ILIKE $1)")
		args = append(args, "%"+search+"%")
	}
	rows, err := DB.Query(facultySelect+where+" ORDER BY f.first_name", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FacultyProfile
	for rows.Next() {
		fp, err := scanFacultyRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, fp)
	}
	return out, rows.Err()
}

// ListActive returns only active faculty (for dropdowns).
func (r *FacultyRepository) ListActive() ([]models.FacultyProfile, error) {
	rows, err := DB.Query(facultySelect+" WHERE f.is_active = true ORDER BY f.first_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FacultyProfile
	for rows.Next() {
		fp, err := scanFacultyRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, fp)
	}
	return out, rows.Err()
}

func (r *FacultyRepository) CreateWithTx(tx *sql.Tx, f *models.Faculty, createdBy int64) (int64, error) {
	var id int64
	err := tx.QueryRow(`
		INSERT INTO faculty (user_id, employee_id, first_name, last_name, phone, department_id, designation, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		f.UserID, f.EmployeeID, f.FirstName, f.LastName, f.Phone, f.DepartmentID, f.Designation, createdBy,
	).Scan(&id)
	return id, err
}

func (r *FacultyRepository) Update(id int64, f *models.Faculty) error {
	_, err := DB.Exec(`
		UPDATE faculty SET employee_id=$1, first_name=$2, last_name=$3,
			phone=$4, department_id=$5, designation=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8`,
		f.EmployeeID, f.FirstName, f.LastName, f.Phone,
		f.DepartmentID, f.Designation, f.IsActive, id,
	)
	return err
}

func (r *FacultyRepository) EmployeeIDExists(empID string) (bool, error) {
	var c int
	err := DB.QueryRow(`SELECT COUNT(*) FROM faculty WHERE employee_id=$1`, empID).Scan(&c)
	return c > 0, err
}

func (r *FacultyRepository) TeachesCourseOffering(facultyID, offeringID int64) (bool, error) {
	var c int
	err := DB.QueryRow(
		`SELECT COUNT(*) FROM course_offerings WHERE id=$1 AND faculty_id=$2`,
		offeringID, facultyID,
	).Scan(&c)
	return c > 0, err
}

func (r *FacultyRepository) GetCurrentCourses(facultyID int64) ([]models.CourseOffering, error) {
	return r.queryCourses(`
		WHERE co.faculty_id = $1 AND s.is_active = true`, facultyID)
}

func (r *FacultyRepository) GetCourseHistory(facultyID int64) ([]models.CourseOffering, error) {
	return r.queryCourses(`
		WHERE co.faculty_id = $1 AND s.is_active = false`, facultyID)
}

func (r *FacultyRepository) queryCourses(where string, args ...interface{}) ([]models.CourseOffering, error) {
	q := `SELECT co.id, co.course_id, c.code, c.name, c.credits, c.course_type,
		co.semester_id, s.name, co.academic_year,
		co.faculty_id, f.first_name||' '||f.last_name, co.section, co.is_active, co.created_at
		FROM course_offerings co
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id ` + where + ` ORDER BY s.academic_year DESC, s.semester_number`

	rows, err := DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CourseOffering
	for rows.Next() {
		var co models.CourseOffering
		if err := rows.Scan(&co.ID, &co.CourseID, &co.CourseCode, &co.CourseName,
			&co.Credits, &co.CourseType, &co.SemesterID, &co.SemesterName,
			&co.AcademicYear, &co.FacultyID, &co.FacultyName,
			&co.Section, &co.IsActive, &co.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, co)
	}
	return out, rows.Err()
}

// Satisfies the interface used by scanFacultyRow with sql.Rows
