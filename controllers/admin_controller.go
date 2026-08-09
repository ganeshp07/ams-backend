package controllers

import (
	"ams-backend/middleware"
	"ams-backend/models"
	"ams-backend/repositories"
	"ams-backend/services"
	"ams-backend/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminRepo    *repositories.AdminRepository
	studentRepo  *repositories.StudentRepository
	facultyRepo  *repositories.FacultyRepository
	deptRepo     *repositories.DepartmentRepository
	programRepo  *repositories.ProgramRepository
	semesterRepo *repositories.SemesterRepository
	courseRepo   *repositories.CourseRepository
	emailRepo    *repositories.EmailSettingsRepository
	adminSvc     *services.AdminService
	emailSvc     *services.EmailService
}

func NewAdminController(
	ar *repositories.AdminRepository,
	sr *repositories.StudentRepository,
	fr *repositories.FacultyRepository,
	dr *repositories.DepartmentRepository,
	pr *repositories.ProgramRepository,
	semr *repositories.SemesterRepository,
	cr *repositories.CourseRepository,
	er *repositories.EmailSettingsRepository,
	as *services.AdminService,
	es *services.EmailService,
) *AdminController {
	return &AdminController{ar, sr, fr, dr, pr, semr, cr, er, as, es}
}

// ── Dashboard ──────────────────────────────────────────────────────────────

func (ctrl *AdminController) GetDashboard(c *gin.Context) {
	stats, err := ctrl.adminRepo.GetDashboardStats()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch stats")
		return
	}
	utils.OK(c, stats)
}

// ── Students ───────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListStudents(c *gin.Context) {
	search := c.Query("search")
	var programID int64
	if pid := c.Query("program_id"); pid != "" {
		programID, _ = strconv.ParseInt(pid, 10, 64)
	}
	var isActive *bool
	if a := c.Query("is_active"); a != "" {
		v := a == "true"
		isActive = &v
	}

	students, err := ctrl.studentRepo.List(search, programID, isActive)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch students")
		return
	}
	if students == nil {
		students = []models.StudentProfile{}
	}
	utils.OK(c, students)
}

func (ctrl *AdminController) GetStudent(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	profile, err := ctrl.studentRepo.FindByID(id)
	if err != nil || profile == nil {
		utils.Fail(c, 404, "student not found")
		return
	}
	utils.OK(c, profile)
}

func (ctrl *AdminController) CreateStudent(c *gin.Context) {
	var input services.CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if strings.TrimSpace(input.RollNumber) == "" || strings.TrimSpace(input.Email) == "" {
		utils.Fail(c, 400, "roll number and email are required")
		return
	}

	createdBy := middleware.GetUserID(c)
	profile, err := ctrl.adminSvc.CreateStudent(&input, createdBy)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Created(c, profile)
}

func (ctrl *AdminController) UpdateStudent(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var s models.Student
	if err := c.ShouldBindJSON(&s); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.studentRepo.Update(id, &s); err != nil {
		utils.Fail(c, 500, "failed to update student")
		return
	}
	profile, _ := ctrl.studentRepo.FindByID(id)
	utils.OK(c, profile)
}

// ── Faculty ────────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListFaculty(c *gin.Context) {
	search := c.Query("search")
	list, err := ctrl.facultyRepo.List(search)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch faculty")
		return
	}
	if list == nil {
		list = []models.FacultyProfile{}
	}
	utils.OK(c, list)
}

func (ctrl *AdminController) GetFaculty(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	profile, err := ctrl.facultyRepo.FindByID(id)
	if err != nil || profile == nil {
		utils.Fail(c, 404, "faculty not found")
		return
	}
	utils.OK(c, profile)
}

func (ctrl *AdminController) CreateFaculty(c *gin.Context) {
	var input services.CreateFacultyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	createdBy := middleware.GetUserID(c)
	profile, err := ctrl.adminSvc.CreateFaculty(&input, createdBy)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Created(c, profile)
}

