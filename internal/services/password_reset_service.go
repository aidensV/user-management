package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"user-management/internal/models"
	"user-management/internal/repositories"
	"user-management/internal/utils"
)

type PasswordResetService struct {
	resetRepo *repositories.PasswordResetRepository
	userRepo  *repositories.UserRepository
}

func NewPasswordResetService() *PasswordResetService {
	return &PasswordResetService{
		resetRepo: repositories.NewPasswordResetRepository(),
		userRepo:  repositories.NewUserRepository(),
	}
}

// GenerateResetToken generates a password reset token for a user
func (s *PasswordResetService) GenerateResetToken(email string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Don't reveal that user doesn't exist for security
		return "", nil
	}

	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	// Create reset token record
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    false,
	}

	if err := s.resetRepo.Create(resetToken); err != nil {
		return "", err
	}

	// Send email (optional, can be disabled in development)
	// resetURL := os.Getenv("RESET_PASSWORD_URL")
	// if resetURL != "" {
	//     _ = utils.SendResetPasswordEmail(user.Email, token, resetURL)
	// }

	return token, nil
}

// ValidateResetToken validates a reset token
func (s *PasswordResetService) ValidateResetToken(token string) (*models.PasswordResetToken, error) {
	resetToken, err := s.resetRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if resetToken == nil {
		return nil, errors.New("invalid or expired token")
	}
	return resetToken, nil
}

// ResetPassword resets user password using a valid token
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	resetToken, err := s.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Mark token as used
	if err := s.resetRepo.MarkAsUsed(resetToken.ID); err != nil {
		return err
	}

	// Invalidate all user sessions (force logout from all devices)
	sessionRepo := repositories.NewSessionRepository()
	_ = sessionRepo.DeleteAllByUserID(user.ID)

	return nil
}
