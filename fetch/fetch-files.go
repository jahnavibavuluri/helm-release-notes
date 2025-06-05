package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func BuildGitHubRawURL(owner, repo, ref, path string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)
}

func DownloadGitHubFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
