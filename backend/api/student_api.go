// Package api student_api.go 学生管理 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// StudentAPI 学生管理相关的 API 结构体
type StudentAPI struct{}

// NewStudentAPI 创建并返回一个新的 StudentAPI 实例
func NewStudentAPI() *StudentAPI {
	return &StudentAPI{}
}

// CreateStudent 对外提供创建学生的 API 接口，调用业务逻辑层处理
func (s *StudentAPI) CreateStudent(student model.Student) error {
	return service.CreateStudent(&student)
}

// GetAllStudents 对外提供获取所有学生信息的 API 接口
func (s *StudentAPI) GetAllStudents() ([]model.Student, error) {
	return service.GetAllStudents()
}

// DeleteStudent 对外提供根据 ID 删除学生的 API 接口
func (s *StudentAPI) DeleteStudent(id uint) error {
	return service.DeleteStudent(id)
}

// SearchStudents 根据关键词搜索学生信息的 API 接口
func (s *StudentAPI) SearchStudents(keyword string) ([]model.Student, error) {
	return service.SearchStudents(keyword)
}

// GetStudentsByPage 分页查询学生列表的 API 接口
func (s *StudentAPI) GetStudentsByPage(page int, pageSize int) (
	model.StudentPageResult,
	error,
) {
	return service.GetStudentsByPage(page, pageSize)
}
