// Package model login.go 登录请求结构
package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
