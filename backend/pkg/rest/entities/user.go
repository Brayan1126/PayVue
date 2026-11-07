package entities

import "github.com/payvue/payvue-backend/pkg/domain/user"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r RegisterRequest) ToDomain() user.RegisterRequest {
	return user.RegisterRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r LoginRequest) ToDomain() user.LoginRequest {
	return user.LoginRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}
