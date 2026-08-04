// Package service validation_service.go 数据校验与错误日志服务
package service

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"Student-Grade-Management-System/backend/utils"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// errorLogPath 校验错误日志文件路径
func errorLogPath() string {
	return filepath.Join(utils.DataDir(), "error.log")
}

// appendErrorLog 将校验错误追加写入错误日志文件
func appendErrorLog(entry model.ErrorLog) error {
	err := os.MkdirAll(utils.DataDir(), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(errorLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if err := appendErrorLog(entry); err != nil {
		println("记录校验错误日志失败:", err.Error())
	}
}

// ReadErrorLogs 读取所有错误日志
func ReadErrorLogs() ([]model.ErrorLog, error) {
	var logs []model.ErrorLog

	file, err := os.Open(errorLogPath())
	if err != nil {
		if os.IsNotExist(err) {
			return logs, nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if entry, ok := parseErrorLogLine(line); ok {
			logs = append(logs, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// parseErrorLogLine 解析单行错误日志
// 行格式：[2006-01-02 15:04:05] 学生=xxx 课程=xxx 分数=87.0 原因=xxx
func parseErrorLogLine(line string) (model.ErrorLog, bool) {
	var entry model.ErrorLog

	if !strings.HasPrefix(line, "[") {
		return entry, false
	}
	end := strings.Index(line, "]")
	if end <= 1 {
		return entry, false
	}

	t, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(line[1:end]))
	if err != nil {
		return entry, false
	}
	entry.Time = t

	rest := line[end+1:]
	var ok bool
	if entry.Student, ok = extractLogField(rest, "学生=", " 课程="); !ok {
		return model.ErrorLog{}, false
	}
	if entry.Course, ok = extractLogField(rest, "课程=", " 分数="); !ok {
		return model.ErrorLog{}, false
	}
	scoreStr, ok := extractLogField(rest, "分数=", " 原因=")
	if !ok {
		return model.ErrorLog{}, false
	}
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return model.ErrorLog{}, false
	}
	entry.Score = score
	entry.Reason, _ = extractLogField(rest, "原因=", "")
	return entry, true
}

// extractLogField 从日志行剩余部分中提取 "前缀...分隔符" 之间的字段值
func extractLogField(s, prefix, sep string) (string, bool) {
	idx := strings.Index(s, prefix)
	if idx < 0 {
		return "", false
	}
	rest := s[idx+len(prefix):]
	if sep == "" {
		return strings.TrimSpace(rest), true
	}
	sepIdx := strings.Index(rest, sep)
	if sepIdx < 0 {
		return "", false
	}
	return strings.TrimSpace(rest[:sepIdx]), true
}
