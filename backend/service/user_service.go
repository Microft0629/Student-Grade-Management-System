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
		Username: username,
		Password: string(hashed),
		Role:     "teacher",
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
	return config.DB.Where("username = ? AND role = ?", username, "teacher").Delete(&model.User{}).Error
}

// GetAllTeachers 获取所有老师账号
func GetAllTeachers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Where("role = ?", "teacher").Find(&users).Error
	return users, err
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
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPwd)); err != nil {
		return errors.New("旧密码错误")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	return config.DB.Model(&model.User{}).Where("id = ?", u.ID).Update("password", string(hashed)).Error
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
