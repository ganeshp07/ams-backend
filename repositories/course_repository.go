package repositories

import (
	"ams-backend/models"
	"database/sql"
)

type CourseRepository struct{}

func NewCourseRepository() *CourseRepository { return &CourseRepository{} }

func (r *CourseRepository) List(programID int64) ([]models.Course, error) {
	args := []interface{}{}
	where := ""
	if programID > 0 {
		where = " WHERE c.program_id = $1"
		args = append(args, programID)
	}
	rows, err := DB.Query(`
		SELECT c.id, c.code, c.name, COALESCE(c.description,''), c.credits, c.course_type,
			c.program_id, COALESCE(p.name,''), c.semester_number, c.is_active, c.created_at, c.updated_at
		FROM courses c LEFT JOIN programs p ON p.id = c.program_id`+where+` ORDER BY c.semester_number, c.code`,
		args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Description, &c.Credits, &c.CourseType,
			&c.ProgramID, &c.ProgramName, &c.SemesterNumber, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CourseRepository) Create(c *models.Course) (int64, error) {
	var id int64
	err := DB.QueryRow(`
		INSERT INTO courses (code, name, description, credits, course_type, program_id, semester_number)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		c.Code, c.Name, c.Description, c.Credits, c.CourseType, c.ProgramID, c.SemesterNumber,
	).Scan(&id)
	return id, err
}

func (r *CourseRepository) Update(id int64, c *models.Course) error {
	_, err := DB.Exec(`
		UPDATE courses SET code=$1, name=$2, description=$3, credits=$4, course_type=$5,
			program_id=$6, semester_number=$7, is_active=$8, updated_at=NOW() WHERE id=$9`,
		c.Code, c.Name, c.Description, c.Credits, c.CourseType,
		c.ProgramID, c.SemesterNumber, c.IsActive, id)
	return err
}

func (r *CourseRepository) CodeExists(code string) (bool, error) {
	var c int
	err := DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE code=$1`, code).Scan(&c)
	return c > 0, err
}

// ListOfferings lists all course offerings (admin view).
func (r *CourseRepository) ListOfferings(semesterID int64) ([]models.CourseOffering, error) {
	args := []interface{}{}
	where := ""
	if semesterID > 0 {
		where = " WHERE co.semester_id = $1"
		args = append(args, semesterID)
	}
	rows, err := DB.Query(`
		SELECT co.id, co.course_id, c.code, c.name, c.credits, c.course_type,
			co.semester_id, s.name, co.academic_year,
			co.faculty_id, f.first_name||' '||f.last_name,
			co.section, co.is_active, co.created_at
		FROM course_offerings co
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id`+where+` ORDER BY s.academic_year DESC, c.code`,
		args...)
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

func (r *CourseRepository) CreateOffering(co *models.CourseOffering) (int64, error) {
	var id int64
	err := DB.QueryRow(`
		INSERT INTO course_offerings (course_id, semester_id, faculty_id, academic_year, section)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		co.CourseID, co.SemesterID, co.FacultyID, co.AcademicYear, co.Section,
	).Scan(&id)
	return id, err
}

func (r *CourseRepository) UpdateOffering(id int64, co *models.CourseOffering) error {
	_, err := DB.Exec(`
		UPDATE course_offerings SET course_id=$1, semester_id=$2, faculty_id=$3,
			academic_year=$4, section=$5, is_active=$6, updated_at=NOW() WHERE id=$7`,
		co.CourseID, co.SemesterID, co.FacultyID, co.AcademicYear, co.Section, co.IsActive, id)
	return err
}

func (r *CourseRepository) GetOfferingByID(id int64) (*models.CourseOffering, error) {
	var co models.CourseOffering
	err := DB.QueryRow(`
		SELECT co.id, co.course_id, c.code, c.name, c.credits, c.course_type,
			co.semester_id, s.name, co.academic_year,
			co.faculty_id, f.first_name||' '||f.last_name,
			co.section, co.is_active, co.created_at
		FROM course_offerings co
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id
		WHERE co.id=$1`, id,
	).Scan(&co.ID, &co.CourseID, &co.CourseCode, &co.CourseName,
		&co.Credits, &co.CourseType, &co.SemesterID, &co.SemesterName,
		&co.AcademicYear, &co.FacultyID, &co.FacultyName,
		&co.Section, &co.IsActive, &co.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &co, err
}

func (r *CourseRepository) GetAvailableForStudent(studentID, semesterID int64) ([]models.CourseOffering, error) {
	rows, err := DB.Query(`
		SELECT co.id, co.course_id, c.code, c.name, c.credits, c.course_type,
			co.semester_id, s.name, co.academic_year,
			co.faculty_id, f.first_name||' '||f.last_name,
			co.section, co.is_active, co.created_at
		FROM course_offerings co
		JOIN courses c ON c.id = co.course_id
		JOIN semesters s ON s.id = co.semester_id
		JOIN faculty f ON f.id = co.faculty_id
		WHERE co.semester_id = $1 AND co.is_active = true
		  AND co.id NOT IN (
			SELECT course_offering_id FROM registrations WHERE student_id = $2
		)
		ORDER BY c.code`, semesterID, studentID)
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

func (r *CourseRepository) IsRegistered(studentID, offeringID int64) (bool, error) {
	var c int
	err := DB.QueryRow(`SELECT COUNT(*) FROM registrations WHERE student_id=$1 AND course_offering_id=$2`,
		studentID, offeringID).Scan(&c)
	return c > 0, err
}

func (r *CourseRepository) RegisterStudent(studentID, offeringID int64) (int64, error) {
	var id int64
	err := DB.QueryRow(`
		INSERT INTO registrations (student_id, course_offering_id) VALUES ($1,$2) RETURNING id`,
		studentID, offeringID).Scan(&id)
	return id, err
}

func (r *CourseRepository) GetStudentsInOffering(offeringID int64) ([]models.CourseStudent, error) {
	rows, err := DB.Query(`
		SELECT s.id, s.roll_number, s.first_name||' '||s.last_name, u.email,
			COALESCE(g.grade,''), g.id, r.id
		FROM registrations r
		JOIN students s ON s.id = r.student_id
		JOIN users u ON u.id = s.user_id
		LEFT JOIN grades g ON g.registration_id = r.id
		WHERE r.course_offering_id = $1
		ORDER BY s.roll_number`, offeringID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CourseStudent
	for rows.Next() {
		var cs models.CourseStudent
		if err := rows.Scan(&cs.StudentID, &cs.RollNumber, &cs.FullName, &cs.Email,
			&cs.Grade, &cs.GradeID, &cs.RegistrationID); err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, rows.Err()
}

func (r *CourseRepository) GetRegistrationForStudent(studentID, offeringID int64) (int64, error) {
	var id int64
	err := DB.QueryRow(`SELECT id FROM registrations WHERE student_id=$1 AND course_offering_id=$2`,
		studentID, offeringID).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}
