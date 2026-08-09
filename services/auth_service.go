package services

import (
	"ams-backend/config"
	"ams-backend/models"
	"ams-backend/repositories"
	"ams-backend/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(r *repositories.UserRepository) *AuthService { return &AuthService{userRepo: r} }

type LoginResult struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Login(email, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if user == nil || !user.IsActive {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	expiry := time.Duration(config.App.JWTExpiry) * time.Hour
	token, err := utils.GenerateJWT(user.ID, user.Role, expiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &LoginResult{Token: token, User: user}, nil
}

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}
