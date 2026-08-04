package service

import "testing"

func TestParseBackupTerm(t *testing.T) {
	cases := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"2025-2026-2_20260803_120000", "2025-2026-2", false},
		{"10871_2025-2026-2_20260803_120000", "2025-2026-2", false},
		{"2024-2025-1_20260530_093425", "2024-2025-1", false},
		{"badname", "", true},
		{"nosep_20260803_120000", "", true},
		{"10871_20260803_120000", "", true},
	}
	for _, c := range cases {
		got, err := parseBackupTerm(c.name)
		if c.wantErr {
			if err == nil {
				t.Errorf("parseBackupTerm(%q) 应报错", c.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseBackupTerm(%q) 返回错误: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseBackupTerm(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}
