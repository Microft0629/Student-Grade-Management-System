// Package service csv_sync_service.go CSV 与 SQLite 同步服务
package service

import (
	"Student-Grade-Management-System/backend/repository"
	csvRepo "Student-Grade-Management-System/backend/repository/csv"
	"Student-Grade-Management-System/backend/utils"
	"os"
	"path/filepath"
	"strings"
)

// SyncStudentsToCSV 从数据库读取全部学生数据并同步导出至 CSV 文件，用于数据备份或导出
func SyncStudentsToCSV() error {
	// 获取全部学生记录
	students, err :=
		repository.GetAllStudents()
	if err != nil {
		return err
	}

	// 将获得的学生列表写入CSV文件
	return csvRepo.SaveStudents(
		students,
	)
}

// SyncCoursesToCSV 从数据库读取全部课程数据并同步导出至 CSV 文件，用于数据备份或导出
func SyncCoursesToCSV() error {
	// 获取所有课程记录
	courses, err := repository.GetAllCourses()
	if err != nil {
		return err
	}

	// 将获得的课程列表写入CSV文件
	return csvRepo.SaveCourses(
		courses,
	)
}

// SyncGradesToCSV 从数据库读取全部成绩并按学期+课程分组同步写入 CSV 文件，
// 同时清理数据库中已不存在对应成绩记录的旧 CSV 文件，保证 CSV 与数据库一致。
func SyncGradesToCSV() error {
	grades, err := repository.GetAllGrades()
	if err != nil {
		return err
	}

	// 按 "学期|课程代码" 分组，避免重复写入同一文件
	group := make(map[string][]csvRepo.GradeCSV)
	for _, grade := range grades {
		if grade.Course.Term == "" {
			continue
		}

		courseKey := grade.Course.CourseCode
		if courseKey == "" {
			courseKey = grade.Course.CourseName
		}
		key := grade.Course.Term + "|" + courseKey

		group[key] = append(group[key], csvRepo.GradeCSV{
			StudentID: grade.Student.StudentID,
			Name:      grade.Student.Name,
			Score:     grade.Score,
		})
	}

	// 遍历分组后的数据，逐组写入对应 CSV 文件
	for key, gradeList := range group {
		// 从组合键中解析出学期和课程代码
		parts := strings.SplitN(key, "|", 2)
		term := parts[0]
		courseCode := parts[1]

		err = csvRepo.SaveCourseGrades(
			term,
			courseCode,
			gradeList,
		)
		if err != nil {
			return err
		}
	}

	// 清理不再对应任何数据库成绩记录的旧 CSV 文件
	return removeStaleGradeFiles(group)
}

// removeStaleGradeFiles 删除 data/grades 下不在当前数据库成绩集合中的旧 CSV 文件
func removeStaleGradeFiles(group map[string][]csvRepo.GradeCSV) error {
	root := filepath.Join(utils.DataDir(), "grades")
	termDirs, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, termDir := range termDirs {
		if !termDir.IsDir() || termDir.Type()&os.ModeSymlink != 0 {
			continue
		}
		dir := filepath.Join(root, termDir.Name())
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() || file.Type()&os.ModeSymlink != 0 || filepath.Ext(file.Name()) != ".csv" {
				continue
			}
			courseCode := strings.TrimSuffix(file.Name(), ".csv")
			key := termDir.Name() + "|" + courseCode
			if _, exists := group[key]; exists {
				continue
			}
			if err := os.Remove(filepath.Join(dir, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
