// Package api auth_api.go 用户登录API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

type AuthAPI struct{}

func NewAuthAPI() *AuthAPI {
	return &AuthAPI{}
}

func (a *AuthAPI) Login(
	username string,
	password string,
) (*model.User, error) {

	req := model.LoginRequest{
		Username: username,
		Password: password,
	}

	return service.Login(req)
}
