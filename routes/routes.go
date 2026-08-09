package routes

import (
	"ams-backend/config"
	"ams-backend/controllers"
	"ams-backend/middleware"
	"ams-backend/repositories"
	"ams-backend/services"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	userRepo      := repositories.NewUserRepository()
	studentRepo   := repositories.NewStudentRepository()
	facultyRepo   := repositories.NewFacultyRepository()
	deptRepo      := repositories.NewDepartmentRepository()
	programRepo   := repositories.NewProgramRepository()
	semesterRepo  := repositories.NewSemesterRepository()
	courseRepo    := repositories.NewCourseRepository()
	gradeRepo     := repositories.NewGradeRepository()
	emailRepo     := repositories.NewEmailSettingsRepository()
	adminRepo     := repositories.NewAdminRepository()

	authSvc  := services.NewAuthService(userRepo)
	emailSvc := services.NewEmailService(emailRepo)
	gradeSvc := services.NewGradeService(gradeRepo, courseRepo, studentRepo, facultyRepo, emailSvc)
	adminSvc := services.NewAdminService(db, userRepo, studentRepo, facultyRepo)

	authCtrl    := controllers.NewAuthController(authSvc, userRepo)
	adminCtrl   := controllers.NewAdminController(adminRepo, studentRepo, facultyRepo, deptRepo,
		programRepo, semesterRepo, courseRepo, emailRepo, adminSvc, emailSvc)
	facultyCtrl := controllers.NewFacultyController(facultyRepo, courseRepo, gradeSvc)
	studentCtrl := controllers.NewStudentController(studentRepo, semesterRepo, courseRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.App.FrontendURL, "http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	// Public
	api.POST("/auth/login", authCtrl.Login)

	// Any authenticated user
	api.GET("/auth/me", middleware.Auth(), authCtrl.Me)

	// Admin routes — middleware applied directly on the group
	admin := api.Group("/admin", middleware.Auth(), middleware.RequireRole("admin"))
	admin.GET("/dashboard", adminCtrl.GetDashboard)

	admin.GET("/students",     adminCtrl.ListStudents)
	admin.GET("/students/:id", adminCtrl.GetStudent)
	admin.POST("/students",    adminCtrl.CreateStudent)
	admin.PUT("/students/:id", adminCtrl.UpdateStudent)

	admin.GET("/faculty",     adminCtrl.ListFaculty)
	admin.GET("/faculty/:id", adminCtrl.GetFaculty)
	admin.POST("/faculty",    adminCtrl.CreateFaculty)
	admin.PUT("/faculty/:id", adminCtrl.UpdateFaculty)

	admin.GET("/departments", adminCtrl.ListDepartments)

	admin.GET("/programs",     adminCtrl.ListPrograms)
	admin.POST("/programs",    adminCtrl.CreateProgram)
	admin.PUT("/programs/:id", adminCtrl.UpdateProgram)

	admin.GET("/semesters",     adminCtrl.ListSemesters)
	admin.POST("/semesters",    adminCtrl.CreateSemester)
	admin.PUT("/semesters/:id", adminCtrl.UpdateSemester)

	admin.GET("/courses",     adminCtrl.ListCourses)
	admin.POST("/courses",    adminCtrl.CreateCourse)
	admin.PUT("/courses/:id", adminCtrl.UpdateCourse)

	admin.GET("/course-offerings",     adminCtrl.ListCourseOfferings)
	admin.POST("/course-offerings",    adminCtrl.CreateCourseOffering)
	admin.PUT("/course-offerings/:id", adminCtrl.UpdateCourseOffering)

	admin.GET("/settings/email",       adminCtrl.GetEmailSettings)
	admin.PUT("/settings/email",       adminCtrl.UpdateEmailSettings)
	admin.POST("/settings/email/test", adminCtrl.TestEmail)

	// Faculty routes
	faculty := api.Group("/faculty", middleware.Auth(), middleware.RequireRole("faculty"))
	faculty.GET("/me",      facultyCtrl.GetProfile)
	faculty.GET("/courses", facultyCtrl.GetCurrentCourses)
	faculty.GET("/courses/history", facultyCtrl.GetCourseHistory)
	faculty.GET("/courses/:courseOfferingId/students",    facultyCtrl.GetCourseStudents)
	faculty.POST("/courses/:courseOfferingId/grades",     facultyCtrl.UploadGrade)
	faculty.PUT("/courses/:courseOfferingId/grades/:sid", facultyCtrl.UploadGrade)
	faculty.POST("/courses/:courseOfferingId/grades/bulk", facultyCtrl.BulkUploadGrades)

	// Student routes
	student := api.Group("", middleware.Auth(), middleware.RequireRole("student"))
	student.GET("/students/me",         studentCtrl.GetProfile)
	student.GET("/students/me/courses", studentCtrl.GetCourses)
	student.GET("/students/me/grades",  studentCtrl.GetGrades)
	student.GET("/semesters/active",    studentCtrl.GetActiveSemester)
	student.GET("/courses/available",   studentCtrl.GetAvailableCourses)
	student.POST("/courses/register",   studentCtrl.RegisterCourse)

	return r
}
