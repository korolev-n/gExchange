package service

import (
	"context"
	"errors"

	"github.com/korolev-n/gExchange/exchanger/internal/domain"
	"github.com/korolev-n/gExchange/exchanger/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hashedPassword))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	return GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
}
