// Package model grade.go 成绩数据类型
package model

type Grade struct {
	ID          uint    `Gorm:"primaryKey"`                   // 主键
	StudentID   uint    `Gorm:"uniqueIndex:idx_grade_unique"` // 关联学生 ID
	CourseID    uint    `Gorm:"uniqueIndex:idx_grade_unique"` // 关联课程 ID
	Score       float64 // 实际成绩
	GradePoint  float64 // 绩点
	CreatorName string  // 录入人用户名
	Student     Student `Gorm:"foreignKey:StudentID"` // 学生数据结构体
	Course      Course  `Gorm:"foreignKey:CourseID"`  // 课程数据结构体
}
