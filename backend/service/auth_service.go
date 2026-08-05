// Package service auth_service.go 用户登录服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Login 校验用户名密码，返回不含密码哈希的用户信息
func Login(req model.LoginRequest) (*model.UserInfo, error) {
	var user model.User
	// 查数据库，寻找用户
	result := config.DB.
		Where("username = ?", req.Username).
		First(&user)

	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}
	// 验证密码
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("密码错误")
	}
	// 返回不含密码哈希的用户信息
	return &model.UserInfo{
		ID:                 user.ID,
		Username:           user.Username,
		Role:               user.Role,
		MustChangePassword: user.MustChangePassword,
	}, nil
}

// Logout 退出登录，清空当前会话
func Logout() {
	ClearCurrentUser()
}
