// Package service excel_service.go Excel 报表导出服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportCourseStatsExcel 导出单课程统计为 Excel 文件，返回文件路径
func ExportCourseStatsExcel(courseID uint) (string, error) {
	stats, err := GetCourseStatistics(courseID)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "课程统计"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	// 标题样式
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "1a1a2e"},
	})
	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "555555"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	valueStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 14, Bold: true, Color: "1a1a1a"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	distHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "ffffff"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"faad14"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	f.SetCellValue(sheet, "A1", "课程成绩统计分析")
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)

	// 基本信息
	info := [][]interface{}{
		{"课程名称", stats.CourseName, "", "学期", stats.Term},
		{"参考人数", stats.StudentCount, "", "及格率", fmt.Sprintf("%.1f%%", stats.PassRate)},
		{"平均分", fmt.Sprintf("%.1f", stats.AverageScore), "", "最高分", fmt.Sprintf("%.1f", stats.HighestScore)},
		{"最低分", fmt.Sprintf("%.1f", stats.LowestScore), "", "", ""},
	}

	for i, row := range info {
		r := i + 3
		f.SetCellStyle(sheet, cellName(1, r), cellName(1, r), labelStyle)
		f.SetCellStyle(sheet, cellName(2, r), cellName(2, r), valueStyle)
		f.SetCellStyle(sheet, cellName(4, r), cellName(4, r), labelStyle)
		f.SetCellStyle(sheet, cellName(5, r), cellName(5, r), valueStyle)
		for j, val := range row {
			if val != "" {
				f.SetCellValue(sheet, cellName(j+1, r), val)
			}
		}
	}

	// 分数段分布表
	distStart := 9
	distHeaders := []string{"分数段", "人数", "占比"}
	for i, h := range distHeaders {
		cell := cellName(i+1, distStart)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, distHeaderStyle)
	}
	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 12)
	f.SetColWidth(sheet, "C", "C", 12)

	order := []string{"90-100", "80-89", "70-79", "60-69", "0-59"}
	for i, k := range order {
		r := distStart + 1 + i
		count := stats.Distribution[k]
		pct := 0.0
		if stats.StudentCount > 0 {
			pct = float64(count) / float64(stats.StudentCount) * 100
		}
		f.SetCellValue(sheet, cellName(1, r), k)
		f.SetCellValue(sheet, cellName(2, r), count)
		f.SetCellValue(sheet, cellName(3, r), fmt.Sprintf("%.1f%%", pct))
		f.SetCellStyle(sheet, cellName(1, r), cellName(3, r), centerStyle)
	}

	// 页脚
	footerRow := distStart + len(order) + 2
	f.SetCellValue(sheet, cellName(1, footerRow), fmt.Sprintf("生成时间：%s", time.Now().Format("2006-01-02 15:04:05")))

	// 列宽
	f.SetColWidth(sheet, "A", "F", 16)

	filePath := filepath.Join("export", fmt.Sprintf("课程统计_%s_%s.xlsx", stats.CourseName, time.Now().Format("20060102_150405")))
	err = os.MkdirAll("export", 0755)
	if err != nil {
		return "", err
	}
	err = f.SaveAs(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// ExportStudentStatsExcel 导出单学生统计为 Excel 文件，返回文件路径
func ExportStudentStatsExcel(studentID uint) (string, error) {
	// 获取该学生的所有成绩
	var grades []model.Grade
	err := config.DB.
		Preload("Student").
		Preload("Course").
		Where("student_id = ?", studentID).
		Find(&grades).Error
	if err != nil {
		return "", fmt.Errorf("查询成绩失败: %w", err)
	}
	if len(grades) == 0 {
		return "", fmt.Errorf("该学生暂无成绩数据")
	}

	student := grades[0].Student

	// 获取排名
	rankings, err := GetStudentRanking()
	if err != nil {
		return "", err
	}
	var rank int
	for _, r := range rankings {
		if r.StudentID == student.StudentID {
			rank = r.Rank
			break
		}
	}

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "学生统计"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "1a1a2e"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "ffffff"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4a90d9"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "555555"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	valueStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 14, Bold: true, Color: "1a1a1a"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	f.SetCellValue(sheet, "A1", "学生成绩统计分析")
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)

	// 基本信息
	var totalScore, totalCredit float64
	for _, g := range grades {
		totalScore += g.Score
		totalCredit += g.Course.Credit
	}
	avgScore := totalScore / float64(len(grades))
	gpa := utils.CalculateStudentGPA(grades)

	info := [][]interface{}{
		{"学号", student.StudentID, "", "姓名", student.Name},
		{"班级", student.ClassName, "", "专业", student.Major},
		{"课程总数", len(grades), "", "总学分", fmt.Sprintf("%.1f", totalCredit)},
		{"总分", fmt.Sprintf("%.1f", totalScore), "", "平均分", fmt.Sprintf("%.1f", avgScore)},
		{"平均绩点", fmt.Sprintf("%.2f", gpa), "", "排名", fmt.Sprintf("第 %d 名", rank)},
	}

	for i, row := range info {
		r := i + 3
		f.SetCellStyle(sheet, cellName(1, r), cellName(1, r), labelStyle)
		f.SetCellStyle(sheet, cellName(2, r), cellName(2, r), valueStyle)
		f.SetCellStyle(sheet, cellName(4, r), cellName(4, r), labelStyle)
		f.SetCellStyle(sheet, cellName(5, r), cellName(5, r), valueStyle)
		for j, val := range row {
			if val != "" {
				f.SetCellValue(sheet, cellName(j+1, r), val)
			}
		}
	}

	// 各科成绩明细表
	detailStart := 10
	detailHeaders := []string{"序号", "课程名称", "学期", "学分", "分数", "绩点"}
	for i, h := range detailHeaders {
		cell := cellName(i+1, detailStart)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, g := range grades {
		r := detailStart + 1 + i
		f.SetCellValue(sheet, cellName(1, r), i+1)
		f.SetCellValue(sheet, cellName(2, r), g.Course.CourseName)
		f.SetCellValue(sheet, cellName(3, r), g.Course.Term)
		f.SetCellValue(sheet, cellName(4, r), g.Course.Credit)
		f.SetCellValue(sheet, cellName(5, r), g.Score)
		f.SetCellValue(sheet, cellName(6, r), g.GradePoint)
		f.SetCellStyle(sheet, cellName(1, r), cellName(6, r), centerStyle)
	}

	// 列宽
	f.SetColWidth(sheet, "A", "F", 16)

	// 页脚
	footerRow := detailStart + len(grades) + 2
	f.SetCellValue(sheet, cellName(1, footerRow), fmt.Sprintf("生成时间：%s", time.Now().Format("2006-01-02 15:04:05")))

	filePath := filepath.Join("export", fmt.Sprintf("学生统计_%s_%s.xlsx", student.Name, time.Now().Format("20060102_150405")))
	err = os.MkdirAll("export", 0755)
	if err != nil {
		return "", err
	}
	err = f.SaveAs(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// cellName 将列号和行号转换为 Excel 单元格名称（如 1,1 → A1）
func cellName(col, row int) string {
	name, _ := excelize.CoordinatesToCellName(col, row)
	return name
}
