package api_test

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"Student-Grade-Management-System/backend/api"
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/service"
)

// boundAPIs 与 main.go 中绑定到 Wails 的 API 实例保持一致
var boundAPIs = []interface{}{
	api.NewAuthAPI(),
	api.NewStudentAPI(),
	api.NewCourseAPI(),
	api.NewGradeAPI(),
	api.NewStatisticsAPI(),
	api.NewGpaAPI(),
	api.NewBackupAPI(),
	api.NewLogAPI(),
	api.NewExcelAPI(),
	api.NewUserAPI(),
}

// exemptMethods 允许未登录调用的方法（登录/登出/查询角色）
var exemptMethods = map[string]bool{
	"Login":              true,
	"Logout":             true,
	"GetCurrentUserRole": true,
}

// TestAllBoundMethodsRequireLogin 反射枚举所有绑定方法，
// 确保未登录调用时都返回"未登录"，防止新增接口漏加守卫。
func TestAllBoundMethodsRequireLogin(t *testing.T) {
	setupGuardTestDB(t)
	errType := reflect.TypeOf((*error)(nil)).Elem()

	for _, instance := range boundAPIs {
		typ := reflect.TypeOf(instance)
		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			if exemptMethods[m.Name] {
				continue
			}
			lastOut := m.Type.NumOut() - 1
			if lastOut < 0 || !m.Type.Out(lastOut).Implements(errType) {
				continue
			}

			args := make([]reflect.Value, m.Type.NumIn()-1)
			for j := range args {
				args[j] = reflect.Zero(m.Type.In(j + 1))
			}

			t.Run(typ.Elem().Name()+"."+m.Name, func(t *testing.T) {
				service.ClearCurrentUser()
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("未登录调用 %s.%s 发生 panic（缺少守卫）: %v", typ.Elem().Name(), m.Name, r)
					}
				}()

				results := m.Func.Call(append([]reflect.Value{reflect.ValueOf(instance)}, args...))
				errVal := results[lastOut]
				if errVal.IsNil() {
					t.Fatalf("未登录调用 %s.%s 未返回错误（缺少守卫）", typ.Elem().Name(), m.Name)
				}
				err := errVal.Interface().(error)
				if !strings.Contains(err.Error(), "未登录") {
					t.Fatalf("未登录调用 %s.%s 返回 %q，应为未登录", typ.Elem().Name(), m.Name, err.Error())
				}
			})
		}
	}
}

func setupGuardTestDB(t *testing.T) {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "sgms_guard_test_*")
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
	service.ClearCurrentUser()

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
