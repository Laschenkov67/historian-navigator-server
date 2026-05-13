package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/historian/backend/internal/models"
	"github.com/historian/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials for login failure.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService — authentication service.
type AuthService struct {
	users  *repository.UserRepo
	secret string
}

// NewAuthService constructor.
func NewAuthService(u *repository.UserRepo, secret string) *AuthService {
	return &AuthService{users: u, secret: secret}
}

// Register new user.
func (s *AuthService) Register(ctx context.Context, email, pwd, name string, role models.Role) (*models.User, string, error) {
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return nil, "", errors.New("user already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	u := &models.User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  string(hash),
		FullName:  name,
		Role:      role,
		CreatedAt: time.Now(),
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, "", err
	}
	token, err := s.issueToken(u)
	return u, token, err
}

// Login checks credentials and issues token.
func (s *AuthService) Login(ctx context.Context, email, pwd string) (*models.User, string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := s.issueToken(u)
	return u, token, err
}

// ParseToken validates token string.
func (s *AuthService) ParseToken(token string) (*models.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) { return []byte(s.secret), nil })
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return &models.User{
		ID:    claims["sub"].(string),
		Email: claims["email"].(string),
		Role:  models.Role(claims["role"].(string)),
	}, nil
}

func (s *AuthService) issueToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  string(u.Role),
		"exp":   time.Now().Add(24 * 7 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secret))
}
