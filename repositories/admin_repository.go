package repositories

import "ams-backend/models"

type AdminRepository struct{}

func NewAdminRepository() *AdminRepository { return &AdminRepository{} }

func (r *AdminRepository) GetDashboardStats() (*models.DashboardStats, error) {
	var stats models.DashboardStats
	err := DB.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM students WHERE is_active = true),
			(SELECT COUNT(*) FROM faculty  WHERE is_active = true),
			(SELECT COUNT(*) FROM programs WHERE is_active = true),
			(SELECT COUNT(*) FROM courses  WHERE is_active = true),
			(SELECT COUNT(*) FROM course_offerings WHERE is_active = true),
			(SELECT COUNT(*) FROM registrations),
			COALESCE((SELECT name FROM semesters WHERE is_active = true LIMIT 1), 'None')
	`).Scan(
		&stats.TotalStudents, &stats.TotalFaculty, &stats.TotalPrograms,
		&stats.TotalCourses, &stats.TotalCourseOfferings, &stats.TotalRegistrations,
		&stats.ActiveSemesterName,
	)
	return &stats, err
}
