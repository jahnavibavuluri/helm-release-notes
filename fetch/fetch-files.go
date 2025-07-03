package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseURL string
	Token   string
}

func (c *Config) BuildGitHubRawURL(owner, repo, ref, path string) string {
	if strings.Contains(c.BaseURL, "github.com") {
		// Public GitHub
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)
	} else {
		// Private GitHub (for MW)
		return fmt.Sprintf("%s/api/v3/repos/%s/%s/contents/%s?ref=%s", c.BaseURL, owner, repo, path, ref)
	}
}

func (c *Config) DownloadGitHubFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %v", url, err)
	}

	// Add authorization header if token is provided
	if c.Token != "" {
		req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	}

	// Headers needed by GH
	req.Header.Add("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("User-Agent", "tf-diff/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s returned status code %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return data, nil
}

func CreateTempFile(content []byte, filename string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "terraform-diff-*")
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(tmpDir, filename)
	err = os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		return "", err
	}

	return tmpFile, nil
}

// NewConfig creates a new GitHub configuration
func NewConfig(baseURL, token string) *Config {
	return &Config{
		BaseURL: baseURL,
		Token:   token,
	}
}

func NewPublicGitHubConfig(token string) *Config {
	return NewConfig("https://github.com", token)
}

func NewPublicGitHubConfigNoAuth() *Config {
	return NewConfig("https://github.com", "")
}

func NewEnterpriseGitHubConfig(baseURL, token string) *Config {
	return NewConfig(baseURL, token)
}
