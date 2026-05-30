// Package api gpa_api.go 绩点规则管理 API
package api

import "Student-Grade-Management-System/backend/service"

// GpaAPI 绩点规则管理 API 结构体
type GpaAPI struct{}

// NewGpaAPI 创建并返回一个新的 GpaAPI 实例
func NewGpaAPI() *GpaAPI {
	return &GpaAPI{}
}

// GetGpaRules 获取当前绩点换算公式（优先从 TXT 加载，否则返回默认）
func (g *GpaAPI) GetGpaRules() (string, error) {
	rules, err := service.LoadGpaRules()
	if err != nil {
		return "", err
	}
	if rules == "" {
		return service.GetDefaultGpaRules(), nil
	}
	return rules, nil
}

// SaveGpaRules 保存自定义绩点换算公式
func (g *GpaAPI) SaveGpaRules(formula string) error {
	return service.SaveGpaRules(formula)
}

// ResetGpaRules 重置为默认绩点换算公式
func (g *GpaAPI) ResetGpaRules() error {
	return service.SaveGpaRules(service.GetDefaultGpaRules())
}

// RecalculateAllGPA 批量重新计算所有学生绩点
func (g *GpaAPI) RecalculateAllGPA() (int, error) {
	return service.RecalculateAllGPA()
}
