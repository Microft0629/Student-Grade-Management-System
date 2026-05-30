// Package service statistics_service.go 成绩统计业务逻辑
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"fmt"
	"sort"
	"strings"
	"time"
)

// GetStudentStatistics 获取所有学生的成绩统计数据（含平均分、GPA、总学分及课程数）
func GetStudentStatistics() ([]model.StudentStatistics, error) {
	var grades []model.Grade
	err := config.DB.
		Preload("Student").
		Preload("Course").
		Find(&grades).Error
	if err != nil {
		return nil, err
	}

	studentGradesMap := make(map[uint][]model.Grade)
	for _, grade := range grades {
		studentGradesMap[grade.StudentID] = append(
			studentGradesMap[grade.StudentID],
			grade,
		)
	}

	var result []model.StudentStatistics

	for _, studentGrades := range studentGradesMap {
		if len(studentGrades) == 0 {
			continue
		}

		student := studentGrades[0].Student

		var totalScore float64
		var totalCredits float64

		for _, grade := range studentGrades {
			totalScore += grade.Score
			totalCredits += grade.Course.Credit
		}

		averageScore := totalScore / float64(len(studentGrades))
		gpa := utils.CalculateStudentGPA(studentGrades)

		result = append(result, model.StudentStatistics{
			StudentName:  student.Name,
			AverageScore: averageScore,
			GPA:          gpa,
			TotalCredits: totalCredits,
			CourseCount:  len(studentGrades),
		})
	}

	return result, nil
}

// GetCourseStatistics 对单课程统计平均分、及格率、分数段分布
func GetCourseStatistics(courseID uint) (*model.CourseStatistics, error) {
	var grades []model.Grade
	err := config.DB.
		Preload("Student").
		Preload("Course").
		Where("course_id = ?", courseID).
		Find(&grades).Error
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}

	if len(grades) == 0 {
		return nil, fmt.Errorf("该课程暂无成绩数据")
	}

	var totalScore, highestScore, lowestScore float64
	highestScore = -1
	lowestScore = 101
	passCount := 0
	dist := map[string]int{
		"0-59":   0,
		"60-69":  0,
		"70-79":  0,
		"80-89":  0,
		"90-100": 0,
	}

	for _, grade := range grades {
		s := grade.Score
		totalScore += s
		if s > highestScore {
			highestScore = s
		}
		if s < lowestScore {
			lowestScore = s
		}
		if s >= 60 {
			passCount++
		}
		switch {
		case s < 60:
			dist["0-59"]++
		case s < 70:
			dist["60-69"]++
		case s < 80:
			dist["70-79"]++
		case s < 90:
			dist["80-89"]++
		default:
			dist["90-100"]++
		}
	}

	stats := &model.CourseStatistics{
		CourseName:   grades[0].Course.CourseName,
		Term:         grades[0].Course.Term,
		StudentCount: len(grades),
		AverageScore: totalScore / float64(len(grades)),
		PassRate:     float64(passCount) / float64(len(grades)) * 100,
		HighestScore: highestScore,
		LowestScore:  lowestScore,
		Distribution: dist,
	}

	return stats, nil
}

