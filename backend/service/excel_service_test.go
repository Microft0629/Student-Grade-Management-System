package service

import "testing"

func TestSanitizeFileName(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"高等数学（上）", "高等数学（上）"},
		{`高/等:数*学?<上>|`, "高_等_数_学__上__"},
		{`  "测试"  `, "_测试_"},
		{"  ", "未命名"},
		{"...", "未命名"},
		{"CON", "_CON"},
		{"con.txt", "con.txt"},
		{"LPT1", "_LPT1"},
	}
	for _, c := range cases {
		if got := sanitizeFileName(c.name); got != c.want {
			t.Errorf("sanitizeFileName(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}
