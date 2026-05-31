// Package repository grade_repository.go 成绩数据访问层
package repository

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"strings"
)

// LoadAssociations 手动加载成绩记录关联的学生和课程（替代 Preload，避免 GORM 关联推断问题）
func LoadAssociations(grades []model.Grade) {
	for i := range grades {
		var student model.Student
		if config.DB.First(&student, grades[i].StudentID).Error == nil {
			grades[i].Student = student
		}
		var course model.Course
		if config.DB.First(&course, grades[i].CourseID).Error == nil {
			grades[i].Course = course
		}
	}
}

// CreateGrade 在数据库中创建一条新的成绩记录
func CreateGrade(grade *model.Grade) error {
	return config.DB.Create(grade).Error
}

// GetAllGrades 从数据库中查询所有成绩记录，并加载关联的学生与课程信息
func GetAllGrades() ([]model.Grade, error) {
	var grades []model.Grade
	err := config.DB.Find(&grades).Error
	if err != nil {
		return nil, err
	}
	LoadAssociations(grades)
	return grades, nil
}

// DeleteGrade 根据指定的 ID 从数据库中删除对应的成绩记录
func DeleteGrade(id uint) error {
	return config.DB.Delete(
		&model.Grade{},
		id,
	).Error
}

// GradeExists 判断成绩是否存在
func GradeExists(studentID uint, courseID uint) (bool, error) {
	var count int64

	err := config.DB.
		Model(&model.Grade{}).
		Where(
			"student_id = ? AND course_id = ?",
			studentID,
			courseID,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SearchGrades 多条件查询成绩
func SearchGrades(
	studentKeyword string,
	courseKeyword string,
	term string,
) ([]model.Grade, error) {
	var grades []model.Grade

	err := config.DB.Find(&grades).Error
	if err != nil {
		return nil, err
	}
	LoadAssociations(grades)

	var result []model.Grade

	for _, grade := range grades {
		matchStudent := studentKeyword == "" ||
			strings.Contains(grade.Student.Name, studentKeyword) ||
			strings.Contains(grade.Student.StudentID, studentKeyword)
		matchCourse := courseKeyword == "" ||
			strings.Contains(grade.Course.CourseName, courseKeyword)
		matchTerm := term == "" ||
			grade.Course.Term == term

		if matchStudent && matchCourse && matchTerm {
			result = append(result, grade)
		}
	}

	return result, nil
}
