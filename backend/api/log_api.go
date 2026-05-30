// Package api log_api.go 操作日志 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// LogAPI 操作日志 API 结构体
type LogAPI struct{}

// NewLogAPI 创建并返回一个新的 LogAPI 实例
func NewLogAPI() *LogAPI {
	return &LogAPI{}
}

// GetAllOperationLogs 获取全部操作日志
func (l *LogAPI) GetAllOperationLogs() ([]model.OperationLog, error) {
	return service.GetAllOperationLogs()
}

// GetOperationLogsByTerm 按学期查询操作日志
func (l *LogAPI) GetOperationLogsByTerm(term string) ([]model.OperationLog, error) {
	return service.GetOperationLogsByTerm(term)
}

// GetOperationLogsByStudent 按学号查询操作日志
func (l *LogAPI) GetOperationLogsByStudent(studentID string) ([]model.OperationLog, error) {
	return service.GetOperationLogsByStudent(studentID)
}

// SearchOperationLogs 多条件追溯操作日志
func (l *LogAPI) SearchOperationLogs(action string, studentID string, course string, term string, startTime string, endTime string) ([]model.OperationLog, error) {
	return service.SearchOperationLogs(action, studentID, course, term, startTime, endTime)
}

// ReadErrorLogs 读取校验错误日志
func (l *LogAPI) ReadErrorLogs() ([]model.ErrorLog, error) {
	return service.ReadErrorLogs()
}
