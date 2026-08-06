// Package csv course_csv.go 课程 CSV 持久化存储
package csv

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

// courseCSVPath 课程数据 CSV 文件的存储路径
func courseCSVPath() string {
	return filepath.Join(utils.DataDir(), "courses.csv")
}

// SaveCourses 将课程列表序列化并写入 CSV 文件，用于数据持久化备份
func SaveCourses(courses []model.Course) error {
	return writeCSVAtomic(courseCSVPath(), func(file *os.File) error {
		writer := csv.NewWriter(file)

		if err := writer.Write([]string{
			"CourseCode",
			"CourseName",
			"Term",
			"Credit",
			"Teacher",
		}); err != nil {
			return err
		}

		for _, course := range courses {
			if err := writer.Write([]string{
				course.CourseCode,
				course.CourseName,
				course.Term,
				strconv.FormatFloat(course.Credit, 'f', 1, 64),
				course.Teacher,
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	})
}

// LoadCourses 从 CSV 文件中加载课程列表，若文件不存在则返回空切片
func LoadCourses() ([]model.Course, error) {
	var courses []model.Course

	// 打开课程数据 CSV 文件
	file, err := os.Open(courseCSVPath())
	if err != nil {
		// 文件不存在时视为正常情况，返回空列表而非错误
		if os.IsNotExist(err) {
			return courses, nil
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

	// 遍历所有记录，跳过首行表头，逐条解析为课程结构体
	for i, record := range records {
		if i == 0 {
			continue // 跳过CSV表头行
		}

		// 解析学分字段，忽略解析错误（默认为0）
		credit, _ := strconv.ParseFloat(
			record[3],
			64,
		)

		var course model.Course

		// 按列顺序赋值课程各字段
		course.CourseCode = record[0]
		course.CourseName = record[1]
		course.Term = record[2]
		course.Credit = credit
		course.Teacher = record[4]

		courses = append(
			courses,
			course,
		)
	}

	return courses, nil
}
