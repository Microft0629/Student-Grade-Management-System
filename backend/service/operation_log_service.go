// Package service operation_log_service.go 操作日志服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"fmt"
	"time"
)

// LogOperation 记录一条成绩修改操作日志
func LogOperation(entry model.OperationLog) {
	entry.Time = time.Now()
	err := config.DB.Create(&entry).Error
	if err != nil {
		println("记录操作日志失败:", err.Error())
	}
}

// GetAllOperationLogs 获取全部操作日志（按时间倒序）
func GetAllOperationLogs() ([]model.OperationLog, error) {
	var logs []model.OperationLog
	err := config.DB.Order("time DESC").Find(&logs).Error
	return logs, err
}

// GetOperationLogsByTerm 按学期查询操作日志
func GetOperationLogsByTerm(term string) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	err := config.DB.
		Where("term = ?", term).
		Order("time DESC").
		Find(&logs).Error
	return logs, err
}

// GetOperationLogsByStudent 按学号查询操作日志
func GetOperationLogsByStudent(studentID string) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	err := config.DB.
		Where("student_id = ?", studentID).
		Order("time DESC").
		Find(&logs).Error
	return logs, err
}

// SearchOperationLogs 按多条件追溯操作日志
func SearchOperationLogs(action string, studentID string, course string, term string, startTime string, endTime string) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	query := config.DB.Order("time DESC")

	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if studentID != "" {
		query = query.Where("student_id LIKE ?", "%"+studentID+"%")
	}
	if course != "" {
		query = query.Where("course LIKE ?", "%"+course+"%")
	}
	if term != "" {
		query = query.Where("term = ?", term)
	}
	if startTime != "" {
		t, err := time.Parse("2006-01-02", startTime)
		if err == nil {
			query = query.Where("time >= ?", t)
		}
	}
	if endTime != "" {
		t, err := time.Parse("2006-01-02", endTime)
		if err == nil {
			query = query.Where("time <= ?", t.Add(24*time.Hour))
		}
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	return logs, nil
}
