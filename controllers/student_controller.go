package controllers

import (
	"ams-backend/middleware"
	"ams-backend/models"
	"ams-backend/repositories"
	"ams-backend/utils"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	studentRepo  *repositories.StudentRepository
	semesterRepo *repositories.SemesterRepository
	courseRepo   *repositories.CourseRepository
}

func NewStudentController(sr *repositories.StudentRepository, semr *repositories.SemesterRepository, cr *repositories.CourseRepository) *StudentController {
	return &StudentController{sr, semr, cr}
}

func (ctrl *StudentController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := ctrl.studentRepo.FindByUserID(userID)
	if err != nil || profile == nil {
		utils.Fail(c, 404, "student profile not found")
		return
	}
	utils.OK(c, profile)
}

func (ctrl *StudentController) GetCourses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	student, err := ctrl.studentRepo.FindByUserID(userID)
	if err != nil || student == nil {
		utils.Fail(c, 404, "student not found")
		return
	}
	courses, err := ctrl.studentRepo.GetCourses(student.ID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch courses")
		return
	}
	if courses == nil {
		courses = []models.StudentCourse{}
	}
	utils.OK(c, courses)
}

func (ctrl *StudentController) GetGrades(c *gin.Context) {
	userID := middleware.GetUserID(c)
	student, err := ctrl.studentRepo.FindByUserID(userID)
	if err != nil || student == nil {
		utils.Fail(c, 404, "student not found")
		return
	}
	grades, err := ctrl.studentRepo.GetGrades(student.ID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch grades")
		return
	}
	if grades == nil {
		grades = []models.StudentCourse{}
	}
	utils.OK(c, grades)
}

func (ctrl *StudentController) GetActiveSemester(c *gin.Context) {
	semester, err := ctrl.semesterRepo.GetActive()
	if err != nil {
		utils.Fail(c, 500, "failed to fetch active semester")
		return
	}
	if semester == nil {
		utils.Fail(c, 404, "no active semester")
		return
	}
	utils.OK(c, semester)
}

func (ctrl *StudentController) GetAvailableCourses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	student, err := ctrl.studentRepo.FindByUserID(userID)
	if err != nil || student == nil {
		utils.Fail(c, 404, "student not found")
		return
	}
	semester, err := ctrl.semesterRepo.GetActive()
	if err != nil || semester == nil {
		utils.Fail(c, 404, "no active semester")
		return
	}
	courses, err := ctrl.courseRepo.GetAvailableForStudent(student.ID, semester.ID)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch available courses")
		return
	}
	if courses == nil {
		courses = []models.CourseOffering{}
	}
	utils.OK(c, courses)
}

func (ctrl *StudentController) RegisterCourse(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body struct {
		CourseOfferingID int64 `json:"course_offering_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.CourseOfferingID == 0 {
		utils.Fail(c, 400, "course_offering_id is required")
		return
	}

	student, err := ctrl.studentRepo.FindByUserID(userID)
	if err != nil || student == nil {
		utils.Fail(c, 404, "student not found")
		return
	}

	offering, err := ctrl.courseRepo.GetOfferingByID(body.CourseOfferingID)
	if err != nil || offering == nil {
		utils.Fail(c, 404, "course offering not found")
		return
	}

	semester, err := ctrl.semesterRepo.FindByID(offering.SemesterID)
	if err != nil || semester == nil {
		utils.Fail(c, 404, "semester not found")
		return
	}

	// Registration window is enforced here, not just on the frontend.
	if !semester.IsRegistrationOpen() {
		utils.Fail(c, 400, "course registration is not open for this semester")
		return
	}

	already, _ := ctrl.courseRepo.IsRegistered(student.ID, body.CourseOfferingID)
	if already {
		utils.Fail(c, 409, "already registered for this course")
		return
	}

	regID, err := ctrl.courseRepo.RegisterStudent(student.ID, body.CourseOfferingID)
	if err != nil {
		utils.Fail(c, 500, "registration failed")
		return
	}
	utils.Created(c, gin.H{"registration_id": regID})
}
