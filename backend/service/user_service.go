// Package service user_service.go 用户管理服务
package service

import (
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CreateTeacher 管理员创建老师账号
func CreateTeacher(username string, password string) error {
	if !IsAdmin() {
		return errors.New("仅管理员可创建账号")
	}
	if username == "" {
		return errors.New("用户名不能为空")
	}
	if len(username) != 7 {
		return errors.New("用户名必须为7位数字")
	}
	for _, c := range username {
		if c < '0' || c > '9' {
			return errors.New("用户名只能包含数字")
		}
	}
	if err := validatePassword(password); err != nil {
		return err
	}

	var existing model.User
	result := config.DB.Where("username = ?", username).First(&existing)
	if result.Error == nil {
		return errors.New("账号已存在")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user := model.User{
		Username:           username,
		Password:           string(hashed),
		Role:               "teacher",
		MustChangePassword: true,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("账号 " + username + " 已存在")
		}
		return err
	}
	return nil
}

// DeleteUser 管理员删除老师账号（不能删除 admin）
func DeleteUser(username string) error {
	if !IsAdmin() {
		return errors.New("仅管理员可删除账号")
	}
	if username == "admin" {
		return errors.New("不能删除管理员账号")
	}

	// 名下还有课程时不允许删除，避免产生无法归属的课程与成绩
	var courseCount int64
	if err := config.DB.Model(&model.Course{}).
		Where("creator_name = ?", username).
		Count(&courseCount).Error; err != nil {
		return err
	}
	if courseCount > 0 {
		return errors.New("该老师名下还有课程，请先删除或转移后再删除账号")
	}

	return config.DB.Where("username = ? AND role = ?", username, "teacher").Delete(&model.User{}).Error
}

// GetAllTeachers 获取所有老师账号（不含密码哈希）
func GetAllTeachers() ([]model.UserInfo, error) {
	if !IsAdmin() {
		return nil, errors.New("仅管理员可查看老师账号")
	}

	var users []model.User
	err := config.DB.Where("role = ?", "teacher").Find(&users).Error
	if err != nil {
		return nil, err
	}
	infos := make([]model.UserInfo, 0, len(users))
	for _, u := range users {
		infos = append(infos, model.UserInfo{
			ID:                 u.ID,
			Username:           u.Username,
			Role:               u.Role,
			MustChangePassword: u.MustChangePassword,
		})
	}
	return infos, nil
}

// ChangePassword 修改当前用户密码
func ChangePassword(oldPwd string, newPwd string) error {
	u := GetCurrentUser()
	if u == nil {
		return errors.New("未登录")
	}
	if err := validatePassword(newPwd); err != nil {
		return err
	}

	// 会话中不保存密码哈希，需要时从数据库读取
	var user model.User
	if err := config.DB.First(&user, u.ID).Error; err != nil {
		return errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return errors.New("旧密码错误")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	return config.DB.Model(&model.User{}).Where("id = ?", u.ID).
		Updates(map[string]interface{}{
			"password":             string(hashed),
			"must_change_password": false,
		}).Error
}

// validatePassword 校验密码长度 8-12 位
func validatePassword(pwd string) error {
	if len(pwd) < 8 {
		return errors.New("密码不能少于8位")
	}
	if len(pwd) > 12 {
		return errors.New("密码不能超过12位")
	}
	return nil
}
