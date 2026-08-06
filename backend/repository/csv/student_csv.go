// Package csv student_csv.go 学生 CSV 持久化存储
package csv

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"encoding/csv"
	"os"
	"path/filepath"
)

// studentCSVPath 学生数据 CSV 文件的存储路径
func studentCSVPath() string {
	return filepath.Join(utils.DataDir(), "students.csv")
}

// SaveStudents 将学生列表数据导出并保存至 CSV 文件，包含表头及全部学生记录
func SaveStudents(students []model.Student) error {
	return writeCSVAtomic(studentCSVPath(), func(file *os.File) error {
		writer := csv.NewWriter(file)

		if err := writer.Write([]string{
			"StudentID",
			"Name",
			"Gender",
			"ClassName",
			"Major",
		}); err != nil {
			return err
		}

		for _, student := range students {
			if err := writer.Write([]string{
				student.StudentID,
				student.Name,
				student.Gender,
				student.ClassName,
				student.Major,
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	})
}

// LoadStudents 从 CSV 文件中加载学生数据，若文件不存在则返回空列表
func LoadStudents() ([]model.Student, error) {
	var students []model.Student

	// 打开学生数据 CSV 文件
	file, err := os.Open(studentCSVPath())
	if err != nil {
		// 文件不存在时视为无数据，返回空列表而非错误
		if os.IsNotExist(err) {
			return students, nil
		}
		return nil, err
	}

	// 关闭文件失败时仅忽略，避免覆盖主流程错误
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// 初始化 CSV 读取器并一次性读取全部记录
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 遍历所有记录，跳过首行表头，逐条解析为学生结构体
	for i, record := range records {
		if i == 0 {
			continue // 跳过CSV表头行
		}

		var student model.Student

		// 按列顺序赋值学生各字段
		student.StudentID = record[0]
		student.Name = record[1]
		student.Gender = record[2]
		student.ClassName = record[3]
		student.Major = record[4]

		students = append(
			students,
			student,
		)
	}

	return students, nil
}
