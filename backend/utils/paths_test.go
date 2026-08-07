package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppBaseDirEnvOverride(t *testing.T) {
	t.Setenv(DataDirEnv, `C:\custom\data`)
	if got := AppBaseDir(); got != `C:\custom\data` {
		t.Errorf("AppBaseDir() = %q, want %q", got, `C:\custom\data`)
	}
}

func TestAppBaseDirFallsBackToCWDWithMarkers(t *testing.T) {
	t.Setenv(DataDirEnv, "")
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir, err := os.MkdirTemp("", "sgms_paths_*")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0755); err != nil {
		t.Fatal(err)
	}

	got := AppBaseDir()
	if filepath.Clean(got) != filepath.Clean(dir) {
		t.Errorf("AppBaseDir() = %q, want %q", got, dir)
	}

	_ = os.Chdir(oldWD)
	// Windows 上目录句柄释放可能有延迟，重试删除
	for i := 0; i < 5; i++ {
		if err := os.RemoveAll(dir); err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
