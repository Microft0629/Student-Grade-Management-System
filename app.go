package main

import (
	"Student-Grade-Management-System/backend/api"
)

// App 结构
type App struct {
	AuthAPI *api.AuthAPI
}

// NewApp 创建 APP
func NewApp() *App {
	return &App{
		AuthAPI: &api.AuthAPI{},
	}
}
