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

// ChangePassword verifies the current password then stores a new bcrypt hash.
func (s *AuthService) ChangePassword(userID int64, currentPwd, newPwd string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPwd)); err != nil {
		return errors.New("current password is incorrect")
	}
	if len(newPwd) < 6 {
		return errors.New("new password must be at least 6 characters")
	}
	hash, err := HashPassword(newPwd)
	if err != nil {
		return errors.New("failed to hash password")
	}
	return s.userRepo.UpdatePassword(userID, hash)
}

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}
