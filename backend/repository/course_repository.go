// Package repository course_repository.go 课程数据访问层
package repository

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
)

// CreateCourse 在数据库中创建一条新的课程记录
func CreateCourse(course *model.Course) error {
	return config.DB.Create(course).Error
}

// GetAllCourses 从数据库中查询并返回所有的课程记录
func GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	err := config.DB.Find(&courses).Error
	return courses, err
}

// DeleteCourse 根据指定的 ID 从数据库中删除对应的课程记录
func DeleteCourse(id uint) error {
	return config.DB.Delete(
		&model.Course{},
		id,
	).Error
}

// GetCourseByCode 根据课程代码查询课程
func GetCourseByCode(courseCode string) (*model.Course, error) {
	var course model.Course

	err := config.DB.
		Where("course_code = ?", courseCode).
		First(&course).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// SearchCoursesByName 按课程名称模糊查询
func SearchCoursesByName(keyword string) ([]model.Course, error) {
	var courses []model.Course

	err := config.DB.
		Where(
			"course_name LIKE ?",
			"%"+keyword+"%",
		).
		Order("course_name ASC").
		Find(&courses).
		Error

	return courses, err
}
