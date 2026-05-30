// Package csv grade_csv.go 成绩 CSV 持久化存储
package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// GradeCSV 单条成绩记录的 CSV 序列化结构
type GradeCSV struct {
	StudentID string // 学号
	Name      string // 姓名
	Score     float64 // 分数
}

// SaveCourseGrades 将某学期某课程的成绩列表写入 CSV 文件
func SaveCourseGrades(
	term string,
	courseCode string,
	grades []GradeCSV,
) error {
	// 按学期创建成绩存储目录，不存在则自动创建
	dir := filepath.Join("data", "grades", term)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, courseCode+".csv") // 以课程代码命名 CSV 文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// 关闭文件失败时仅忽略，避免覆盖主流程错误
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file) // 初始化 CSV 写入器
	defer writer.Flush()          // 确保函数返回前将缓冲区数据刷新至文件

	// 写入 CSV 表头
	err = writer.Write([]string{
		"StudentID",
		"Name",
		"Score",
	})
	if err != nil {
		return err
	}

	// 逐条写入成绩记录
	for _, grade := range grades {
		record := []string{
			grade.StudentID,
			grade.Name,
			fmt.Sprintf("%.2f", grade.Score),
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadCourseGrades 读取指定学期和课程的成绩CSV文件
func LoadCourseGrades(term, courseCode string) ([]GradeCSV, error) {
	var grades []GradeCSV

	filePath := filepath.Join("data", "grades", term, courseCode+".csv")

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return grades, nil
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 { // 跳过表头
			continue
		}

		score, _ := strconv.ParseFloat(record[2], 64)

		grades = append(grades, GradeCSV{
			StudentID: record[0],
			Name:      record[1],
			Score:     score,
		})
	}

	return grades, nil
}
