// Package model course_statistics.go 课程统计数据模型
package model

// CourseStatistics 单门课程的成绩统计数据
type CourseStatistics struct {
	CourseName   string  // 课程名称
	Term         string  // 学期
	StudentCount int     // 参考学生人数
	AverageScore float64 // 平均分
	PassRate     float64 // 及格率（百分制）
	HighestScore float64 // 最高分
	LowestScore  float64 // 最低分
	// 分数段分布：0-59, 60-69, 70-79, 80-89, 90-100
	Distribution map[string]int // 分数段 → 人数
}

// ScoreDistribution 分数段分布
type ScoreDistribution struct {
	Range0_59  int // 0-59 人数
	Range60_69 int // 60-69 人数
	Range70_79 int // 70-79 人数
	Range80_89 int // 80-89 人数
	Range90_100 int // 90-100 人数
}

// StudentRanking 学生排名数据
type StudentRanking struct {
	StudentName   string  // 学生姓名
	StudentID     string  // 学号
	TotalScore    float64 // 多课程总分
	AverageScore  float64 // 平均分
	GPA           float64 // 平均绩点
	CourseCount   int     // 课程数
	Rank          int     // 排名
}
