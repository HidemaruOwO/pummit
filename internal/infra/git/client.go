package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	Dir string
}

func NewClient(dir string) *Client {
	return &Client{Dir: dir}
}

func (c *Client) StagedFiles() ([]string, error) {
	output, err := c.run("diff", "--name-only", "--cached")
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return []string{}, nil
	}

	return strings.Split(trimmed, "\n"), nil
}

func (c *Client) CurrentBranch() (string, error) {
	output, err := c.run("branch", "--show-current")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func (c *Client) Commit(message string) error {
	_, err := c.run("commit", "-m", message)
	return err
}

func (c *Client) AddFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"add"}, files...)
	_, err := c.run(args...)
	return err
}

func (c *Client) ChangedFiles() ([]string, error) {
	output, err := c.run("status", "--porcelain=v1", "-u")
	if err != nil {
		return nil, err
	}

	raw := bytes.TrimRight([]byte(output), "\n\r")
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	files := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
		}
		path := line[3:]
		if strings.Contains(path, " -> ") {
			parts := strings.SplitN(path, " -> ", 2)
			if len(parts) == 2 {
				path = parts[1]
			}
		}
		path = strings.TrimSpace(path)
		if path != "" {
			files = append(files, path)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if trimmed == "" {
			return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
		}
		return "", fmt.Errorf("git %s failed: %w: %s", strings.Join(args, " "), err, trimmed)
	}

	return string(output), nil
}
