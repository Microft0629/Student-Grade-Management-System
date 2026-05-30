// Package repository grade_repository.go 成绩数据访问层
package repository

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"strings"
)

// CreateGrade 在数据库中创建一条新的成绩记录
func CreateGrade(grade *model.Grade) error {
	return config.DB.Create(grade).Error
}

// GetAllGrades 从数据库中查询所有成绩记录，并预加载关联的学生与课程信息
func GetAllGrades() ([]model.Grade, error) {
	var grades []model.Grade
	err := config.DB.
		Preload("Student").
		Preload("Course").
		Find(&grades).Error

	return grades, err
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

	// 预加载关联对象
	query := config.DB.
		Preload("Student").
		Preload("Course")
	err := query.Find(&grades).Error
	if err != nil {
		return nil, err
	}

	var result []model.Grade

	// 学生/课程关键词做模糊匹配
	for _, grade := range grades {
		// 学生关键词同时匹配姓名与学号，为空则跳过该条件
		matchStudent := studentKeyword == "" ||
			strings.Contains(
				grade.Student.Name,
				studentKeyword,
			) ||
			strings.Contains(
				grade.Student.StudentID,
				studentKeyword,
			)
		// 课程关键词仅匹配课程名称，为空则跳过该条件
		matchCourse := courseKeyword == "" ||
			strings.Contains(
				grade.Course.CourseName,
				courseKeyword,
			)
		// 学期关键词仅匹配学期名，为空则跳过该条件
		matchTerm := term == "" ||
			grade.Course.Term == term

		// 所有非空条件均满足时才纳入结果
		if matchStudent && matchCourse && matchTerm {
			result = append(
				result,
				grade,
			)
		}
	}

	return result, nil
}
