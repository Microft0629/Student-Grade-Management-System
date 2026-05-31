// Package api log_api.go 操作日志 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
	"errors"
)

// LogAPI 操作日志 API 结构体
type LogAPI struct{}

// NewLogAPI 创建并返回一个新的 LogAPI 实例
func NewLogAPI() *LogAPI {
	return &LogAPI{}
}

func (l *LogAPI) checkAdmin() error {
	if !service.IsAdmin() {
		return errors.New("仅管理员可查看操作日志")
	}
	return nil
}

// GetAllOperationLogs 获取全部操作日志（仅管理员）
func (l *LogAPI) GetAllOperationLogs() ([]model.OperationLog, error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}
	return service.GetAllOperationLogs()
}

// GetOperationLogsByTerm 按学期查询操作日志（仅管理员）
func (l *LogAPI) GetOperationLogsByTerm(term string) ([]model.OperationLog, error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}
	return service.GetOperationLogsByTerm(term)
}

// GetOperationLogsByStudent 按学号查询操作日志（仅管理员）
func (l *LogAPI) GetOperationLogsByStudent(studentID string) ([]model.OperationLog, error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}
	return service.GetOperationLogsByStudent(studentID)
}

// SearchOperationLogs 多条件追溯操作日志（仅管理员）
func (l *LogAPI) SearchOperationLogs(action string, studentID string, course string, term string, startTime string, endTime string) ([]model.OperationLog, error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}
	return service.SearchOperationLogs(action, studentID, course, term, startTime, endTime)
}

// ReadErrorLogs 读取校验错误日志（仅管理员）
func (l *LogAPI) ReadErrorLogs() ([]model.ErrorLog, error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}
	return service.ReadErrorLogs()
}
