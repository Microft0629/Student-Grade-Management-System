package service

import "testing"

func TestParseErrorLogLine(t *testing.T) {
	line := "[2026-08-02 10:00:00] 学生=2024001 课程=10871 分数=87.0 原因=成绩必须在0-100之间，当前值: 87.0"
	entry, ok := parseErrorLogLine(line)
	if !ok {
		t.Fatal("parseErrorLogLine 应能解析合法行")
	}
	if entry.Time.Format("2006-01-02 15:04:05") != "2026-08-02 10:00:00" {
		t.Errorf("时间解析错误: %v", entry.Time)
	}
	if entry.Student != "2024001" {
		t.Errorf("学生解析错误: %q", entry.Student)
	}
	if entry.Course != "10871" {
		t.Errorf("课程解析错误: %q", entry.Course)
	}
	if entry.Score != 87.0 {
		t.Errorf("分数解析错误: %v", entry.Score)
	}
	if entry.Reason != "成绩必须在0-100之间，当前值: 87.0" {
		t.Errorf("原因解析错误: %q", entry.Reason)
	}
}

func TestParseErrorLogLineInvalid(t *testing.T) {
	invalidLines := []string{
		"",
		"not a log line",
		"[bad-time] 学生=2024001 课程=10871 分数=87.0 原因=xx",
		"[2026-08-02 10:00:00] 学生=2024001 课程=10871 分数=notnum 原因=xx",
		"[2026-08-02 10:00:00] 学生=2024001 课程=10871",
	}
	for _, line := range invalidLines {
		if _, ok := parseErrorLogLine(line); ok {
			t.Errorf("非法行不应解析成功: %q", line)
		}
	}
}
