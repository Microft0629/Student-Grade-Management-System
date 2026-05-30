// Package model statistics.go 学生成绩统计数据模型
package model

type StudentStatistics struct {
	StudentName  string  // 学生姓名
	AverageScore float64 // 平均分
	GPA          float64 // GPA
	TotalCredits float64 // 总学分
	CourseCount  int     // 课程总数
}
