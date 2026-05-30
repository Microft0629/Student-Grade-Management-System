// Package service gpa_rule_service.go 绩点换算规则管理服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"fmt"
	"os"
	"strings"
)

// GpaRulesPath 绩点换算规则文件路径
const GpaRulesPath = "data/gpa_rules.txt"

// SaveGpaRules 将绩点换算公式保存至 TXT 文件
func SaveGpaRules(formula string) error {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(GpaRulesPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString(formula)
	return err
}

// LoadGpaRules 从 TXT 文件读取绩点换算公式
func LoadGpaRules() (string, error) {
	data, err := os.ReadFile(GpaRulesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// GetDefaultGpaRules 返回系统默认绩点换算公式描述
func GetDefaultGpaRules() string {
	return "# 绩点换算公式\n" +
		"# score >= 60: 绩点 = round(score/10 - 5, 1)\n" +
		"# score <  60: 绩点 = 0\n"
}

// RecalculateAllGPA 批量重新计算所有学生绩点并更新成绩文件
func RecalculateAllGPA() (int, error) {
	var grades []model.Grade
	err := config.DB.Preload("Student").Preload("Course").Find(&grades).Error
	if err != nil {
		return 0, fmt.Errorf("查询成绩失败: %w", err)
	}

	count := 0
	for _, grade := range grades {
		newGP := utils.CalculateGradePoint(grade.Score)
		if grade.GradePoint != newGP {
			grade.GradePoint = newGP
			err = config.DB.Save(&grade).Error
			if err != nil {
				return count, fmt.Errorf("更新绩点失败[ID=%d]: %w", grade.ID, err)
			}
			count++
		}
	}

	if count > 0 {
		err = SyncGradesToCSV()
		if err != nil {
			return count, fmt.Errorf("同步CSV失败: %w", err)
		}
	}

	return count, nil
}
