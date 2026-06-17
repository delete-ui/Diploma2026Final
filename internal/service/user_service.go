package service

import (
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/repository"
	"GolangBackendDiploma26/internal/transport/email"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        *repository.UserRepository
	logger      *zap.Logger
	jwtSecret   []byte
	jwtTTL      time.Duration
	emailSender email.Sender
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger, jwtSecret string, jwtTTL time.Duration, emailSender email.Sender) *UserService {
	return &UserService{
		repo:        repo,
		logger:      logger,
		jwtSecret:   []byte(jwtSecret),
		jwtTTL:      jwtTTL,
		emailSender: emailSender,
	}
}

func generateCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *UserService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("generate code: %w", err)
	}
	expires := time.Now().Add(1 * time.Hour)

	user := &models.User{
		Email:                 req.Email,
		Username:              req.Username,
		PasswordHash:          string(hashed),
		Role:                  "user",
		IsVerified:            false,
		VerificationCode:      &code,
		VerificationExpiresAt: &expires,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.emailSender.SendVerificationCode(user.Email, user.Username, code); err != nil {
		s.logger.Error("failed to send verification email", zap.Error(err))
	}

	s.logger.Info("verification code generated",
		zap.String("email", user.Email),
		zap.String("code", code),
		zap.Time("expires", expires),
	)

	return user, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, email, code string) error {
	err := s.repo.VerifyEmail(ctx, email, code)
	if err != nil {
		s.logger.Warn("email verification failed", zap.String("email", email), zap.Error(err))
		return err
	}
	s.logger.Info("email verified successfully", zap.String("email", email))
	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if !user.IsVerified {
		return nil, fmt.Errorf("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   now.Add(s.jwtTTL).Unix(),
		"iat":   now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.jwtTTL.Seconds()),
		User:        *user,
	}, nil
}

func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	if user == nil {
		s.logger.Info("password reset requested for non-existent email", zap.String("email", email))
		return nil
	}

	code, err := generateCode()
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	expires := time.Now().Add(1 * time.Hour)

	if err := s.repo.SetResetCode(ctx, email, code, expires); err != nil {
		return fmt.Errorf("set reset code: %w", err)
	}

	if err := s.emailSender.SendPasswordResetCode(email, user.Username, code); err != nil {
		s.logger.Error("failed to send password reset email", zap.Error(err))
	}

	s.logger.Info("password reset code generated",
		zap.String("email", email),
		zap.String("code", code),
	)

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	if err := s.repo.VerifyResetCode(ctx, email, code); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.repo.UpdatePassword(ctx, email, string(hashed))
}
