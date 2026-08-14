package controllers

import (
	"ams-backend/middleware"
	"ams-backend/repositories"
	"ams-backend/services"
	"ams-backend/utils"
	"strings"

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
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Fail(c, 400, "Please provide email and password")
		return
	}

	body.Email    = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)

	if body.Email == "" {
		utils.Fail(c, 400, "Email address is required")
		return
	}
	if body.Password == "" {
		utils.Fail(c, 400, "Password is required")
		return
	}

	result, err := ctrl.authSvc.Login(body.Email, body.Password)
	if err != nil {
		// Use 400 not 401 so frontend api.js does not auto-redirect to /login
		utils.Fail(c, 400, "Invalid email or password. Please check your credentials.")
		return
	}

	if !result.User.IsActive {
		utils.Fail(c, 400, "Your account has been deactivated. Please contact the administrator.")
		return
	}

	utils.OK(c, result)
}

func (ctrl *AuthController) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil || user == nil {
		utils.Fail(c, 404, "User not found")
		return
	}
	utils.OK(c, user)
}
