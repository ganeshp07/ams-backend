package models

type DashboardStats struct {
	TotalStudents        int    `json:"total_students"`
	TotalFaculty         int    `json:"total_faculty"`
	TotalPrograms        int    `json:"total_programs"`
	TotalCourses         int    `json:"total_courses"`
	TotalCourseOfferings int    `json:"total_course_offerings"`
	TotalRegistrations   int    `json:"total_registrations"`
	ActiveSemesterName   string `json:"active_semester_name"`
}
