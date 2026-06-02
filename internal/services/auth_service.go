package services

import (
	"errors"
	"log"

	"github.com/google/uuid"

	"user-management/internal/models"
	"user-management/internal/repositories"
	"user-management/internal/utils"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	sessionRepo      *repositories.SessionRepository
	loginHistoryRepo *repositories.LoginHistoryRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:         repositories.NewUserRepository(),
		sessionRepo:      repositories.NewSessionRepository(),
		loginHistoryRepo: repositories.NewLoginHistoryRepository(),
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Find user by email OR username
	var user *models.User
	var err error

	if containsAtSymbol(req.Email) {
		user, err = s.userRepo.GetByEmail(req.Email)
	} else {
		user, err = s.userRepo.GetByUsername(req.Email)
	}

	if err != nil {
		s.saveLoginHistory(nil, req.Email, "failed", ipAddress, userAgent, "internal server error")
		return nil, errors.New("internal server error")
	}

	if user == nil {
		s.saveLoginHistory(nil, req.Email, "failed", ipAddress, userAgent, "user not found")
		return nil, errors.New("invalid email/username or password")
	}

	// Check if user is soft deleted
	if user.DeletedAt != nil {
		s.saveLoginHistory(&user.ID, req.Email, "failed", ipAddress, userAgent, "account deleted")
		return nil, errors.New("account not found")
	}

	// Check if account is locked
	if user.IsLocked {
		s.saveLoginHistory(&user.ID, req.Email, "locked", ipAddress, userAgent, "account locked")
		return nil, errors.New("account is locked. please contact administrator")
	}

	// Check if account is active
	if !user.IsActive {
		s.saveLoginHistory(&user.ID, req.Email, "failed", ipAddress, userAgent, "account inactive")
		return nil, errors.New("account is inactive. please contact administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		s.userRepo.IncrementFailedLoginAttempts(user.ID)

		if user.FailedLoginAttempts+1 >= 5 {
			s.userRepo.LockAccount(user.ID)
			s.saveLoginHistory(&user.ID, req.Email, "locked", ipAddress, userAgent, "too many failed attempts")
			return nil, errors.New("account has been locked due to too many failed attempts")
		}
		s.saveLoginHistory(&user.ID, req.Email, "failed", ipAddress, userAgent, "invalid password")
		return nil, errors.New("invalid email/username or password")
	}

	// Reset failed login attempts
	s.userRepo.ResetFailedLoginAttempts(user.ID)

	// Generate JWT token
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Email, user.Role, user.Name)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		s.saveLoginHistory(&user.ID, req.Email, "failed", ipAddress, userAgent, "token generation failed")
		return nil, errors.New("failed to generate token")
	}

	// Save session
	tokenHash := utils.HashToken(token)
	session := &models.UserSession{
		UserID:    user.ID,
		TokenHash: tokenHash,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		log.Printf("Failed to create session: %v", err)
		return nil, errors.New("failed to create session")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID, ipAddress)

	// Save successful login history
	s.saveLoginHistory(&user.ID, req.Email, "success", ipAddress, userAgent, "")

	return &models.LoginResponse{
		Token:     token,
		User:      user.ToResponse(),
		ExpiresAt: expiresAt,
	}, nil
}

// saveLoginHistory saves login attempt record
func (s *AuthService) saveLoginHistory(userID *uuid.UUID, email, status, ipAddress, userAgent, failureReason string) {
	history := &models.LoginHistory{
		UserID:        userID,
		Email:         email,
		Status:        status,
		IPAddress:     &ipAddress,
		UserAgent:     &userAgent,
		FailureReason: nil,
	}
	if failureReason != "" {
		history.FailureReason = &failureReason
	}
	_ = s.loginHistoryRepo.Create(history)
}

func containsAtSymbol(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '@' {
			return true
		}
	}
	return false
}

// Logout invalidates a token
func (s *AuthService) Logout(tokenString string) error {
	tokenHash := utils.HashToken(tokenString)
	return s.sessionRepo.DeleteByTokenHash(tokenHash)
}

// GetUserByToken gets user info from token
func (s *AuthService) GetUserByToken(tokenString string) (*models.UserResponse, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	tokenHash := utils.HashToken(tokenString)
	session, err := s.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil || session == nil {
		return nil, errors.New("session not found or expired")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

// ValidateToken validates a token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	tokenHash := utils.HashToken(tokenString)
	session, err := s.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil || session == nil {
		return nil, errors.New("token not found or already logged out")
	}

	s.sessionRepo.UpdateLastActivity(tokenHash)

	return claims, nil
}
