// Package csv atomic.go CSV 原子写入辅助
package csv

import (
	"os"
	"path/filepath"
)

// writeCSVAtomic 将内容先写入同目录临时文件，再原子替换目标文件，
// 避免写入中途失败时破坏原文件。
func writeCSVAtomic(path string, write func(*os.File) error) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if err := write(tmp); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
