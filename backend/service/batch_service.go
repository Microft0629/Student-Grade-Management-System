// Package service batch_service.go 批量操作服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"Student-Grade-Management-System/backend/utils"
	"fmt"
	"strings"
	"time"
)

// BatchAdjustResult 批量加减分操作结果
type BatchAdjustResult struct {
	AffectedCount int // 受影响记录数
	Details       []string // 逐条详情
}

// BatchAdjustScores 按课程+分数段对成绩进行批量加减分调整
func BatchAdjustScores(courseID uint, minScore float64, maxScore float64, delta float64) (*BatchAdjustResult, error) {
	var grades []model.Grade
	err := config.DB.
		Where("course_id = ? AND score >= ? AND score <= ?", courseID, minScore, maxScore).
		Find(&grades).Error
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)

	result := &BatchAdjustResult{}
	for _, grade := range grades {
		oldScore := grade.Score
		newScore := oldScore + delta

		// 调整后分数范围校验
		if newScore < 0 {
			newScore = 0
		}
		if newScore > 100 {
			newScore = 100
		}

		grade.Score = newScore
		grade.GradePoint = utils.CalculateGradePoint(newScore)

		err = config.DB.Save(&grade).Error
		if err != nil {
			return result, fmt.Errorf("更新成绩失败[ID=%d]: %w", grade.ID, err)
		}

		result.AffectedCount++
		result.Details = append(result.Details, fmt.Sprintf(
			"%s(%s) %s: %.1f -> %.1f",
			grade.Student.Name,
			grade.Student.StudentID,
			grade.Course.CourseName,
			oldScore,
			newScore,
		))

		// 记录操作日志
		LogOperation(model.OperationLog{
			Time:      time.Now(),
			Action:    "批量调整",
			Student:   grade.Student.Name,
			StudentID: grade.Student.StudentID,
			Course:    grade.Course.CourseName,
			Term:      grade.Course.Term,
			OldScore:  oldScore,
			NewScore:  newScore,
			Detail:    fmt.Sprintf("分数段[%.0f-%.0f]批量%+.0f分", minScore, maxScore, delta),
		})
	}

	if result.AffectedCount > 0 {
		err = SyncGradesToCSV()
		if err != nil {
			return result, fmt.Errorf("同步CSV失败: %w", err)
		}
	}

	return result, nil
}

// AggregatedGrade 跨课程/跨学期的成绩汇总记录
type AggregatedGrade struct {
	StudentID   string  // 学号
	StudentName string  // 姓名
	CourseName  string  // 课程名称
	Term        string  // 学期
	Score       float64 // 分数
	GradePoint  float64 // 绩点
	Credit      float64 // 学分
}

// AggregateGrades 跨课程/跨学期的成绩数据汇总
func AggregateGrades(term string, courseKeyword string) ([]AggregatedGrade, error) {
	var grades []model.Grade
	query := config.DB

	if term != "" {
		query = query.Where("courses.term = ?", term)
	}
	if courseKeyword != "" {
		query = query.Joins("JOIN courses ON courses.id = grades.course_id").
			Where("courses.term = ? OR courses.course_name LIKE ?", term, "%"+courseKeyword+"%")
	}

	// 简化查询：先获取全部再内存筛选
	err := config.DB.Find(&grades).Error
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)

	var result []AggregatedGrade
	for _, grade := range grades {
		// 学期筛选
		if term != "" && grade.Course.Term != term {
			continue
		}
		// 课程关键词筛选
		if courseKeyword != "" && !contains(grade.Course.CourseName, courseKeyword) {
			continue
		}

		result = append(result, AggregatedGrade{
			StudentID:   grade.Student.StudentID,
			StudentName: grade.Student.Name,
			CourseName:  grade.Course.CourseName,
			Term:        grade.Course.Term,
			Score:       grade.Score,
			GradePoint:  grade.GradePoint,
			Credit:      grade.Course.Credit,
		})
	}

	return result, nil
}

// ExportTranscript 导出标准化成绩单
func ExportTranscript(term string) (string, error) {
	var grades []model.Grade
	query := config.DB

	if term != "" {
		query = query.Joins("JOIN courses ON courses.id = grades.course_id").
			Where("courses.term = ?", term)
	}

	err := query.Find(&grades).Error
	if err != nil {
		return "", fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)

	// 按学生分组
	studentMap := make(map[uint][]model.Grade)
	var studentIDs []uint
	for _, grade := range grades {
		if _, exists := studentMap[grade.StudentID]; !exists {
			studentIDs = append(studentIDs, grade.StudentID)
		}
		studentMap[grade.StudentID] = append(studentMap[grade.StudentID], grade)
	}

	// 生成文字版成绩单（标签格式，避免中文字符列对齐问题）
	var sb strings.Builder
	sb.WriteString("========================================\n")
	sb.WriteString("        标 准 化 成 绩 单\n")
	if term != "" {
		sb.WriteString(fmt.Sprintf("        学期：%s\n", term))
	}
	sb.WriteString("========================================\n\n")

	for _, sid := range studentIDs {
		studentGrades := studentMap[sid]
		if len(studentGrades) == 0 {
			continue
		}

		student := studentGrades[0].Student
		sb.WriteString(fmt.Sprintf("▸ 学号: %s  姓名: %s\n", student.StudentID, student.Name))

		var totalScore, totalCredit float64
		for _, grade := range studentGrades {
			sb.WriteString(fmt.Sprintf("  %s | 学分: %.1f | 分数: %.1f | 绩点: %.1f\n",
				grade.Course.CourseName,
				grade.Course.Credit,
				grade.Score,
				grade.GradePoint,
			))
			totalScore += grade.Score
			totalCredit += grade.Course.Credit
		}

		avgScore := totalScore / float64(len(studentGrades))
		gpa := utils.CalculateStudentGPA(studentGrades)

		sb.WriteString(fmt.Sprintf("  → 平均分: %.1f | 绩点: %.2f | 总学分: %.1f\n",
			avgScore, gpa, totalCredit))
		sb.WriteString("\n")
	}

	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("  打印时间：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("========================================\n")

	return sb.String(), nil
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
