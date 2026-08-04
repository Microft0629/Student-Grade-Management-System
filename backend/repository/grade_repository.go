// Package repository grade_repository.go 成绩数据访问层
package repository

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
)

// LoadAssociations 批量加载成绩记录关联的学生和课程（两次 IN 查询，避免逐条查询）
func LoadAssociations(grades []model.Grade) {
	if len(grades) == 0 {
		return
	}

	studentIDSet := make(map[uint]struct{})
	courseIDSet := make(map[uint]struct{})
	for _, grade := range grades {
		studentIDSet[grade.StudentID] = struct{}{}
		courseIDSet[grade.CourseID] = struct{}{}
	}

	studentIDs := make([]uint, 0, len(studentIDSet))
	for id := range studentIDSet {
		studentIDs = append(studentIDs, id)
	}
	courseIDs := make([]uint, 0, len(courseIDSet))
	for id := range courseIDSet {
		courseIDs = append(courseIDs, id)
	}

	studentMap := make(map[uint]model.Student, len(studentIDs))
	if len(studentIDs) > 0 {
		var students []model.Student
		if err := config.DB.Where("id IN ?", studentIDs).Find(&students).Error; err == nil {
			for _, student := range students {
				studentMap[student.ID] = student
			}
		}
	}

	courseMap := make(map[uint]model.Course, len(courseIDs))
	if len(courseIDs) > 0 {
		var courses []model.Course
		if err := config.DB.Where("id IN ?", courseIDs).Find(&courses).Error; err == nil {
			for _, course := range courses {
				courseMap[course.ID] = course
			}
		}
	}

	for i := range grades {
		if student, ok := studentMap[grades[i].StudentID]; ok {
			grades[i].Student = student
		}
		if course, ok := courseMap[grades[i].CourseID]; ok {
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

	query := config.DB.Model(&model.Grade{}).
		Select("grades.*").
		Joins("JOIN students ON students.id = grades.student_id").
		Joins("JOIN courses ON courses.id = grades.course_id")

	if studentKeyword != "" {
		like := "%" + studentKeyword + "%"
		query = query.Where("(students.name LIKE ? OR students.student_id LIKE ?)", like, like)
	}
	if courseKeyword != "" {
		query = query.Where("courses.course_name LIKE ?", "%"+courseKeyword+"%")
	}
	if term != "" {
		query = query.Where("courses.term = ?", term)
	}

	if err := query.Find(&grades).Error; err != nil {
		return nil, err
	}
	LoadAssociations(grades)
	return grades, nil
}
