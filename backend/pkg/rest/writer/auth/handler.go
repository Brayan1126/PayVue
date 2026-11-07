package auth

import (
	"github.com/payvue/payvue-backend/pkg/domain/user"
	"github.com/payvue/payvue-backend/pkg/rest"
)

type handler struct {
	userService user.Service
}

func NewHandler(userService user.Service) rest.Handler {
	return &handler{
		userService: userService,
	}
}
