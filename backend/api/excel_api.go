// Package api excel_api.go Excel 导出 API
package api

import (
	"Student-Grade-Management-System/backend/service"
	"os/exec"
	"runtime"
)

// ExcelAPI Excel 导出 API 结构体
type ExcelAPI struct{}

// NewExcelAPI 创建并返回一个新的 ExcelAPI 实例
func NewExcelAPI() *ExcelAPI {
	return &ExcelAPI{}
}

// ExportCourseStats 导出单课程统计 Excel 文件并打开，返回文件路径
func (e *ExcelAPI) ExportCourseStats(courseID uint) (string, error) {
	path, err := service.ExportCourseStatsExcel(courseID)
	if err != nil {
		return "", err
	}
	openFile(path)
	return path, nil
}

// ExportStudentStats 导出单学生统计 Excel 文件并打开，返回文件路径
func (e *ExcelAPI) ExportStudentStats(studentID uint) (string, error) {
	path, err := service.ExportStudentStatsExcel(studentID)
	if err != nil {
		return "", err
	}
	openFile(path)
	return path, nil
}

// openFile 使用系统默认程序打开文件
func openFile(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	_ = cmd.Start()
}
