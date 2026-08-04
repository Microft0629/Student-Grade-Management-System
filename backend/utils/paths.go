// Package utils paths.go 应用数据目录解析
package utils

import (
	"os"
	"path/filepath"
)

// AppBaseDir 返回应用数据根目录。
// 优先使用包含 data 或 database 目录的当前工作目录（兼容 wails dev 与从项目目录启动），
// 否则使用可执行文件所在目录，避免从其他工作目录启动时把数据写到错误位置。
func AppBaseDir() string {
	wd, _ := os.Getwd()
	if hasDataMarkers(wd) {
		return wd
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe)
	}
	return wd
}

// hasDataMarkers 判断目录下是否已存在 data 或 database 目录
func hasDataMarkers(dir string) bool {
	_, err1 := os.Stat(filepath.Join(dir, "data"))
	_, err2 := os.Stat(filepath.Join(dir, "database"))
	return err1 == nil || err2 == nil
}

// DatabaseDir 数据库目录
func DatabaseDir() string {
	return filepath.Join(AppBaseDir(), "database")
}

// DataDir 数据文件目录
func DataDir() string {
	return filepath.Join(AppBaseDir(), "data")
}

// BackupDir 备份目录
func BackupDir() string {
	return filepath.Join(AppBaseDir(), "backup")
}

// ExportDir 导出目录
func ExportDir() string {
	return filepath.Join(AppBaseDir(), "export")
}
