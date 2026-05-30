// Package service backup_service.go 数据备份与恢复服务
package service

import (
	"Student-Grade-Management-System/backend/repository"
	csvRepo "Student-Grade-Management-System/backend/repository/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const backupRoot = "backup"

// BackupByTerm 备份指定学期的全部成绩数据
func BackupByTerm(term string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(backupRoot, fmt.Sprintf("%s_%s", term, timestamp))

	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 扫描成绩文件，找到属于指定学期的文件
	files, err := csvRepo.ScanGradeFiles()
	if err != nil {
		return "", err
	}

	backupCount := 0
	for _, f := range files {
		if f.Term != term {
			continue
		}

		srcPath := filepath.Join("data", "grades", f.Term, f.CourseCode+".csv")
		dstPath := filepath.Join(backupDir, f.CourseCode+".csv")

		err = copyFile(srcPath, dstPath)
		if err != nil {
			return backupDir, fmt.Errorf("备份文件[%s]失败: %w", f.CourseCode, err)
		}
		backupCount++
	}

	if backupCount == 0 {
		return backupDir, fmt.Errorf("学期[%s]没有可备份的成绩数据", term)
	}

	return backupDir, nil
}

// BackupByCourse 备份指定课程的成绩数据
func BackupByCourse(term string, courseCode string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(backupRoot, fmt.Sprintf("%s_%s_%s", courseCode, term, timestamp))

	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建备份目录失败: %w", err)
	}

	srcPath := filepath.Join("data", "grades", term, courseCode+".csv")
	dstPath := filepath.Join(backupDir, courseCode+".csv")

	// 检查源文件是否存在
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return backupDir, fmt.Errorf("课程[%s]在学期[%s]下没有成绩数据", courseCode, term)
	}

	err = copyFile(srcPath, dstPath)
	if err != nil {
		return backupDir, fmt.Errorf("备份文件失败: %w", err)
	}

	return backupDir, nil
}

// ListBackups 列出所有备份目录
func ListBackups() ([]string, error) {
	_, err := os.Stat(backupRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(backupRoot)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// RestoreFromBackup 从备份目录恢复数据
func RestoreFromBackup(backupName string) error {
	backupDir := filepath.Join(backupRoot, backupName)

	_, err := os.Stat(backupDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("备份[%s]不存在", backupName)
	}

	// 从备份目录名解析学期信息（格式：term_timestamp 或 courseCode_term_timestamp）
	parts := strings.Split(backupName, "_")
	if len(parts) < 2 {
		return fmt.Errorf("备份名称格式不正确: %s", backupName)
	}

	// 扫描备份目录下的 CSV 文件
	csvFiles, err := filepath.Glob(filepath.Join(backupDir, "*.csv"))
	if err != nil {
		return fmt.Errorf("扫描备份文件失败: %w", err)
	}

	if len(csvFiles) == 0 {
		return fmt.Errorf("备份[%s]中没有成绩文件", backupName)
	}

	for _, csvFile := range csvFiles {
		// 从文件名提取课程代码
		courseCode := strings.TrimSuffix(filepath.Base(csvFile), ".csv")

		// 从备份目录名提取学期
		var term string
		if len(parts) >= 2 && strings.Contains(parts[1], "-") {
			// 格式：courseCode_term_timestamp
			if len(parts) >= 3 && strings.Contains(parts[1], "-") {
				term = parts[1]
			} else {
				// 格式：term_timestamp
				term = parts[0]
			}
		}

		if term == "" {
			continue
		}

		// 确保目标目录存在
		dstDir := filepath.Join("data", "grades", term)
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}

		dstPath := filepath.Join(dstDir, courseCode+".csv")
		err = copyFile(csvFile, dstPath)
		if err != nil {
			return fmt.Errorf("恢复文件[%s]失败: %w", courseCode, err)
		}

		// 重新加载该课程的成绩到数据库
		rows, err := csvRepo.LoadCourseGrades(term, courseCode)
		if err != nil {
			return fmt.Errorf("读取恢复的成绩文件[%s]失败: %w", courseCode, err)
		}

		course, err := repository.GetCourseByCode(courseCode)
		if err != nil {
			continue
		}

		for _, row := range rows {
			student, err := repository.GetStudentByStudentID(row.StudentID)
			if err != nil {
				continue
			}

			// 删除旧成绩并插入新成绩（这里简化处理：仅在不存在时插入）
			exists, _ := repository.GradeExists(student.ID, course.ID)
			if exists {
				continue
			}

			// 创建成绩记录（批量导入逻辑在 load_grade_csv_service.go 中）
			// 此处利用 LoadGradesFromCSV 完成
		}
	}

	// 重新加载所有 CSV 数据到数据库
	return LoadAllCSVData()
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = sourceFile.Close() }()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destFile.Close() }()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
