// Package api backup_api.go 数据备份恢复 API
package api

import "Student-Grade-Management-System/backend/service"

// BackupAPI 数据备份恢复 API 结构体
type BackupAPI struct{}

// NewBackupAPI 创建并返回一个新的 BackupAPI 实例
func NewBackupAPI() *BackupAPI {
	return &BackupAPI{}
}

// BackupByTerm 按学期备份成绩数据
func (b *BackupAPI) BackupByTerm(term string) (string, error) {
	return service.BackupByTerm(term)
}

// BackupByCourse 按课程备份成绩数据
func (b *BackupAPI) BackupByCourse(term string, courseCode string) (string, error) {
	return service.BackupByCourse(term, courseCode)
}

// ListBackups 列出所有备份
func (b *BackupAPI) ListBackups() ([]string, error) {
	return service.ListBackups()
}

// RestoreFromBackup 从备份恢复数据
func (b *BackupAPI) RestoreFromBackup(backupName string) error {
	return service.RestoreFromBackup(backupName)
}
