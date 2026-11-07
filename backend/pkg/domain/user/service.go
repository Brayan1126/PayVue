package user

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrDatabaseError      = errors.New("database error")
	ErrHashingPassword    = errors.New("error hashing password")
)

type Service interface {
	Register(ctx context.Context, request RegisterRequest) (*User, error)
	Login(ctx context.Context, request LoginRequest) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type service struct {
	*Container
}

func New(container *Container) Service {
	return &service{
		Container: container,
	}
}

func (s *service) Register(ctx context.Context, request RegisterRequest) (*User, error) {
	// Verificar si el email ya existe
	existingUser, err := s.Repository.GetUserByEmail(ctx, request.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrHashingPassword
	}

	user := &User{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *service) Login(ctx context.Context, request LoginRequest) (*User, error) {
	// Buscar usuario por email
	user, err := s.Repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.Repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
