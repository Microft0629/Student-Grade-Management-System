// Package model course.go 课程数据模型
package model

type Course struct {
	ID         uint    `Gorm:"primaryKey"` // 主键
	CourseCode string  `Gorm:"unique"`     // 课程代码
	CourseName string  // 课程名称
	Term       string  // 学期
	Credit     float64 // 学分
	Teacher    string  // 任课教师
}
