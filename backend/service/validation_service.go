// Package service validation_service.go 数据校验与错误日志服务
package service

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"fmt"
	"os"
	"time"
)

// ErrorLogPath 校验错误日志文件路径
const ErrorLogPath = "data/error.log"

// appendErrorLog 将校验错误追加写入错误日志文件
func appendErrorLog(entry model.ErrorLog) error {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(ErrorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	line := fmt.Sprintf("[%s] 学生=%s 课程=%s 分数=%.1f 原因=%s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Student,
		entry.Course,
		entry.Score,
		entry.Reason,
	)

	_, err = file.WriteString(line)
	return err
}

// ValidateScore 校验分数是否在 0-100 范围内
func ValidateScore(score float64) error {
	if score < 0 || score > 100 {
		return fmt.Errorf("成绩必须在0-100之间，当前值: %.1f", score)
	}
	return nil
}

// ValidateStudentExists 校验学号是否存在，返回对应的数据库学生ID
func ValidateStudentExists(studentID string) (uint, error) {
	student, err := repository.GetStudentByStudentID(studentID)
	if err != nil {
		return 0, fmt.Errorf("学号[%s]不存在", studentID)
	}
	return student.ID, nil
}

// LogValidationError 记录校验失败的错误日志
func LogValidationError(student, course string, score float64, reason string) {
	entry := model.ErrorLog{
		Time:    time.Now(),
		Student: student,
		Course:  course,
		Score:   score,
		Reason:  reason,
	}
	_ = appendErrorLog(entry)
}

// ReadErrorLogs 读取所有错误日志
func ReadErrorLogs() ([]model.ErrorLog, error) {
	var logs []model.ErrorLog

	file, err := os.Open(ErrorLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return logs, nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var entry model.ErrorLog
	for {
		_, err := fmt.Fscanf(file, "[%s] 学生=%s 课程=%s 分数=%f 原因=%[^\n]\n",
			&entry.Time, &entry.Student, &entry.Course, &entry.Score, &entry.Reason)
		if err != nil {
			break
		}
		logs = append(logs, entry)
	}

	return logs, nil
}
