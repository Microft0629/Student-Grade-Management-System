// Package api grade_api.go 成绩管理 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// GradeAPI 成绩管理相关的 API 结构体
type GradeAPI struct{}

// NewGradeAPI 创建并返回一个新的 GradeAPI 实例
func NewGradeAPI() *GradeAPI {
	return &GradeAPI{}
}

// CreateGrade 对外提供录入成绩的 API 接口，调用业务逻辑层处理（含绩点计算）
func (g *GradeAPI) CreateGrade(grade model.Grade) error {
	return service.CreateGrade(&grade)
}

// UpdateGrade 对外提供修改成绩的 API 接口
func (g *GradeAPI) UpdateGrade(id uint, newScore float64) error {
	return service.UpdateGrade(id, newScore)
}

// GetAllGrades 对外提供获取所有成绩及关联信息的 API 接口
func (g *GradeAPI) GetAllGrades() ([]model.Grade, error) {
	return service.GetAllGrades()
}

// DeleteGrade 对外提供根据 ID 删除成绩记录的 API 接口
func (g *GradeAPI) DeleteGrade(id uint) error {
	return service.DeleteGrade(id)
}

// SearchGrades 对外提供查询成绩的 API 接口
func (g *GradeAPI) SearchGrades(
	studentKeyword string,
	courseKeyword string,
	term string,
) ([]model.Grade, error) {
	return service.SearchGrades(
		studentKeyword,
		courseKeyword,
		term,
	)
}

// BatchImportGrades 对外提供批量导入成绩的 API 接口，返回成功数和错误列表
func (g *GradeAPI) BatchImportGrades(grades []model.Grade) (int, []string, error) {
	return service.BatchImportGrades(grades)
}

// BatchAdjustScores 按课程+分数段批量加减分调整
func (g *GradeAPI) BatchAdjustScores(courseID uint, minScore float64, maxScore float64, delta float64) (*service.BatchAdjustResult, error) {
	return service.BatchAdjustScores(courseID, minScore, maxScore, delta)
}

// AggregateGrades 跨课程/跨学期成绩汇总
func (g *GradeAPI) AggregateGrades(term string, courseKeyword string) ([]service.AggregatedGrade, error) {
	return service.AggregateGrades(term, courseKeyword)
}
