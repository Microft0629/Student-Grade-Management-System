// Package api course_api.go 课程管理 API
package api

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/service"
)

// CourseAPI 课程管理相关的 API 结构体
type CourseAPI struct{}

// NewCourseAPI 创建并返回一个新的 CourseAPI 实例
func NewCourseAPI() *CourseAPI {
	return &CourseAPI{}
}

// CreateCourse 对外提供创建课程的 API 接口，调用业务逻辑层处理
func (c *CourseAPI) CreateCourse(course model.Course) error {
	return service.CreateCourse(&course)
}

// GetAllCourses 对外提供获取所有课程信息的 API 接口
func (c *CourseAPI) GetAllCourses() ([]model.Course, error) {
	return service.GetAllCourses()
}

// DeleteCourse 对外提供根据 ID 删除课程的 API 接口
func (c *CourseAPI) DeleteCourse(id uint) error {
	return service.DeleteCourse(id)
}

// SearchCourses 对外提供搜索课程的 API 接口
func (c *CourseAPI) SearchCourses(keyword string) ([]model.Course, error) {
	return service.SearchCourses(keyword)
}
