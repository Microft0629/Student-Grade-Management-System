// Package csv student_csv.go 学生 CSV 持久化存储
package csv

import (
	"Student-Grade-Management-System/backend/model"
	"encoding/csv"
	"os"
)

// StudentCSVPath 学生数据 CSV 文件的存储路径
const StudentCSVPath = "data/students.csv"

// SaveStudents 将学生列表数据导出并保存至 CSV 文件，包含表头及全部学生记录
func SaveStudents(students []model.Student) error {
	// 确保 data 目录存在，若不存在则创建，权限设为 0755
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}

	// 创建或覆盖目标 CSV 文件
	file, err := os.Create(StudentCSVPath)
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

	// 写入 CSV 表头字段
	err = writer.Write([]string{
		"StudentID",
		"Name",
		"Gender",
		"ClassName",
		"Major",
	})
	if err != nil {
		return err
	}

	// 遍历学生列表，逐条转换为字符串记录并写入 CSV
	for _, student := range students {
		record := []string{
			student.StudentID,
			student.Name,
			student.Gender,
			student.ClassName,
			student.Major,
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadStudents 从 CSV 文件中加载学生数据，若文件不存在则返回空列表
func LoadStudents() ([]model.Student, error) {
	var students []model.Student

	// 打开学生数据 CSV 文件
	file, err := os.Open(StudentCSVPath)
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
