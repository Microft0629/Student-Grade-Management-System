// Package service excel_service.go Excel 报表导出服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
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
		Where("student_id = ?", studentID).
		Find(&grades).Error
	if err != nil {
		return "", fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)
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

// ExportTranscriptExcel 导出标准化成绩单为 Excel 文件
func ExportTranscriptExcel(term string) (string, error) {
	var grades []model.Grade
	query := config.DB

	if term != "" {
		query = query.Joins("JOIN courses ON courses.id = grades.course_id").
			Where("courses.term = ?", term)
	}

	err := query.Find(&grades).Error
	if err != nil {
		return "", fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "成绩单"
	f.SetSheetName("Sheet1", sheet)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "1a1a2e"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "ffffff"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4a90d9"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "ffffff", Style: 1},
			{Type: "right", Color: "ffffff", Style: 1},
		},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "e0e0e0", Style: 1},
			{Type: "right", Color: "e0e0e0", Style: 1},
			{Type: "bottom", Color: "e0e0e0", Style: 1},
		},
	})
	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// 标题
	f.SetCellValue(sheet, "A1", "标准化成绩单")
	if term != "" {
		f.SetCellValue(sheet, "A1", fmt.Sprintf("标准化成绩单（%s）", term))
	}
	f.MergeCell(sheet, "A1", "G1")
	f.SetCellStyle(sheet, "A1", "G1", titleStyle)
	f.SetRowHeight(sheet, 1, 36)

	// 按学生分组
	studentMap := make(map[uint][]model.Grade)
	var studentIDs []uint
	for _, g := range grades {
		if _, exists := studentMap[g.StudentID]; !exists {
			studentIDs = append(studentIDs, g.StudentID)
		}
		studentMap[g.StudentID] = append(studentMap[g.StudentID], g)
	}

	row := 3
	headers := []string{"学号", "姓名", "课程名称", "学期", "学分", "分数", "绩点"}
	colWidths := []float64{14, 12, 22, 16, 10, 10, 10}

	for i, h := range headers {
		cell := cellName(i+1, row)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	for i, w := range colWidths {
		f.SetColWidth(sheet, cellName(i+1, 1), cellName(i+1, 1), w)
	}
	f.SetRowHeight(sheet, row, 24)
	row++

	for _, sid := range studentIDs {
		studentGrades := studentMap[sid]
		if len(studentGrades) == 0 {
			continue
		}
		student := studentGrades[0].Student

		var totalScore, totalCredit float64
		startRow := row

		for _, g := range studentGrades {
			f.SetCellValue(sheet, cellName(1, row), student.StudentID)
			f.SetCellValue(sheet, cellName(2, row), student.Name)
			f.SetCellValue(sheet, cellName(3, row), g.Course.CourseName)
			f.SetCellValue(sheet, cellName(4, row), g.Course.Term)
			f.SetCellValue(sheet, cellName(5, row), g.Course.Credit)
			f.SetCellValue(sheet, cellName(6, row), g.Score)
			f.SetCellValue(sheet, cellName(7, row), g.GradePoint)
			for c := 1; c <= 7; c++ {
				f.SetCellStyle(sheet, cellName(c, row), cellName(c, row), cellStyle)
			}
			totalScore += g.Score
			totalCredit += g.Course.Credit
			row++
		}

		// 汇总行
		avgScore := totalScore / float64(len(studentGrades))
		gpa := utils.CalculateStudentGPA(studentGrades)

		f.MergeCell(sheet, cellName(1, row), cellName(2, row))
		f.SetCellValue(sheet, cellName(1, row), "小计")
		f.SetCellValue(sheet, cellName(5, row), totalCredit)
		f.SetCellValue(sheet, cellName(6, row), fmt.Sprintf("均%.1f", avgScore))
		f.SetCellValue(sheet, cellName(7, row), fmt.Sprintf("%.2f", gpa))
		for c := 1; c <= 7; c++ {
			f.SetCellStyle(sheet, cellName(c, row), cellName(c, row), boldStyle)
		}
		row++

		// 空行分隔
		if startRow > 3+len(headers) {
			row++ // 学生之间加空行
		}
	}

	f.SetRowHeight(sheet, row, 8)
	row++
	f.SetCellValue(sheet, cellName(1, row), fmt.Sprintf("打印时间：%s", time.Now().Format("2006-01-02 15:04:05")))
	f.MergeCell(sheet, cellName(1, row), cellName(4, row))

	filePath := filepath.Join("export", fmt.Sprintf("成绩单_%s.xlsx", time.Now().Format("20060102_150405")))
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

// ExportOperationLogsExcel 导出全部操作日志为 Excel 文件
func ExportOperationLogsExcel() (string, error) {
	logs, err := GetAllOperationLogs()
	if err != nil {
		return "", fmt.Errorf("查询操作日志失败: %w", err)
	}

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "操作日志"
	f.SetSheetName("Sheet1", sheet)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "1a1a2e"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "ffffff"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4a90d9"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "e0e0e0", Style: 1},
			{Type: "right", Color: "e0e0e0", Style: 1},
			{Type: "bottom", Color: "e0e0e0", Style: 1},
		},
	})

	f.SetCellValue(sheet, "A1", "操作日志")
	f.MergeCell(sheet, "A1", "J1")
	f.SetCellStyle(sheet, "A1", "J1", titleStyle)
	f.SetRowHeight(sheet, 1, 36)

	headers := []string{"时间", "操作人", "操作类型", "学号", "学生", "课程", "学期", "旧分", "新分", "详情"}
	for i, h := range headers {
		cell := cellName(i+1, 3)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 3, 24)

	for i, log := range logs {
		row := 4 + i
		f.SetCellValue(sheet, cellName(1, row), log.Time.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheet, cellName(2, row), log.Operator)
		f.SetCellValue(sheet, cellName(3, row), log.Action)
		f.SetCellValue(sheet, cellName(4, row), log.StudentID)
		f.SetCellValue(sheet, cellName(5, row), log.Student)
		f.SetCellValue(sheet, cellName(6, row), log.Course)
		f.SetCellValue(sheet, cellName(7, row), log.Term)
		f.SetCellValue(sheet, cellName(8, row), log.OldScore)
		f.SetCellValue(sheet, cellName(9, row), log.NewScore)
		f.SetCellValue(sheet, cellName(10, row), log.Detail)
		for c := 1; c <= 10; c++ {
			f.SetCellStyle(sheet, cellName(c, row), cellName(c, row), cellStyle)
		}
	}

	cols := []float64{18, 12, 10, 12, 10, 16, 14, 8, 8, 24}
	for i, w := range cols {
		f.SetColWidth(sheet, cellName(i+1, 1), cellName(i+1, 1), w)
	}

	filePath := filepath.Join("export", fmt.Sprintf("操作日志_%s.xlsx", time.Now().Format("20060102_150405")))
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
