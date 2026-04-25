package mcp

import (
	"testing"

	"github.com/HidemaruOwO/pummit/internal/usecase"
)

func TestFormatDoctorReport(t *testing.T) {
	report := usecase.DoctorReport{
		System: []usecase.DoctorCheck{{Name: "GOOS", Status: "OK", Message: "windows"}},
		Checks: []usecase.DoctorCheck{{Name: "Config", Status: "OK", Message: "Configuration is valid"}},
	}
	got := formatDoctorReport(report)
	want := "GOOS: windows\nConfig [OK]: Configuration is valid"
	if got != want {
		t.Fatalf("formatDoctorReport() = %q, want %q", got, want)
	}
}
