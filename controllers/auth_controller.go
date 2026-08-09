package controllers

import (
	"ams-backend/middleware"
	"ams-backend/repositories"
	"ams-backend/services"
	"ams-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authSvc  *services.AuthService
	userRepo *repositories.UserRepository
}

func NewAuthController(a *services.AuthService, u *repositories.UserRepository) *AuthController {
	return &AuthController{a, u}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" || body.Password == "" {
		utils.Fail(c, 400, "email and password are required")
		return
	}

	result, err := ctrl.authSvc.Login(body.Email, body.Password)
	if err != nil {
		utils.Fail(c, 401, err.Error())
		return
	}
	utils.OK(c, result)
}

func (ctrl *AuthController) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil || user == nil {
		utils.Fail(c, 404, "user not found")
		return
	}
	utils.OK(c, user)
}
