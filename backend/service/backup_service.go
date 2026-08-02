// Package service backup_service.go 数据备份与恢复服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	csvRepo "Student-Grade-Management-System/backend/repository/csv"
	"Student-Grade-Management-System/backend/utils"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
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
// 备份目录名支持两种格式：term_timestamp（按学期备份）和 courseCode_term_timestamp（按课程备份）。
// 恢复采用覆盖语义：备份中每个课程对应的现有成绩会被替换为备份文件中的成绩。
func RestoreFromBackup(backupName string) error {
	if backupName == "" || filepath.Base(backupName) != backupName || strings.Contains(backupName, "..") {
		return fmt.Errorf("备份名称不合法: %s", backupName)
	}

	backupDir := filepath.Join(backupRoot, backupName)

	info, err := os.Stat(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("备份[%s]不存在", backupName)
		}
		return fmt.Errorf("读取备份目录失败: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("备份[%s]不是目录", backupName)
	}

	// 从备份目录名解析学期信息（格式：term_timestamp 或 courseCode_term_timestamp）
	term, err := parseBackupTerm(backupName)
	if err != nil {
		return err
	}

	// 扫描备份目录下的 CSV 文件
	csvFiles, err := filepath.Glob(filepath.Join(backupDir, "*.csv"))
	if err != nil {
		return fmt.Errorf("扫描备份文件失败: %w", err)
	}
	if len(csvFiles) == 0 {
		return fmt.Errorf("备份[%s]中没有成绩文件", backupName)
	}

	restoredCount := 0
	for _, csvFile := range csvFiles {
		// 从文件名提取课程代码
		courseCode := strings.TrimSuffix(filepath.Base(csvFile), ".csv")

		// 课程不存在则跳过（学生/课程主数据不在成绩备份范围内）
		course, err := repository.GetCourseByCode(courseCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return fmt.Errorf("查询课程[%s]失败: %w", courseCode, err)
		}

		// 确保目标目录存在
		dstDir := filepath.Join("data", "grades", term)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}

		dstPath := filepath.Join(dstDir, courseCode+".csv")
		if err := copyFile(csvFile, dstPath); err != nil {
			return fmt.Errorf("恢复文件[%s]失败: %w", courseCode, err)
		}

		// 从已复制回的成绩文件读取备份行
		rows, err := csvRepo.LoadCourseGrades(term, courseCode)
		if err != nil {
			return fmt.Errorf("读取恢复的成绩文件[%s]失败: %w", courseCode, err)
		}

		// 覆盖语义：先删除该课程现有成绩，再按备份文件重新写入
		if err := config.DB.Where("course_id = ?", course.ID).Delete(&model.Grade{}).Error; err != nil {
			return fmt.Errorf("删除课程[%s]现有成绩失败: %w", courseCode, err)
		}

		for _, row := range rows {
			student, err := repository.GetStudentByStudentID(row.StudentID)
			if err != nil {
				continue // 学生不存在则跳过该条记录
			}

			grade := model.Grade{
				StudentID:   student.ID,
				CourseID:    course.ID,
				Score:       row.Score,
				GradePoint:  utils.CalculateGradePoint(row.Score),
				CreatorName: "admin",
			}
			if err := repository.CreateGrade(&grade); err != nil {
				return fmt.Errorf("恢复成绩[学号%s-课程%s]失败: %w", row.StudentID, courseCode, err)
			}
		}
		restoredCount++
	}

	if restoredCount == 0 {
		return fmt.Errorf("备份[%s]中没有可恢复的成绩数据（对应课程均不存在）", backupName)
	}

	// 从数据库重新同步全部成绩 CSV，保证与数据库一致
	return SyncGradesToCSV()
}

// parseBackupTerm 从备份目录名解析学期
// term_timestamp：parts[0] 为学期；courseCode_term_timestamp：parts[1] 为学期
func parseBackupTerm(backupName string) (string, error) {
	parts := strings.Split(backupName, "_")
	if len(parts) < 2 {
		return "", fmt.Errorf("备份名称格式不正确: %s", backupName)
	}
	if strings.Contains(parts[0], "-") {
		return parts[0], nil
	}
	if strings.Contains(parts[1], "-") {
		return parts[1], nil
	}
	return "", fmt.Errorf("无法从备份名称中解析学期: %s", backupName)
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
