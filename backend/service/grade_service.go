// Package service grade_service.go 成绩业务逻辑层
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

// CreateGrade 根据分数计算绩点，并调用数据访问层创建成绩记录
func CreateGrade(grade *model.Grade) error {
	// 学生校验
	if grade.StudentID == 0 {
		return errors.New("请选择学生")
	}

	// 课程校验
	if grade.CourseID == 0 {
		return errors.New("请选择课程")
	}

	// 成绩范围校验
	if err := ValidateScore(grade.Score); err != nil {
		LogValidationError(
			fmt.Sprintf("ID=%d", grade.StudentID),
			fmt.Sprintf("ID=%d", grade.CourseID),
			grade.Score,
			err.Error(),
		)
		return err
	}

	// 成绩存在性校验
	exists, err :=
		repository.GradeExists(
			grade.StudentID,
			grade.CourseID,
		)
	if err != nil {
		return fmt.Errorf("检查成绩是否存在失败: %w", err)
	}
	if exists {
		return errors.New("该学生该课程成绩已存在")
	}

	// 自动计算绩点
	grade.GradePoint = utils.CalculateGradePoint(grade.Score)
	err = repository.CreateGrade(grade)
	if err != nil {
		return err
	}

	// 记录操作日志（需查询关联学生和课程信息）
	logCreateOperation(grade)

	return SyncGradesToCSV()
}

// UpdateGrade 修改已有成绩记录并重新计算绩点
func UpdateGrade(id uint, newScore float64) error {
	if err := ValidateScore(newScore); err != nil {
		return err
	}

	// 查询全部成绩记录以获取待修改项（含关联信息）
	grades, err := repository.GetAllGrades()
	if err != nil {
		return fmt.Errorf("查询成绩失败: %w", err)
	}

	var target *model.Grade
	for i := range grades {
		if grades[i].ID == id {
			target = &grades[i]
			break
		}
	}
	if target == nil {
		return errors.New("成绩记录不存在")
	}

	oldScore := target.Score
	target.Score = newScore
	target.GradePoint = utils.CalculateGradePoint(newScore)

	// 保存修改（通过直接使用 config.DB 保存，避免通过 repository 的局限）
	err = updateGradeDirect(target)
	if err != nil {
		return fmt.Errorf("更新成绩失败: %w", err)
	}

	// 记录操作日志
	LogOperation(model.OperationLog{
		Time:      time.Now(),
		Action:    "修改",
		Student:   target.Student.Name,
		StudentID: target.Student.StudentID,
		Course:    target.Course.CourseName,
		Term:      target.Course.Term,
		OldScore:  oldScore,
		NewScore:  newScore,
		Detail:    fmt.Sprintf("分数由 %.1f 修改为 %.1f", oldScore, newScore),
	})

	return SyncGradesToCSV()
}

// GetAllGrades 获取所有包含关联信息的成绩记录
func GetAllGrades() ([]model.Grade, error) {
	return repository.GetAllGrades()
}

// DeleteGrade 删除指定成绩记录并同步更新 CSV 备份文件
func DeleteGrade(id uint) error {
	// 先查询待删除记录以便记录日志
	grades, _ := repository.GetAllGrades()
	var target *model.Grade
	for i := range grades {
		if grades[i].ID == id {
			target = &grades[i]
			break
		}
	}

	// 从数据库中删除成绩
	err := repository.DeleteGrade(id)
	if err != nil {
		return err
	}

	// 记录操作日志
	if target != nil {
		LogOperation(model.OperationLog{
			Time:      time.Now(),
			Action:    "删除",
			Student:   target.Student.Name,
			StudentID: target.Student.StudentID,
			Course:    target.Course.CourseName,
			Term:      target.Course.Term,
			OldScore:  target.Score,
			NewScore:  0,
			Detail:    "删除成绩记录",
		})
	}

	// 数据库删除成功后，重新同步全量成绩到 CSV 以保持数据一致
	return SyncGradesToCSV()
}

// SearchGrades 查询成绩
func SearchGrades(
	studentKeyword string,
	courseKeyword string,
	term string,
) ([]model.Grade, error) {
	return repository.SearchGrades(
		studentKeyword,
		courseKeyword,
		term,
	)
}

// BatchImportGrades 批量导入成绩，校验失败不写入并提示错误原因
func BatchImportGrades(grades []model.Grade) (int, []string, error) {
	var errors_ []string
	successCount := 0

	for i, grade := range grades {
		// 分数范围校验
		if err := ValidateScore(grade.Score); err != nil {
			errMsg := fmt.Sprintf("第%d条: %v", i+1, err)
			errors_ = append(errors_, errMsg)
			LogValidationError(
				fmt.Sprintf("ID=%d", grade.StudentID),
				fmt.Sprintf("ID=%d", grade.CourseID),
				grade.Score,
				err.Error(),
			)
			continue
		}

		// 成绩存在性校验
		exists, err := repository.GradeExists(grade.StudentID, grade.CourseID)
		if err != nil {
			errors_ = append(errors_, fmt.Sprintf("第%d条: 检查失败 - %v", i+1, err))
			continue
		}
		if exists {
			errors_ = append(errors_, fmt.Sprintf("第%d条: 该学生该课程成绩已存在", i+1))
			continue
		}

		// 计算绩点并写入
		grade.GradePoint = utils.CalculateGradePoint(grade.Score)
		err = repository.CreateGrade(&grade)
		if err != nil {
			errors_ = append(errors_, fmt.Sprintf("第%d条: 写入失败 - %v", i+1, err))
			continue
		}
		successCount++
	}

	if successCount > 0 {
		err := SyncGradesToCSV()
		if err != nil {
			return successCount, errors_, fmt.Errorf("同步CSV失败: %w", err)
		}
	}

	return successCount, errors_, nil
}

// logCreateOperation 记录新增成绩的操作日志
func logCreateOperation(grade *model.Grade) {
	// 查询关联的学生和课程信息
	student := grade.Student
	course := grade.Course

	LogOperation(model.OperationLog{
		Action:    "新增",
		Student:   student.Name,
		StudentID: student.StudentID,
		Course:    course.CourseName,
		Term:      course.Term,
		OldScore:  0,
		NewScore:  grade.Score,
		Detail:    fmt.Sprintf("新增成绩 %.1f，绩点 %.1f", grade.Score, grade.GradePoint),
	})
}

// updateGradeDirect 直接更新数据库中的成绩记录
func updateGradeDirect(grade *model.Grade) error {
	return config.DB.Save(grade).Error
}
