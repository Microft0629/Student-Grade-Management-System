// Package api user_api.go 用户管理 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// UserAPI 用户管理 API 结构体
type UserAPI struct{}

// NewUserAPI 创建并返回一个新的 UserAPI 实例
func NewUserAPI() *UserAPI {
	return &UserAPI{}
}

// CreateTeacher 管理员创建老师账号
func (u *UserAPI) CreateTeacher(username string, password string) error {
	return service.CreateTeacher(username, password)
}

// DeleteUser 管理员删除老师账号
func (u *UserAPI) DeleteUser(username string) error {
	return service.DeleteUser(username)
}

// GetAllTeachers 获取所有老师账号列表
func (u *UserAPI) GetAllTeachers() ([]model.User, error) {
	return service.GetAllTeachers()
}

// ChangePassword 修改当前用户密码
func (u *UserAPI) ChangePassword(oldPwd string, newPwd string) error {
	return service.ChangePassword(oldPwd, newPwd)
}
