// Package config database.go 数据库设置
package config

import (
	"Student-Grade-Management-System/backend/model"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局数据库对象

// InitDatabase 初始化数据库
func InitDatabase() {
	var err error

	DB, err = gorm.Open( // 建立 GORM 数据库连接
		sqlite.Open("./database/student.db"),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	log.Println("数据库连接成功")

	err = DB.AutoMigrate( // 创建数据表
		&model.User{},
		&model.Student{},
		&model.Course{},
		&model.Grade{},
		&model.OperationLog{},
	)

	if err != nil {
		log.Fatal("数据表创建失败", err)
	}

	log.Println("数据表创建成功")
}

// CreateDefaultAdmin 创建管理员
func CreateDefaultAdmin() {
	var count int64
	DB.Model(&model.User{}).
		Where("username = ?", "admin").
		Count(&count)

	if count > 0 {
		// 已有管理员则更新密码为 12345678
		password, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
		DB.Model(&model.User{}).Where("username = ?", "admin").Update("password", string(password))
		return
	}

	password, _ := bcrypt.GenerateFromPassword(
		[]byte("12345678"),
		bcrypt.DefaultCost,
	)

	admin := model.User{
		Username: "admin",
		Password: string(password),
		Role:     "admin",
	}

	DB.Create(&admin)

	log.Println("默认管理员创建成功")
}