func (ctrl *AdminController) UpdateFaculty(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var f models.Faculty
	if err := c.ShouldBindJSON(&f); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.facultyRepo.Update(id, &f); err != nil {
		utils.Fail(c, 500, "failed to update faculty")
		return
	}
	profile, _ := ctrl.facultyRepo.FindByID(id)
	utils.OK(c, profile)
}

// ── Departments ────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListDepartments(c *gin.Context) {
	depts, err := ctrl.deptRepo.List()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch departments")
		return
	}
	if depts == nil {
		depts = []models.Department{}
	}
	utils.OK(c, depts)
}

// ── Programs ───────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListPrograms(c *gin.Context) {
	list, err := ctrl.programRepo.List()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch programs")
		return
	}
	if list == nil {
		list = []models.Program{}
	}
	utils.OK(c, list)
}

func (ctrl *AdminController) CreateProgram(c *gin.Context) {
	var p models.Program
	if err := c.ShouldBindJSON(&p); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	id, err := ctrl.programRepo.Create(&p)
	if err != nil {
		utils.Fail(c, 500, "failed to create program")
		return
	}
	created, _ := ctrl.programRepo.FindByID(id)
	utils.Created(c, created)
}

func (ctrl *AdminController) UpdateProgram(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var p models.Program
	if err := c.ShouldBindJSON(&p); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.programRepo.Update(id, &p); err != nil {
		utils.Fail(c, 500, "failed to update program")
		return
	}
	updated, _ := ctrl.programRepo.FindByID(id)
	utils.OK(c, updated)
}

// ── Semesters ──────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListSemesters(c *gin.Context) {
	list, err := ctrl.semesterRepo.List()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch semesters")
		return
	}
	if list == nil {
		list = []models.Semester{}
	}
	utils.OK(c, list)
}

func (ctrl *AdminController) CreateSemester(c *gin.Context) {
	var s models.Semester
	if err := c.ShouldBindJSON(&s); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	id, err := ctrl.semesterRepo.Create(&s)
	if err != nil {
		utils.Fail(c, 500, "failed to create semester")
		return
	}
	created, _ := ctrl.semesterRepo.FindByID(id)
	utils.Created(c, created)
}

func (ctrl *AdminController) UpdateSemester(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var s models.Semester
	if err := c.ShouldBindJSON(&s); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.semesterRepo.Update(id, &s); err != nil {
		utils.Fail(c, 500, "failed to update semester")
		return
	}
	updated, _ := ctrl.semesterRepo.FindByID(id)
	utils.OK(c, updated)
}

// ── Courses ────────────────────────────────────────────────────────────────

func (ctrl *AdminController) ListCourses(c *gin.Context) {
	var programID int64
	if pid := c.Query("program_id"); pid != "" {
		programID, _ = strconv.ParseInt(pid, 10, 64)
	}
	list, err := ctrl.courseRepo.List(programID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch courses")
		return
	}
	if list == nil {
		list = []models.Course{}
	}
	utils.OK(c, list)
}

func (ctrl *AdminController) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	id, err := ctrl.courseRepo.Create(&course)
	if err != nil {
		utils.Fail(c, 500, "failed to create course")
		return
	}
	utils.Created(c, gin.H{"id": id})
}

func (ctrl *AdminController) UpdateCourse(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.courseRepo.Update(id, &course); err != nil {
		utils.Fail(c, 500, "failed to update course")
		return
	}
	utils.OKMsg(c, "course updated", gin.H{"id": id})
}

// ── Course Offerings ───────────────────────────────────────────────────────

func (ctrl *AdminController) ListCourseOfferings(c *gin.Context) {
	var semesterID int64
	if sid := c.Query("semester_id"); sid != "" {
		semesterID, _ = strconv.ParseInt(sid, 10, 64)
	}
	list, err := ctrl.courseRepo.ListOfferings(semesterID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch course offerings")
		return
	}
	if list == nil {
		list = []models.CourseOffering{}
	}
	utils.OK(c, list)
}

