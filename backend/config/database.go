// Package config database.go 数据库设置
package config

import (
	"Student-Grade-Management-System/backend/model"
	"Student-Grade-Management-System/backend/utils"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局数据库对象

// InitDatabase 初始化数据库
func InitDatabase() {
	var err error

	// 确保 database 目录存在
	if err := os.MkdirAll(utils.DatabaseDir(), 0755); err != nil {
		log.Fatal("创建数据库目录失败：", err)
	}

	DB, err = gorm.Open( // 建立 GORM 数据库连接
		sqlite.Open(filepath.Join(utils.DatabaseDir(), "student.db")),
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

// CreateDefaultAdmin 创建默认管理员（仅当 admin 账号不存在时创建，绝不重置已有密码）
func CreateDefaultAdmin() {
	var count int64
	DB.Model(&model.User{}).
		Where("username = ?", "admin").
		Count(&count)

	if count > 0 {
		// 已有管理员账号则不做任何操作，避免覆盖用户修改过的密码
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("创建默认管理员密码失败：", err)
	}

	admin := model.User{
		Username:           "admin",
		Password:           string(password),
		Role:               "admin",
		MustChangePassword: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatal("创建默认管理员失败：", err)
	}

	log.Println("默认管理员创建成功")
}
