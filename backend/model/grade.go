// Package model grade.go 成绩数据类型
package model

type Grade struct {
	ID          uint    `gorm:"primaryKey"`                   // 主键
	StudentID   uint    `gorm:"uniqueIndex:idx_grade_unique"` // 关联学生表主键
	CourseID    uint    `gorm:"uniqueIndex:idx_grade_unique"` // 关联课程表主键
	Score       float64 // 实际成绩
	GradePoint  float64 // 绩点
	CreatorName string  // 录入人用户名
	Student     Student `gorm:"foreignKey:StudentID;references:ID"` // 学生数据
	Course      Course  `gorm:"foreignKey:CourseID;references:ID"`  // 课程数据
}
