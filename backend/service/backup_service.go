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

// BackupByTerm 备份指定学期的全部成绩数据
func BackupByTerm(term string) (string, error) {
	if !IsAdmin() {
		return "", errors.New("仅管理员可备份数据")
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(utils.BackupDir(), fmt.Sprintf("%s_%s", term, timestamp))

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

		srcPath := filepath.Join(utils.DataDir(), "grades", f.Term, f.CourseCode+".csv")
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
	if !IsAdmin() {
		return "", errors.New("仅管理员可备份数据")
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(utils.BackupDir(), fmt.Sprintf("%s_%s_%s", courseCode, term, timestamp))

	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建备份目录失败: %w", err)
	}

	srcPath := filepath.Join(utils.DataDir(), "grades", term, courseCode+".csv")
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
	if !IsAdmin() {
		return nil, errors.New("仅管理员可查看备份")
	}

	_, err := os.Stat(utils.BackupDir())
	if os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(utils.BackupDir())
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
	if !IsAdmin() {
		return errors.New("仅管理员可恢复数据")
	}

	if backupName == "" || filepath.Base(backupName) != backupName || strings.Contains(backupName, "..") {
		return fmt.Errorf("备份名称不合法: %s", backupName)
	}

	backupDir := filepath.Join(utils.BackupDir(), backupName)

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

	// 恢复前自动备份当前成绩，避免恢复过程异常时数据无法找回
	if err := backupBeforeRestore(term, csvFiles); err != nil {
		return fmt.Errorf("恢复前自动备份失败: %w", err)
	}

	// 在事务中删除并重建成绩，避免中途失败留下半恢复状态
	restoredCount := 0
	err = config.DB.Transaction(func(tx *gorm.DB) error {
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

			// 直接从备份文件读取成绩行
			rows, err := csvRepo.LoadCourseGradesFromFile(csvFile)
			if err != nil {
				return fmt.Errorf("读取恢复的成绩文件[%s]失败: %w", courseCode, err)
			}

			// 覆盖语义：先删除该课程现有成绩，再按备份文件重新写入
			if err := tx.Where("course_id = ?", course.ID).Delete(&model.Grade{}).Error; err != nil {
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
				if err := tx.Create(&grade).Error; err != nil {
					return fmt.Errorf("恢复成绩[学号%s-课程%s]失败: %w", row.StudentID, courseCode, err)
				}
			}
			restoredCount++
		}
		return nil
	})
	if err != nil {
		return err
	}
	if restoredCount == 0 {
		return fmt.Errorf("备份[%s]中没有可恢复的成绩数据（对应课程均不存在）", backupName)
	}

	// 从数据库重新同步全部成绩 CSV；同步失败只告警，不视为恢复失败
	if err := SyncGradesToCSV(); err != nil {
		println("CSV同步失败[恢复]:", err.Error())
	}
	return nil
}

// backupBeforeRestore 把当前成绩文件复制到 backup/恢复前_{学期}_{时间戳} 目录作为安全备份
func backupBeforeRestore(term string, csvFiles []string) error {
	safetyDir := filepath.Join(utils.BackupDir(), fmt.Sprintf("恢复前_%s_%s", term, time.Now().Format("20060102_150405")))
	if err := os.MkdirAll(safetyDir, 0755); err != nil {
		return err
	}

	copied := 0
	for _, csvFile := range csvFiles {
		courseCode := strings.TrimSuffix(filepath.Base(csvFile), ".csv")
		srcPath := filepath.Join(utils.DataDir(), "grades", term, courseCode+".csv")
		if _, err := os.Stat(srcPath); err != nil {
			continue
		}
		if err := copyFile(srcPath, filepath.Join(safetyDir, courseCode+".csv")); err != nil {
			return err
		}
		copied++
	}

	if copied == 0 {
		_ = os.RemoveAll(safetyDir) // 当前没有可备份的成绩文件，清理空目录
	}
	return nil
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
