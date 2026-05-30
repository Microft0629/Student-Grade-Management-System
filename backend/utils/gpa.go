// Package utils gpa.go GPA 计算工具
package utils

import (
	"Student-Grade-Management-System/backend/model"
	"math"
)

// CalculateGradePoint 根据成绩计算单科绩点
func CalculateGradePoint(score float64) float64 {
	var gp float64
	if score >= 60 {
		gp = score/10.0 - 5.0
	} else {
		gp = 0
	}
	return math.Round(gp*10) / 10 // 保留一位小数
}

// CalculateStudentGPA 计算学生 GPA
func CalculateStudentGPA(grades []model.Grade) float64 {
	var totalPoints, totalCredits float64
	// 计算：1.（单科成绩绩点 * 课程学分）的总和；2.课程学分的总和
	for _, grade := range grades {
		credit := grade.Course.Credit
		totalPoints += grade.GradePoint * credit
		totalCredits += credit
	}
	if totalCredits == 0 {
		return 0
	}
	gpa := totalPoints / totalCredits
	return math.Round(gpa*100) / 100 // 保留两位小数
}
