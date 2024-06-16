package user

import (
	userService "github.com/aradwann/eenergy/service/v1/user"
)

type UserHandler struct {
	service userService.UserService
	UnimplementedUserServiceServer
}

func NewUserHandler(service userService.UserService) *UserHandler {
	return &UserHandler{service: service}
}
