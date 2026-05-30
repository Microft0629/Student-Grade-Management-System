// Package csv grade_scan.go 成绩 CSV 目录扫描工具
package csv

import (
	"os"
	"path/filepath"
	"strings"
)

// GradeFileInfo 成绩文件元数据
type GradeFileInfo struct {
	Term       string // 学期
	CourseCode string // 课程代码
}

// ScanGradeFiles 扫描 data/grades 下所有 CSV 成绩文件并解析元数据
func ScanGradeFiles() ([]GradeFileInfo, error) {
	var result []GradeFileInfo

	root := filepath.Join("data", "grades")
	_, err := os.Stat(root)

	// 目录不存在时返回空切片，不视为错误
	if os.IsNotExist(err) {
		return result, nil
	}

	err = filepath.Walk(
		root,
		func(
			path string,
			info os.FileInfo,
			err error,
		) error {
			if err != nil {
				return err
			}
			// 跳过目录
			if info.IsDir() {
				return nil
			}
			// 仅处理 .csv 文件
			if filepath.Ext(path) != ".csv" {
				return nil
			}
			// 从父目录名提取学期
			term :=
				filepath.Base(
					filepath.Dir(path),
				)
			// 从文件名提取课程代码
			courseCode :=
				strings.TrimSuffix(
					info.Name(),
					".csv",
				)
			result = append(
				result,
				GradeFileInfo{
					Term:       term,
					CourseCode: courseCode,
				},
			)
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
