// Package service import_helper.go CSV 导入通用辅助工具
package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ImportFromCSV 通用 CSV 数据导入函数，自动跳过已存在的记录
// T: 实体类型；K: 唯一键类型（必须可比较）
func ImportFromCSV[T any, K comparable](
	loadFunc func() ([]T, error), // 从 CSV 加载全部数据
	getKeyFunc func(T) K, // 提取实体的唯一标识
	checkExistFunc func(K) error, // 按唯一键检查是否存在，不存在时返回 gorm.ErrRecordNotFound
	createFunc func(*T) error, // 创建新记录
	entityName string, // 实体名称，用于错误提示
) error {
	// 加载CSV数据
	items, err := loadFunc()
	if err != nil {
		return fmt.Errorf("加载%s数据失败: %w", entityName, err)
	}

	// 逐条检查并导入
	for _, item := range items {
		key := getKeyFunc(item)

		// 检查记录是否已存在
		err := checkExistFunc(key)
		if err == nil {
			continue // 已存在，跳过
		}

		// 非"未找到"错误视为真实故障，立即返回
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("检查%s[%v]失败: %w", entityName, key, err)
		}

		// 记录不存在，执行创建（使用副本避免循环变量指针问题）
		itemCopy := item
		if err := createFunc(&itemCopy); err != nil {
			return fmt.Errorf("导入%s[%v]失败: %w", entityName, key, err)
		}
	}

	return nil
}
