package usecase

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
)

type DoctorCheck struct {
	Name    string
	Status  string
	Message string
}

type DoctorReport struct {
	System []DoctorCheck
	Checks []DoctorCheck
}

type DoctorService struct {
	config *ConfigService
	store  *rootconfig.Store
}

func NewDoctorService(config *ConfigService, store *rootconfig.Store) *DoctorService {
	return &DoctorService{config: config, store: store}
}

func (s *DoctorService) Run() DoctorReport {
	report := DoctorReport{
		System: []DoctorCheck{
			{Name: "GOOS", Status: "OK", Message: runtime.GOOS},
			{Name: "GOARCH", Status: "OK", Message: runtime.GOARCH},
		},
	}

	configDir, err := s.store.ConfigDir()
	if err != nil {
		report.Checks = append(report.Checks, DoctorCheck{Name: "Config", Status: "ERROR", Message: err.Error()})
	} else {
		report.System = append(report.System, DoctorCheck{Name: "ConfigDir", Status: "OK", Message: configDir})
		if _, err := s.config.Validate(); err != nil {
			report.Checks = append(report.Checks, DoctorCheck{Name: "Config", Status: "ERROR", Message: err.Error()})
		} else {
			report.Checks = append(report.Checks, DoctorCheck{Name: "Config", Status: "OK", Message: "Configuration is valid"})
		}
	}

	report.Checks = append(report.Checks, s.gitCheck())
	report.Checks = append(report.Checks, s.repoCheck())
	report.Checks = append(report.Checks, s.networkCheck())

	return report
}

func (s *DoctorService) gitCheck() DoctorCheck {
	cmd := exec.Command("git", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return DoctorCheck{Name: "Git", Status: "ERROR", Message: err.Error()}
	}
	return DoctorCheck{Name: "Git", Status: "OK", Message: string(output)}
}

func (s *DoctorService) repoCheck() DoctorCheck {
	wd, err := os.Getwd()
	if err != nil {
		return DoctorCheck{Name: "Repository", Status: "ERROR", Message: err.Error()}
	}
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = wd
	if err := cmd.Run(); err != nil {
		return DoctorCheck{Name: "Repository", Status: "WARNING", Message: "Not inside a git repository"}
	}
	return DoctorCheck{Name: "Repository", Status: "OK", Message: "Inside a git repository"}
}

func (s *DoctorService) networkCheck() DoctorCheck {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://api.github.com", nil)
	if err != nil {
		return DoctorCheck{Name: "Network", Status: "ERROR", Message: err.Error()}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DoctorCheck{Name: "Network", Status: "WARNING", Message: "GitHub is unreachable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return DoctorCheck{Name: "Network", Status: "OK", Message: "GitHub is reachable"}
	}
	return DoctorCheck{Name: "Network", Status: "WARNING", Message: resp.Status}
}
