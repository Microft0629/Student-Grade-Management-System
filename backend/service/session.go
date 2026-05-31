// Package service session.go 当前登录用户会话
package service

import "Student-Grade-Management-System/backend/model"

// currentSession 存储当前登录用户（单用户桌面应用，全局一份即可）
var currentSession struct {
	User     *model.User
	LoggedIn bool
}

// SetCurrentUser 设置当前登录用户
func SetCurrentUser(user *model.User) {
	currentSession.User = user
	currentSession.LoggedIn = true
}

// GetCurrentUser 获取当前登录用户，未登录返回 nil
func GetCurrentUser() *model.User {
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

// CurrentOperator 返回当前操作人标识，用于日志记录
func CurrentOperator() string {
	u := GetCurrentUser()
	if u == nil {
		return "未知"
	}
	return u.Username
}
