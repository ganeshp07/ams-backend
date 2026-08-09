package controllers

import (
	"ams-backend/middleware"
	"ams-backend/models"
	"ams-backend/repositories"
	"ams-backend/services"
	"ams-backend/utils"

	"github.com/gin-gonic/gin"
)

type FacultyController struct {
	facultyRepo  *repositories.FacultyRepository
	courseRepo   *repositories.CourseRepository
	gradeSvc     *services.GradeService
}

func NewFacultyController(fr *repositories.FacultyRepository, cr *repositories.CourseRepository, gs *services.GradeService) *FacultyController {
	return &FacultyController{fr, cr, gs}
}

func (ctrl *FacultyController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := ctrl.facultyRepo.FindByUserID(userID)
	if err != nil || profile == nil {
		utils.Fail(c, 404, "faculty profile not found")
		return
	}
	utils.OK(c, profile)
}

func (ctrl *FacultyController) GetCurrentCourses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	faculty, err := ctrl.facultyRepo.FindByUserID(userID)
	if err != nil || faculty == nil {
		utils.Fail(c, 404, "faculty not found")
		return
	}
	courses, err := ctrl.facultyRepo.GetCurrentCourses(faculty.ID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch courses")
		return
	}
	if courses == nil {
		courses = []models.CourseOffering{}
	}
	utils.OK(c, courses)
}

func (ctrl *FacultyController) GetCourseHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	faculty, err := ctrl.facultyRepo.FindByUserID(userID)
	if err != nil || faculty == nil {
		utils.Fail(c, 404, "faculty not found")
		return
	}
	history, err := ctrl.facultyRepo.GetCourseHistory(faculty.ID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch course history")
		return
	}
	if history == nil {
		history = []models.CourseOffering{}
	}
	utils.OK(c, history)
}

func (ctrl *FacultyController) GetCourseStudents(c *gin.Context) {
	offeringID := parseID(c, "courseOfferingId")
	if offeringID == 0 {
		return
	}

	// Verify faculty teaches this course
	userID := middleware.GetUserID(c)
	faculty, err := ctrl.facultyRepo.FindByUserID(userID)
	if err != nil || faculty == nil {
		utils.Fail(c, 404, "faculty not found")
		return
	}
	ok, err := ctrl.facultyRepo.TeachesCourseOffering(faculty.ID, offeringID)
	if err != nil || !ok {
		utils.Fail(c, 403, "you are not assigned to this course")
		return
	}

	students, err := ctrl.courseRepo.GetStudentsInOffering(offeringID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch students")
		return
	}
	if students == nil {
		students = []models.CourseStudent{}
	}
	utils.OK(c, students)
}

func (ctrl *FacultyController) UploadGrade(c *gin.Context) {
	offeringID := parseID(c, "courseOfferingId")
	if offeringID == 0 {
		return
	}

	var body struct {
		StudentID int64  `json:"student_id"`
		Grade     string `json:"grade"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.StudentID == 0 || body.Grade == "" {
		utils.Fail(c, 400, "student_id and grade are required")
		return
	}

	userID := middleware.GetUserID(c)
	grade, err := ctrl.gradeSvc.UploadGrade(userID, offeringID, body.StudentID, body.Grade)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.OKMsg(c, "grade saved", grade)
}

func (ctrl *FacultyController) BulkUploadGrades(c *gin.Context) {
	offeringID := parseID(c, "courseOfferingId")
	if offeringID == 0 {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, 400, "no file uploaded; field name must be 'file'")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		utils.Fail(c, 500, "failed to open uploaded file")
		return
	}
	defer file.Close()

	userID := middleware.GetUserID(c)
	result, err := ctrl.gradeSvc.BulkUpload(userID, offeringID, file)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.OKMsg(c, "bulk upload completed", result)
}
