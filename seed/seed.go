package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	log.Println("Seeding database...")

	seedDepartments(db)
	seedPrograms(db)
	seedSemesters(db)
	seedCourses(db)
	seedAdmin(db)

	log.Println("Seed completed successfully.")
}

func hash(plain string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash: %v", err)
	}
	return string(b)
}

func seedDepartments(db *sql.DB) {
	depts := []struct{ name, code string }{
		{"Computer Science and Engineering", "CSE"},
		{"Information Technology", "IT"},
		{"Electronics and Communication Engineering", "ECE"},
		{"Electrical and Electronics Engineering", "EEE"},
		{"Mechanical Engineering", "MECH"},
		{"Civil Engineering", "CIVIL"},
	}
	for _, d := range depts {
		_, err := db.Exec(
			`INSERT INTO departments (name, code) VALUES ($1, $2) ON CONFLICT (code) DO NOTHING`,
			d.name, d.code)
		if err != nil {
			log.Printf("dept %s: %v", d.code, err)
		}
	}
	log.Println("  Departments seeded")
}

func seedPrograms(db *sql.DB) {
	programs := []struct {
		name, code, desc string
		duration         int
		deptCode         string
	}{
		{"B.Tech Computer Science and Engineering", "BTECH-CSE", "Bachelor of Technology in CSE", 4, "CSE"},
		{"B.Tech Information Technology", "BTECH-IT", "Bachelor of Technology in IT", 4, "IT"},
		{"B.Tech Electronics and Communication Engineering", "BTECH-ECE", "Bachelor of Technology in ECE", 4, "ECE"},
		{"B.Tech Electrical and Electronics Engineering", "BTECH-EEE", "Bachelor of Technology in EEE", 4, "EEE"},
		{"B.Tech Mechanical Engineering", "BTECH-MECH", "Bachelor of Technology in Mechanical", 4, "MECH"},
		{"B.Tech Civil Engineering", "BTECH-CIVIL", "Bachelor of Technology in Civil", 4, "CIVIL"},
	}
	for _, p := range programs {
		var deptID int64
		if err := db.QueryRow(`SELECT id FROM departments WHERE code = $1`, p.deptCode).Scan(&deptID); err != nil {
			log.Printf("dept not found %s: %v", p.deptCode, err)
			continue
		}
		_, err := db.Exec(`
			INSERT INTO programs (name, code, description, duration_years, department_id)
			VALUES ($1,$2,$3,$4,$5) ON CONFLICT (code) DO NOTHING`,
			p.name, p.code, p.desc, p.duration, deptID)
		if err != nil {
			log.Printf("program %s: %v", p.code, err)
		}
	}
	log.Println("  Programs seeded")
}

func seedSemesters(db *sql.DB) {
	now := time.Now()
	semesters := []struct {
		academicYear    string
		semNumber       int
		name            string
		regStart, regEnd time.Time
		isActive        bool
	}{
		{
			academicYear: "2026-2027", semNumber: 4,
			name:     "Fourth Semester",
			regStart: now.Add(-48 * time.Hour),
			regEnd:   now.Add(30 * 24 * time.Hour),
			isActive: true,
		},
		{"2025-2026", 3, "Third Semester", now.Add(-200 * 24 * time.Hour), now.Add(-100 * 24 * time.Hour), false},
		{"2025-2026", 2, "Second Semester", now.Add(-400 * 24 * time.Hour), now.Add(-300 * 24 * time.Hour), false},
		{"2024-2025", 1, "First Semester", now.Add(-600 * 24 * time.Hour), now.Add(-500 * 24 * time.Hour), false},
	}
	for _, s := range semesters {
		_, err := db.Exec(`
			INSERT INTO semesters (academic_year, semester_number, name, registration_start, registration_end, is_active)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT DO NOTHING`,
			s.academicYear, s.semNumber, s.name, s.regStart, s.regEnd, s.isActive)
		if err != nil {
			log.Printf("semester %s sem%d: %v", s.academicYear, s.semNumber, err)
		}
	}
	log.Println("  Semesters seeded")
}