func (ctrl *AdminController) CreateCourseOffering(c *gin.Context) {
	var co models.CourseOffering
	if err := c.ShouldBindJSON(&co); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if co.Section == "" {
		co.Section = "A"
	}
	id, err := ctrl.courseRepo.CreateOffering(&co)
	if err != nil {
		utils.Fail(c, 500, "failed to create course offering")
		return
	}
	created, _ := ctrl.courseRepo.GetOfferingByID(id)
	utils.Created(c, created)
}

func (ctrl *AdminController) UpdateCourseOffering(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}
	var co models.CourseOffering
	if err := c.ShouldBindJSON(&co); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}
	if err := ctrl.courseRepo.UpdateOffering(id, &co); err != nil {
		utils.Fail(c, 500, "failed to update course offering")
		return
	}
	updated, _ := ctrl.courseRepo.GetOfferingByID(id)
	utils.OK(c, updated)
}

// ── Email Settings ─────────────────────────────────────────────────────────

func (ctrl *AdminController) GetEmailSettings(c *gin.Context) {
	settings, err := ctrl.emailRepo.Get()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch email settings")
		return
	}
	// Mask the password before returning
	masked := "••••••••••••"
	if settings.SmtpPassword == "" {
		masked = ""
	}
	resp := models.EmailSettingsResponse{
		ID:           settings.ID,
		SmtpHost:     settings.SmtpHost,
		SmtpPort:     settings.SmtpPort,
		SmtpUsername: settings.SmtpUsername,
		SmtpPassword: masked,
		FromEmail:    settings.FromEmail,
		FromName:     settings.FromName,
		IsEnabled:    settings.IsEnabled,
		UpdatedAt:    settings.UpdatedAt,
	}
	utils.OK(c, resp)
}

func (ctrl *AdminController) UpdateEmailSettings(c *gin.Context) {
	var body struct {
		SmtpHost     string `json:"smtp_host"`
		SmtpPort     int    `json:"smtp_port"`
		SmtpUsername string `json:"smtp_username"`
		SmtpPassword string `json:"smtp_password"` // empty = keep existing
		FromEmail    string `json:"from_email"`
		FromName     string `json:"from_name"`
		IsEnabled    bool   `json:"is_enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Fail(c, 400, "invalid request body")
		return
	}

	// Load existing to get ID
	existing, err := ctrl.emailRepo.Get()
	if err != nil {
		utils.Fail(c, 500, "failed to load email settings")
		return
	}

	updated := &models.EmailSettings{
		ID:           existing.ID,
		SmtpHost:     body.SmtpHost,
		SmtpPort:     body.SmtpPort,
		SmtpUsername: body.SmtpUsername,
		FromEmail:    body.FromEmail,
		FromName:     body.FromName,
		IsEnabled:    body.IsEnabled,
	}
	// Only update password if a real new value was sent (not the masked placeholder)
	if body.SmtpPassword != "" && body.SmtpPassword != "••••••••••••" {
		updated.SmtpPassword = body.SmtpPassword
	}

	if err := ctrl.emailRepo.Update(updated); err != nil {
		utils.Fail(c, 500, "failed to update email settings")
		return
	}
	utils.OKMsg(c, "email settings updated", nil)
}

func (ctrl *AdminController) TestEmail(c *gin.Context) {
	var body struct {
		To string `json:"to"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.To == "" {
		utils.Fail(c, 400, "recipient email (to) is required")
		return
	}
	if err := ctrl.emailSvc.TestConnection(body.To); err != nil {
		utils.Fail(c, 500, "test email failed: "+err.Error())
		return
	}
	utils.OKMsg(c, "test email sent to "+body.To, nil)
}

// parseID reads a URL path parameter as int64.
func parseID(c *gin.Context, param string) int64 {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil || id <= 0 {
		utils.Fail(c, 400, "invalid "+param)
		return 0
	}
	return id
}
