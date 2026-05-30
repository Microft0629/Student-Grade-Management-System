// Package model error_log.go 校验错误日志数据模型
package model

import "time"

// ErrorLog 单条校验错误日志
type ErrorLog struct {
	Time    time.Time // 错误发生时间
	Student string    // 学生标识（学号或姓名）
	Course  string    // 课程标识
	Score   float64   // 触发错误的分数值
	Reason  string    // 错误原因描述
}
