package repositories

import (
	"ams-backend/models"
	"database/sql"
	"fmt"
	"strings"
)

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository { return &StudentRepository{} }

const studentProfileSelect = `
	SELECT s.id, s.user_id, u.email, s.roll_number,
		s.first_name, s.last_name, s.first_name||' '||s.last_name,
		COALESCE(s.phone,''), TO_CHAR(s.date_of_birth,'YYYY-MM-DD'),
		COALESCE(s.gender,''), COALESCE(s.address,''),
		s.program_id, COALESCE(p.name,''), COALESCE(p.code,''),
		s.admission_year, s.current_semester, s.is_active
	FROM students s
	JOIN users u ON u.id = s.user_id
	LEFT JOIN programs p ON p.id = s.program_id`

func scanProfile(row *sql.Row) (*models.StudentProfile, error) {
	var sp models.StudentProfile
	err := row.Scan(
		&sp.ID, &sp.UserID, &sp.Email, &sp.RollNumber,
		&sp.FirstName, &sp.LastName, &sp.FullName,
		&sp.Phone, &sp.DateOfBirth,
		&sp.Gender, &sp.Address,
		&sp.ProgramID, &sp.ProgramName, &sp.ProgramCode,
		&sp.AdmissionYear, &sp.CurrentSemester, &sp.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sp, err
}

func (r *StudentRepository) FindByUserID(userID int64) (*models.StudentProfile, error) {
	q := studentProfileSelect + " WHERE s.user_id = $1"
	return scanProfile(DB.QueryRow(q, userID))
}

func (r *StudentRepository) FindByID(id int64) (*models.StudentProfile, error) {
	q := studentProfileSelect + " WHERE s.id = $1"
	return scanProfile(DB.QueryRow(q, id))
}

func (r *StudentRepository) FindByRollNumber(roll string) (*models.Student, error) {
	var s models.Student
	err := DB.QueryRow(`SELECT id, user_id, roll_number, first_name, last_name FROM students WHERE roll_number = $1`, roll).
		Scan(&s.ID, &s.UserID, &s.RollNumber, &s.FirstName, &s.LastName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *StudentRepository) List(search string, programID int64, isActive *bool) ([]models.StudentProfile, error) {
	args := []interface{}{}
	conds := []string{}
	i := 1

	if search != "" {
		conds = append(conds, fmt.Sprintf(
			"(s.roll_number ILIKE $%d OR s.first_name ILIKE $%d OR s.last_name ILIKE $%d OR u.email ILIKE $%d)",
			i, i, i, i))
		args = append(args, "%"+search+"%")
		i++
	}
	if programID > 0 {
		conds = append(conds, fmt.Sprintf("s.program_id = $%d", i))
		args = append(args, programID)
		i++
	}
	if isActive != nil {
		conds = append(conds, fmt.Sprintf("s.is_active = $%d", i))
		args = append(args, *isActive)
	}

	q := `SELECT s.id, s.user_id, u.email, s.roll_number,
		s.first_name, s.last_name, s.first_name||' '||s.last_name,
		COALESCE(s.phone,''), TO_CHAR(s.date_of_birth,'YYYY-MM-DD'),
		COALESCE(s.gender,''), COALESCE(s.address,''),
		s.program_id, COALESCE(p.name,''), COALESCE(p.code,''),
		s.admission_year, s.current_semester, s.is_active
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN programs p ON p.id = s.program_id`

	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY s.roll_number"

	rows, err := DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.StudentProfile
	for rows.Next() {
		var sp models.StudentProfile
		if err := rows.Scan(
			&sp.ID, &sp.UserID, &sp.Email, &sp.RollNumber,
			&sp.FirstName, &sp.LastName, &sp.FullName,
			&sp.Phone, &sp.DateOfBirth,
			&sp.Gender, &sp.Address,
			&sp.ProgramID, &sp.ProgramName, &sp.ProgramCode,
			&sp.AdmissionYear, &sp.CurrentSemester, &sp.IsActive,
		); err != nil {
			return nil, err
		}
		out = append(out, sp)
	}
	return out, rows.Err()
}

func (r *StudentRepository) CreateWithTx(tx *sql.Tx, s *models.Student, createdBy int64) (int64, error) {
	var id int64
	err := tx.QueryRow(`
		INSERT INTO students (user_id, roll_number, first_name, last_name, phone,
			date_of_birth, gender, address, program_id, admission_year, current_semester, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		s.UserID, s.RollNumber, s.FirstName, s.LastName, s.Phone,
		s.DateOfBirth, s.Gender, s.Address, s.ProgramID, s.AdmissionYear, s.CurrentSemester, createdBy,
	).Scan(&id)
	return id, err
}

func (r *StudentRepository) Update(id int64, s *models.Student) error {
	_, err := DB.Exec(`
		UPDATE students SET roll_number=$1, first_name=$2, last_name=$3, phone=$4,
			date_of_birth=$5, gender=$6, address=$7, program_id=$8,
			admission_year=$9, current_semester=$10, is_active=$11, updated_at=NOW()
		WHERE id=$12`,
		s.RollNumber, s.FirstName, s.LastName, s.Phone,
		s.DateOfBirth, s.Gender, s.Address, s.ProgramID,
		s.AdmissionYear, s.CurrentSemester, s.IsActive, id,
	)
	return err
}

func (r *StudentRepository) RollExists(roll string) (bool, error) {
	var c int
	err := DB.QueryRow(`SELECT COUNT(*) FROM students WHERE roll_number=$1`, roll).Scan(&c)
	return c > 0, err
}

func (r *StudentRepository) GetCourses(studentID int64) ([]models.StudentCourse, error) {
	rows, err := DB.Query(`
		SELECT r.id, r.course_offering_id,
			c.code, c.name, c.credits, c.course_type,
			s.name, co.academic_year,
			f.first_name||' '||f.last_name,
			COALESCE(g.grade,'')
		FROM registrations r
		JOIN course_offerings co ON co.id = r.course_offering_id
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id
		LEFT JOIN grades g ON g.registration_id = r.id
		WHERE r.student_id = $1
		ORDER BY s.academic_year DESC, s.semester_number ASC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.StudentCourse
	for rows.Next() {
		var sc models.StudentCourse
		if err := rows.Scan(&sc.RegistrationID, &sc.CourseOfferingID,
			&sc.CourseCode, &sc.CourseName, &sc.Credits, &sc.CourseType,
			&sc.SemesterName, &sc.AcademicYear, &sc.FacultyName, &sc.Grade); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

func (r *StudentRepository) GetGrades(studentID int64) ([]models.StudentCourse, error) {
	rows, err := DB.Query(`
		SELECT r.id, r.course_offering_id,
			c.code, c.name, c.credits, c.course_type,
			s.name, co.academic_year,
			f.first_name||' '||f.last_name,
			g.grade
		FROM registrations r
		JOIN course_offerings co ON co.id = r.course_offering_id
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id
		JOIN grades g ON g.registration_id = r.id
		WHERE r.student_id = $1
		ORDER BY s.academic_year DESC, s.semester_number ASC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.StudentCourse
	for rows.Next() {
		var sc models.StudentCourse
		if err := rows.Scan(&sc.RegistrationID, &sc.CourseOfferingID,
			&sc.CourseCode, &sc.CourseName, &sc.Credits, &sc.CourseType,
			&sc.SemesterName, &sc.AcademicYear, &sc.FacultyName, &sc.Grade); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}
