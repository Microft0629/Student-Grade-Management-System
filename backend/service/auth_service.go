// Package service auth_service.go 用户登录服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Login(req model.LoginRequest) (*model.User, error) {
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
	// 返回用户信息
	return &user, nil
}
