// Package service csv_load_service.go CSV 数据加载服务
package service

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	csvRepo "Student-Grade-Management-System/backend/repository/csv"
)

// LoadStudentsFromCSV 从 CSV 加载学生数据，已存在的学生自动跳过
func LoadStudentsFromCSV() error {
	return ImportFromCSV[model.Student, string](
		csvRepo.LoadStudents,
		func(s model.Student) string { return s.StudentID },
		func(id string) error {
			// 仅传递错误，ImportFromCSV 内部通过 errors.Is 判断是否为 ErrRecordNotFound
			_, err := repository.GetStudentByStudentID(id)
			return err
		},
		repository.CreateStudent,
		"学生",
	)
}

// LoadCoursesFromCSV 从CSV加载课程数据，已存在的课程自动跳过
func LoadCoursesFromCSV() error {
	return ImportFromCSV[model.Course, string](
		csvRepo.LoadCourses,
		func(c model.Course) string { return c.CourseCode },
		func(code string) error {
			// 仅传递错误，ImportFromCSV 内部通过 errors.Is 判断是否为 ErrRecordNotFound
			_, err := repository.GetCourseByCode(code)
			return err
		},
		repository.CreateCourse,
		"课程",
	)
}

// LoadAllCSVData 加载全部 CSV 数据
func LoadAllCSVData() error {
	err := LoadStudentsFromCSV()
	if err != nil {
		return err
	}

	err = LoadCoursesFromCSV()
	if err != nil {
		return err
	}

	err = LoadGradesFromCSV()
	if err != nil {
		return err
	}

	return nil
}