func seedCourses(db *sql.DB) {
	var programID int64
	if err := db.QueryRow(`SELECT id FROM programs WHERE code = 'BTECH-CSE'`).Scan(&programID); err != nil {
		log.Printf("CSE program not found, skipping courses: %v", err)
		return
	}

	type courseRow struct {
		code, name, courseType string
		credits, semNum        int
	}

	courses := []courseRow{
		// Semester 1
		{"MA101", "Engineering Mathematics I", "Theory", 4, 1},
		{"CS101", "Programming in C", "Theory", 3, 1},
		{"PH101", "Engineering Physics", "Theory", 3, 1},
		{"CH101", "Engineering Chemistry", "Theory", 3, 1},
		{"EN101", "English", "Theory", 2, 1},
		{"CS102", "Programming Laboratory", "Laboratory", 2, 1},
		// Semester 2
		{"MA201", "Engineering Mathematics II", "Theory", 4, 2},
		{"CS201", "Data Structures", "Theory", 4, 2},
		{"CS202", "Digital Logic", "Theory", 3, 2},
		{"CS203", "Object Oriented Programming", "Theory", 3, 2},
		{"EG201", "Engineering Graphics", "Theory", 2, 2},
		{"CS204", "Data Structures Laboratory", "Laboratory", 2, 2},
		// Semester 3
		{"MA301", "Discrete Mathematics", "Theory", 4, 3},
		{"CS301", "Computer Organization", "Theory", 3, 3},
		{"CS302", "Database Fundamentals", "Theory", 3, 3},
		{"CS303", "Operating Systems Fundamentals", "Theory", 3, 3},
		{"CS304", "Java Programming", "Theory", 3, 3},
		{"CS305", "Database Laboratory", "Laboratory", 2, 3},
		// Semester 4
		{"CS401", "Database Management Systems", "Theory", 4, 4},
		{"CS402", "Operating Systems", "Theory", 4, 4},
		{"CS403", "Computer Networks", "Theory", 3, 4},
		{"MA401", "Engineering Mathematics", "Theory", 4, 4},
		{"CS404", "Operating Systems Laboratory", "Laboratory", 2, 4},
		// Semester 5
		{"CS501", "Software Engineering", "Theory", 4, 5},
		{"CS502", "Web Technologies", "Theory", 3, 5},
		{"CS503", "Theory of Computation", "Theory", 4, 5},
		{"CS504", "Artificial Intelligence", "Theory", 3, 5},
		{"CS505", "Cloud Computing", "Theory", 3, 5},
		// Semester 6
		{"CS601", "Machine Learning", "Theory", 4, 6},
		{"CS602", "Distributed Systems", "Theory", 3, 6},
		{"CS603", "Cyber Security", "Theory", 3, 6},
		{"CS604", "Mobile Application Development", "Theory", 3, 6},
		{"CS605", "Project Phase I", "Project", 2, 6},
		// Semester 7
		{"CS701", "Big Data Analytics", "Theory", 4, 7},
		{"CS702", "Advanced AI", "Theory", 3, 7},
		{"CS703", "DevOps", "Theory", 3, 7},
		{"CS704", "Project Phase II", "Project", 4, 7},
		{"CS705", "Professional Elective I", "Elective", 3, 7},
		// Semester 8
		{"CS801", "Major Project", "Project", 8, 8},
		{"CS802", "Professional Elective II", "Elective", 3, 8},
		{"CS803", "Internship / Seminar", "Elective", 4, 8},
	}

	for _, c := range courses {
		_, err := db.Exec(`
			INSERT INTO courses (code, name, course_type, credits, program_id, semester_number)
			VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (code) DO NOTHING`,
			c.code, c.name, c.courseType, c.credits, programID, c.semNum)
		if err != nil {
			log.Printf("course %s: %v", c.code, err)
		}
	}
	log.Println(fmt.Sprintf("  %d courses seeded for BTECH-CSE", len(courses)))
}

func seedAdmin(db *sql.DB) {
	passHash := hash("admin@123")
	var userID int64
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ('developer@gmail.com', $1, 'admin')
		ON CONFLICT (email) DO UPDATE SET password_hash=$1
		RETURNING id`, passHash).Scan(&userID)
	if err != nil {
		log.Fatalf("admin seed: %v", err)
	}
	log.Printf("  Admin seeded: developer@gmail.com / admin@123 (user_id=%d)", userID)
}
