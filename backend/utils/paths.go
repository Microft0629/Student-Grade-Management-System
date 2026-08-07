// Package utils paths.go 应用数据目录解析
package utils

import (
	"os"
	"path/filepath"
)

// DataDirEnv 环境变量名，显式指定应用数据根目录（测试与部署可用）
const DataDirEnv = "SGMS_DATA_DIR"

// AppBaseDir 返回应用数据根目录，解析优先级：
//  1. 环境变量 SGMS_DATA_DIR（显式指定，测试/部署可用）
//  2. 可执行文件所在目录（若已包含 data/database，部署形态优先）
//  3. 当前工作目录（若已包含 data/database，兼容 wails dev 从项目目录启动）
//  4. 可执行文件所在目录 / 当前工作目录兜底
func AppBaseDir() string {
	if dir := os.Getenv(DataDirEnv); dir != "" {
		return filepath.Clean(dir)
	}

	exeDir := ""
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
		if hasDataMarkers(exeDir) {
			return exeDir
		}
	}

	wd, _ := os.Getwd()
	if hasDataMarkers(wd) {
		return wd
	}
	if exeDir != "" {
		return exeDir
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
