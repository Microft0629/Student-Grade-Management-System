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
) (*model.UserInfo, error) {

	req := model.LoginRequest{
		Username: username,
		Password: password,
	}

	user, err := service.Login(req)
	if err != nil {
		return nil, err
	}
	service.SetCurrentUser(user)
	return user, nil
}

// Logout 退出登录，清空后端会话
func (a *AuthAPI) Logout() error {
	service.Logout()
	return nil
}

// GetCurrentUserRole 获取当前登录用户的角色
func (a *AuthAPI) GetCurrentUserRole() string {
	u := service.GetCurrentUser()
	if u == nil {
		return ""
	}
	return u.Role
}
