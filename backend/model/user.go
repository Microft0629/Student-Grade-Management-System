// Package model user.go 用户（user）模型
package model

type User struct {
	ID       uint   `gorm:"primaryKey"` // 主键
	Username string `gorm:"unique"`     // 唯一账号
	Password string // 加密密码
	Role     string // 权限
}
