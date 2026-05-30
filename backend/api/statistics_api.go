// Package api statistics_api.go 成绩统计 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// StatisticsAPI 学生成绩统计相关的 API 结构体
type StatisticsAPI struct{}

// NewStatisticsAPI 创建并返回一个新的 StatisticsAPI 实例
func NewStatisticsAPI() *StatisticsAPI {
	return &StatisticsAPI{}
}

// GetStudentStatistics 对外提供获取所有学生成绩统计数据的 API 接口
func (s *StatisticsAPI) GetStudentStatistics() (
	[]model.StudentStatistics,
	error,
) {
	return service.GetStudentStatistics()
}

// GetCourseStatistics 获取单课程统计（平均分、及格率、分数段分布）
func (s *StatisticsAPI) GetCourseStatistics(courseID uint) (*model.CourseStatistics, error) {
	return service.GetCourseStatistics(courseID)
}

// GetStudentRanking 获取学生排名（按平均绩点降序）
func (s *StatisticsAPI) GetStudentRanking() ([]model.StudentRanking, error) {
	return service.GetStudentRanking()
}

// GenerateStatisticsReport 生成文字版统计报表
func (s *StatisticsAPI) GenerateStatisticsReport(term string) (string, error) {
	return service.GenerateStatisticsReport(term)
}
