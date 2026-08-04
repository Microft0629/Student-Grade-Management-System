package utils

import (
	"testing"

	"Student-Grade-Management-System/backend/model"
)

func TestCalculateGradePoint(t *testing.T) {
	cases := []struct {
		score float64
		want  float64
	}{
		{0, 0},
		{59, 0},
		{60, 1},
		{65, 1.5},
		{74, 2.4},
		{80, 3},
		{85, 3.5},
		{95, 4.5},
		{100, 5},
	}
	for _, c := range cases {
		if got := CalculateGradePoint(c.score); got != c.want {
			t.Errorf("CalculateGradePoint(%v) = %v, want %v", c.score, got, c.want)
		}
	}
}

func TestCalculateStudentGPA(t *testing.T) {
	grades := []model.Grade{
		{Course: model.Course{Credit: 4}, GradePoint: 4},
		{Course: model.Course{Credit: 2}, GradePoint: 2},
	}
	// (4*4 + 2*2) / (4+2) = 20/6 ≈ 3.3333 → 3.33
	if got := CalculateStudentGPA(grades); got != 3.33 {
		t.Errorf("CalculateStudentGPA = %v, want 3.33", got)
	}

	if got := CalculateStudentGPA(nil); got != 0 {
		t.Errorf("CalculateStudentGPA(nil) = %v, want 0", got)
	}
}
