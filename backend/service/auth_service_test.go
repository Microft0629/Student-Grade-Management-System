package service

import (
	"os"
	"strings"
	"testing"
	"time"

	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/model"
)

func setupAuthTestDB(t *testing.T) {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "sgms_auth_test_*")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll("database", 0755); err != nil {
		t.Fatal(err)
	}
	config.InitDatabase()
	config.CreateDefaultAdmin()
	ClearCurrentUser()

	t.Cleanup(func() {
		if config.DB != nil {
			if sqlDB, err := config.DB.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}
		_ = os.Chdir(oldWD)
		// Windows 上 SQLite 文件句柄释放可能有延迟，重试删除
		for i := 0; i < 5; i++ {
			if err := os.RemoveAll(dir); err == nil {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
}

func TestLoginReturnsUserInfoWithoutPassword(t *testing.T) {
	setupAuthTestDB(t)
	info, err := Login(model.LoginRequest{Username: "admin", Password: "12345678"})
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if info.Username != "admin" || info.Role != "admin" {
		t.Errorf("用户信息错误: %+v", info)
	}
	if !info.MustChangePassword {
		t.Errorf("首次创建的管理员应要求修改密码: %+v", info)
	}
}

func TestLogoutClearsSession(t *testing.T) {
	setupAuthTestDB(t)
	info, err := Login(model.LoginRequest{Username: "admin", Password: "12345678"})
	if err != nil {
		t.Fatal(err)
	}
	SetCurrentUser(info)
	if err := RequireLogin(); err != nil {
		t.Fatalf("登录后 RequireLogin 应通过: %v", err)
	}

	Logout()
	if GetCurrentUser() != nil {
		t.Errorf("退出后会话应被清空")
	}
	if err := RequireLogin(); err == nil || !strings.Contains(err.Error(), "未登录") {
		t.Errorf("退出后 RequireLogin 应返回未登录: %v", err)
	}
}

func TestChangePasswordClearsMustChange(t *testing.T) {
	setupAuthTestDB(t)
	info, err := Login(model.LoginRequest{Username: "admin", Password: "12345678"})
	if err != nil {
		t.Fatal(err)
	}
	SetCurrentUser(info)

	if err := ChangePassword("12345678", "newpass123"); err != nil {
		t.Fatalf("修改密码失败: %v", err)
	}

	var user model.User
	if err := config.DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.MustChangePassword {
		t.Errorf("修改密码后 must_change_password 应清除")
	}

	info2, err := Login(model.LoginRequest{Username: "admin", Password: "newpass123"})
	if err != nil {
		t.Fatalf("新密码登录失败: %v", err)
	}
	if info2.MustChangePassword {
		t.Errorf("登录返回的 must_change_password 应为 false")
	}
	if _, err := Login(model.LoginRequest{Username: "admin", Password: "12345678"}); err == nil {
		t.Errorf("旧密码应已失效")
	}
}

func TestRequireAdmin(t *testing.T) {
	setupAuthTestDB(t)

	admin, err := Login(model.LoginRequest{Username: "admin", Password: "12345678"})
	if err != nil {
		t.Fatal(err)
	}
	SetCurrentUser(admin)
	if err := RequireAdmin(); err != nil {
		t.Errorf("管理员 RequireAdmin 应通过: %v", err)
	}

	if err := CreateTeacher("1000001", "12345678"); err != nil {
		t.Fatal(err)
	}
	teacher, err := Login(model.LoginRequest{Username: "1000001", Password: "12345678"})
	if err != nil {
		t.Fatal(err)
	}
	SetCurrentUser(teacher)
	if err := RequireAdmin(); err == nil || !strings.Contains(err.Error(), "仅管理员") {
		t.Errorf("老师 RequireAdmin 应返回仅管理员: %v", err)
	}

	Logout()
	if err := RequireAdmin(); err == nil || !strings.Contains(err.Error(), "未登录") {
		t.Errorf("未登录 RequireAdmin 应返回未登录: %v", err)
	}
}
