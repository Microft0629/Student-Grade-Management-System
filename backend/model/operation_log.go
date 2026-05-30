// Package model operation_log.go 操作日志数据模型
package model

import "time"

// OperationLog 成绩修改操作日志
type OperationLog struct {
	ID        uint      `Gorm:"primaryKey"` // 主键
	Time      time.Time // 操作时间
	Action    string    // 操作类型：新增/修改/删除/批量调整/导入
	Student   string    // 学生姓名
	StudentID string    // 学号
	Course    string    // 课程名称
	Term      string    // 学期
	OldScore  float64   // 修改前分数（新增为0）
	NewScore  float64   // 修改后分数（删除为0）
	Detail    string    // 操作详情
}