// GetStudentRanking 对学生统计多课程总分、平均绩点并排名
func GetStudentRanking() ([]model.StudentRanking, error) {
	var grades []model.Grade
	err := config.DB.
		Preload("Student").
		Preload("Course").
		Find(&grades).Error
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}

	type studentData struct {
		student    model.Student
		totalScore float64
		gradeCount int
		grades     []model.Grade
	}
	studentMap := make(map[uint]*studentData)

	for _, grade := range grades {
		if _, exists := studentMap[grade.StudentID]; !exists {
			studentMap[grade.StudentID] = &studentData{
				student: grade.Student,
				grades:  []model.Grade{},
			}
		}
		sd := studentMap[grade.StudentID]
		sd.totalScore += grade.Score
		sd.gradeCount++
		sd.grades = append(sd.grades, grade)
	}

	var rankings []model.StudentRanking
	for _, sd := range studentMap {
		if sd.gradeCount == 0 {
			continue
		}
		rankings = append(rankings, model.StudentRanking{
			StudentName:  sd.student.Name,
			StudentID:    sd.student.StudentID,
			TotalScore:   sd.totalScore,
			AverageScore: sd.totalScore / float64(sd.gradeCount),
			GPA:          utils.CalculateStudentGPA(sd.grades),
			CourseCount:  sd.gradeCount,
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].GPA > rankings[j].GPA
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings, nil
}

// GenerateStatisticsReport 生成文字版统计报表（标签格式，避免中文字符列对齐问题）
func GenerateStatisticsReport(term string) (string, error) {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("      学 生 成 绩 统 计 报 表\n")
	if term != "" {
		sb.WriteString(fmt.Sprintf("      学期：%s\n", term))
	}
	sb.WriteString("========================================\n\n")

	rankings, err := GetStudentRanking()
	if err != nil {
		return "", err
	}

	// 一、学生排名
	if term != "" {
		var filtered []model.StudentRanking
		for _, r := range rankings {
			var grades []model.Grade
			config.DB.
				Preload("Student").
				Preload("Course").
				Where("student_id = (SELECT id FROM students WHERE student_id = ?)", r.StudentID).
				Where("course_id IN (SELECT id FROM courses WHERE term = ?)", term).
				Find(&grades)

			if len(grades) == 0 {
				continue
			}

			var totalScore float64
			for _, g := range grades {
				totalScore += g.Score
			}

			filtered = append(filtered, model.StudentRanking{
				StudentName:  r.StudentName,
				StudentID:    r.StudentID,
				TotalScore:   totalScore,
				AverageScore: totalScore / float64(len(grades)),
				GPA:          utils.CalculateStudentGPA(grades),
				CourseCount:  len(grades),
			})
		}

		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].GPA > filtered[j].GPA
		})
		for i := range filtered {
			filtered[i].Rank = i + 1
		}
		rankings = filtered
	}

	sb.WriteString("【一、学生综合排名（按平均绩点）】\n\n")
	for _, r := range rankings {
		medal := ""
		if r.Rank == 1 {
			medal = " [🥇]"
		} else if r.Rank == 2 {
			medal = " [🥈]"
		} else if r.Rank == 3 {
			medal = " [🥉]"
		}
		sb.WriteString(fmt.Sprintf(
			"第%d名%s  学号: %s  姓名: %s\n"+
				"  总分: %.1f  平均分: %.1f  绩点: %.2f  课程数: %d\n",
			r.Rank, medal, r.StudentID, r.StudentName,
			r.TotalScore, r.AverageScore, r.GPA, r.CourseCount,
		))
	}
	sb.WriteString("\n")

	// 二、课程统计
	sb.WriteString("【二、各课程成绩统计】\n\n")

	courses, err := getAllCoursesForTerm(term)
	if err != nil {
		return "", err
	}

	for _, course := range courses {
		stats, err := GetCourseStatistics(course.ID)
		if err != nil {
			continue
		}

		sb.WriteString(fmt.Sprintf("▸ %s（%s）\n", stats.CourseName, stats.Term))
		sb.WriteString(fmt.Sprintf("  参考人数: %-4d  平均分: %-6.1f  及格率: %-6.1f%%\n",
			stats.StudentCount, stats.AverageScore, stats.PassRate))
		sb.WriteString(fmt.Sprintf("  最高分: %-6.1f  最低分: %-6.1f\n",
			stats.HighestScore, stats.LowestScore))
		sb.WriteString(fmt.Sprintf("  分布: 0-59[%d]  60-69[%d]  70-79[%d]  80-89[%d]  90-100[%d]\n",
			stats.Distribution["0-59"],
			stats.Distribution["60-69"],
			stats.Distribution["70-79"],
			stats.Distribution["80-89"],
			stats.Distribution["90-100"],
		))
		sb.WriteString("\n")
	}

	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("  报表生成时间：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("========================================\n")

	return sb.String(), nil
}

func getAllCoursesForTerm(term string) ([]model.Course, error) {
	var courses []model.Course
	query := config.DB
	if term != "" {
		query = query.Where("term = ?", term)
	}
	err := query.Find(&courses).Error
	return courses, err
}
