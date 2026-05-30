// Package repository student_repository.go 学生数据访问层
package repository

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
)

// CreateStudent 在数据库中创建一条新的学生记录
func CreateStudent(student *model.Student) error {
	return config.DB.Create(student).Error
}

// GetAllStudents 从数据库中查询并返回所有的学生记录
func GetAllStudents() ([]model.Student, error) {
	var students []model.Student
	err := config.DB.Find(&students).Error
	return students, err
}

// DeleteStudent 根据指定的 ID 从数据库中删除对应的学生记录
func DeleteStudent(id uint) error {
	return config.DB.Delete(
		&model.Student{},
		id,
	).Error
}

// SearchStudents 根据姓名搜索学生
func SearchStudents(keyword string) ([]model.Student, error) {
	var students []model.Student

	err := config.DB.
		Where(
			"name LIKE ?",
			"%"+keyword+"%",
		).
		Find(&students).Error

	return students, err
}

// GetStudentByStudentID 根据学号查询学生记录
func GetStudentByStudentID(studentID string) (*model.Student, error) {
	var student model.Student

	err := config.DB.
		Where("student_id = ?", studentID).
		First(&student).
		Error

	if err != nil {
		return nil, err
	}

	return &student, nil
}

// GetStudentsByPage 分页查询学生列表，返回包含数据列表、总数及分页参数的结果
func GetStudentsByPage(page int, pageSize int) (
	model.StudentPageResult,
	error,
) {
	var result model.StudentPageResult

	// 防御性校验：页码最小为1
	if page < 1 {
		page = 1
	}

	// 防御性校验：每页条数最小为10
	if pageSize < 1 {
		pageSize = 10
	}

	// 防止恶意超大查询：每页条数最大限制为100
	if pageSize > 100 {
		pageSize = 100
	}

	var students []model.Student
	var total int64

	// 查询满足条件的学生总记录数，用于前端分页组件渲染
	err := config.DB.
		Model(&model.Student{}).
		Count(&total).Error

	if err != nil {
		return result, err
	}

	// 根据页码和每页大小计算数据库查询偏移量
	offset := (page - 1) * pageSize

	// 按ID倒序执行分页查询，获取当前页的学生数据
	err = config.DB.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&students).Error

	if err != nil {
		return result, err
	}

	// 组装分页查询结果并返回
	result = model.StudentPageResult{
		List:     students,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return result, nil
}
