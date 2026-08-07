// Package model student.go 学生数据模型
package model

import "time"

type Student struct {
	ID        uint      `gorm:"primaryKey"` // 主键
	StudentID string    `gorm:"unique"`     // 学号（唯一）
	Name      string    // 姓名
	Gender    string    // 性别
	ClassName string    // 班级
	Major     string    // 专业
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
