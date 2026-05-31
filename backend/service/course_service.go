// Package service course_service.go 课程业务逻辑层
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/repository"
	"errors"
	"strings"
)

// CreateCourse 创建新课程记录，成功后自动同步更新 CSV 备份文件
func CreateCourse(course *model.Course) error {
	if course.CourseCode == "" {
		return errors.New("课程代码不能为空")
	}
	for _, c := range course.CourseCode {
		if c != '.' && (c < '0' || c > '9') {
			return errors.New("课程代码只能包含数字和 .")
		}
	}

	if course.CourseName == "" {
		return errors.New("课程名称不能为空")
	}

	if course.Term == "" {
		return errors.New("请选择学期")
	}

	if course.Credit <= 0 {
		return errors.New("学分必须大于0")
	}
	if course.Teacher == "" {
		return errors.New("任课教师不能为空")
	}
	// 记录创建人
	course.CreatorName = CurrentOperator()

	// 将课程数据持久化至数据库
	err := repository.CreateCourse(course)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("课程代码 " + course.CourseCode + " 已存在，请使用其他代码")
		}
		return err
	}

	// 数据库写入成功后，同步刷新 CSV 文件以保持数据一致性
	return SyncCoursesToCSV()
}

// GetAllCourses 获取所有课程的业务数据
func GetAllCourses() ([]model.Course, error) {
	return repository.GetAllCourses()
}

// DeleteCourse 根据 ID 删除课程及关联成绩，成功后自动同步更新 CSV 备份文件
func DeleteCourse(id uint) error {
	// 老师只能删除自己创建的课程，管理员可删除全部
	if !IsAdmin() {
		var course model.Course
		if err := config.DB.First(&course, id).Error; err != nil {
			return errors.New("课程不存在")
		}
		if course.CreatorName != CurrentOperator() {
			return errors.New("只能删除自己创建的课程")
		}
	}

	// 先删除该课程的所有成绩记录
	config.DB.Where("course_id = ?", id).Delete(&model.Grade{})
	// 再删除课程记录
	err := repository.DeleteCourse(id)
	if err != nil {
		return err
	}

	// 数据库删除成功后，同步刷新 CSV 文件以保持数据一致性
	return SyncCoursesToCSV()
}

// SearchCourses 按课程名称搜索
func SearchCourses(keyword string) ([]model.Course, error) {
	return repository.SearchCoursesByName(keyword)
}
