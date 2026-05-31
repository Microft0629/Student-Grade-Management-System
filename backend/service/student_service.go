// Package service student_service.go 学生业务逻辑层
package service

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"errors"
)

// CreateStudent 创建新学生记录（仅管理员），成功后自动同步更新 CSV 备份文件
func CreateStudent(student *model.Student) error {
	if !IsAdmin() {
		return errors.New("仅管理员可新增学生")
	}
	if student.StudentID == "" {
		return errors.New("学号不能为空")
	}
	if student.Name == "" {
		return errors.New("姓名不能为空")
	}
	if student.Gender == "" {
		return errors.New("性别不能为空")
	}
	if student.ClassName == "" {
		return errors.New("班级不能为空")
	}
	if student.Major == "" {
		return errors.New("专业不能为空")
	}

	// 将学生数据持久化至数据库
	err := repository.CreateStudent(student)
	if err != nil {
		return err
	}

	// 数据库写入成功后，同步刷新 CSV 文件以保持数据一致性
	return SyncStudentsToCSV()
}

// GetAllStudents 获取所有学生的业务数据
func GetAllStudents() ([]model.Student, error) {
	return repository.GetAllStudents()
}

// DeleteStudent 根据 ID 删除学生记录（仅管理员），成功后自动同步更新 CSV 备份文件
func DeleteStudent(id uint) error {
	if !IsAdmin() {
		return errors.New("仅管理员可删除学生")
	}
	// 从数据库中删除指定 ID 的学生记录
	err := repository.DeleteStudent(id)
	if err != nil {
		return err
	}

	// 数据库删除成功后，同步刷新 CSV 文件以保持数据一致性
	return SyncStudentsToCSV()
}

// SearchStudents 根据关键词搜索学生信息，将请求委托给数据访问层处理
func SearchStudents(keyword string) ([]model.Student, error) {
	return repository.SearchStudents(keyword)
}

// GetStudentsByPage 分页查询学生列表
func GetStudentsByPage(page int, pageSize int) (
	model.StudentPageResult,
	error,
) {
	return repository.GetStudentsByPage(page, pageSize)
}
