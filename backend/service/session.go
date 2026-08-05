// Package service session.go 当前登录用户会话
package service

import (
	"errors"

	"Student-Grade-Management-System/backend/model"
)

// currentSession 存储当前登录用户（单用户桌面应用，全局一份即可）
var currentSession struct {
	User     *model.UserInfo
	LoggedIn bool
}

// SetCurrentUser 设置当前登录用户
func SetCurrentUser(user *model.UserInfo) {
	currentSession.User = user
	currentSession.LoggedIn = true
}

// ClearCurrentUser 清除当前登录用户（退出登录）
func ClearCurrentUser() {
	currentSession.User = nil
	currentSession.LoggedIn = false
}

// IsLoggedIn 是否已登录
func IsLoggedIn() bool {
	return currentSession.LoggedIn && currentSession.User != nil
}

// GetCurrentUser 获取当前登录用户，未登录返回 nil
func GetCurrentUser() *model.UserInfo {
	if !currentSession.LoggedIn {
		return nil
	}
	return currentSession.User
}

// IsAdmin 当前用户是否为管理员
func IsAdmin() bool {
	u := GetCurrentUser()
	return u != nil && u.Role == "admin"
}

// RequireLogin 校验当前是否已登录
func RequireLogin() error {
	if !IsLoggedIn() {
		return errors.New("未登录")
	}
	return nil
}

// RequireAdmin 校验当前是否为管理员（未登录时优先返回未登录错误）
func RequireAdmin() error {
	if err := RequireLogin(); err != nil {
		return err
	}
	if !IsAdmin() {
		return errors.New("仅管理员可执行此操作")
	}
	return nil
}

// CurrentOperator 返回当前操作人标识，用于日志记录
func CurrentOperator() string {
	u := GetCurrentUser()
	if u == nil {
		return "未知"
	}
	return u.Username
}
