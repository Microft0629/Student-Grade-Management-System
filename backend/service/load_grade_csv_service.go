// Package service load_grade_csv_service.go 成绩 CSV 导入服务
package service

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	csvRepo "Student-Grade-Management-System/backend/repository/csv"
	"Student-Grade-Management-System/backend/utils"
)

// LoadGradesFromCSV 从 CSV 文件批量导入成绩到数据库（已存在则跳过）
func LoadGradesFromCSV() error {
	files, err := csvRepo.ScanGradeFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		rows, err := csvRepo.LoadCourseGrades(file.Term, file.CourseCode)
		if err != nil {
			return err
		}

		course, err := repository.GetCourseByCode(file.CourseCode)
		if err != nil {
			return err
		}

		for _, row := range rows {
			student, err := repository.GetStudentByStudentID(row.StudentID)
			if err != nil {
				continue // 学生不存在则跳过该条记录
			}

			exists, err := repository.GradeExists(student.ID, course.ID)
			if err != nil {
				return err
			}
			if exists {
				continue // 成绩已存在则跳过
			}

			grade := model.Grade{
				StudentID:   student.ID,
				CourseID:    course.ID,
				Score:       row.Score,
				GradePoint:  utils.CalculateGradePoint(row.Score),
				CreatorName: "admin",
			}

			if err := repository.CreateGrade(&grade); err != nil {
				return err
			}
		}
	}

	return nil
}
