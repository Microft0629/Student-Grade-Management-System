// Package csv grade_csv.go 成绩 CSV 持久化存储
package csv

import (
	"Student-Grade-Management-System/backend/utils"
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
	filePath := filepath.Join(utils.DataDir(), "grades", term, courseCode+".csv")
	return writeCSVAtomic(filePath, func(file *os.File) error {
		writer := csv.NewWriter(file)

		if err := writer.Write([]string{
			"StudentID",
			"Name",
			"Score",
		}); err != nil {
			return err
		}

		for _, grade := range grades {
			if err := writer.Write([]string{
				grade.StudentID,
				grade.Name,
				fmt.Sprintf("%.2f", grade.Score),
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	})
}

// LoadCourseGrades 读取指定学期和课程的成绩CSV文件
func LoadCourseGrades(term, courseCode string) ([]GradeCSV, error) {
	return LoadCourseGradesFromFile(filepath.Join(utils.DataDir(), "grades", term, courseCode+".csv"))
}

// LoadCourseGradesFromFile 从指定 CSV 文件读取成绩记录（文件不存在时返回空列表）
func LoadCourseGradesFromFile(filePath string) ([]GradeCSV, error) {
	var grades []GradeCSV

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
