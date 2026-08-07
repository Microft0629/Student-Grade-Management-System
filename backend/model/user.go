// Package model user.go 用户（user）模型
package model

import "time"

type User struct {
	ID                 uint      `gorm:"primaryKey"` // 主键
	Username           string    `gorm:"unique"`     // 唯一账号
	Password           string    // 加密密码
	Role               string    // 权限
	MustChangePassword bool      // 首次登录后是否必须修改密码
	CreatedAt          time.Time // 创建时间
	UpdatedAt          time.Time // 更新时间
}

// UserInfo 返回给前端的用户信息（不含密码哈希等敏感字段）
type UserInfo struct {
	ID                 uint   // 主键
	Username           string // 账号
	Role               string // 权限
	MustChangePassword bool   // 是否必须修改密码
}
