// Package service batch_service.go 批量操作服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"Student-Grade-Management-System/backend/utils"
	"errors"
	"fmt"
	"time"
)

// BatchAdjustResult 批量加减分操作结果
type BatchAdjustResult struct {
	AffectedCount int // 受影响记录数
	Details       []string // 逐条详情
}

// BatchAdjustScores 按课程+分数段对成绩进行批量加减分调整
func BatchAdjustScores(courseID uint, minScore float64, maxScore float64, delta float64) (*BatchAdjustResult, error) {
	if err := RequireLogin(); err != nil {
		return nil, err
	}

	// 课程必须存在
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("课程不存在")
	}

	// 分数范围校验
	if minScore < 0 || maxScore > 100 || minScore > maxScore {
		return nil, errors.New("分数范围无效：需在0-100之间且最小值不能大于最大值")
	}

	// 老师只能调整自己创建的课程，管理员不受限制
	if !IsAdmin() && course.CreatorName != CurrentOperator() {
		return nil, errors.New("只能调整自己创建的课程的成绩")
	}

	var grades []model.Grade
	err = config.DB.
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

		// 只更新分数和绩点字段，避免 GORM 连带保存 Student/Course 关联导致业务字段被覆盖
		err = config.DB.Model(&model.Grade{}).
			Where("id = ?", grade.ID).
			Updates(map[string]interface{}{
				"score":       newScore,
				"grade_point": utils.CalculateGradePoint(newScore),
			}).Error
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
	query := config.DB.Model(&model.Grade{}).
		Select("grades.*").
		Joins("JOIN students ON students.id = grades.student_id").
		Joins("JOIN courses ON courses.id = grades.course_id")

	if term != "" {
		query = query.Where("courses.term = ?", term)
	}
	if courseKeyword != "" {
		query = query.Where("courses.course_name LIKE ?", "%"+courseKeyword+"%")
	}

	if err := query.Find(&grades).Error; err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}
	repository.LoadAssociations(grades)

	result := make([]AggregatedGrade, 0, len(grades))
	for _, grade := range grades {
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
