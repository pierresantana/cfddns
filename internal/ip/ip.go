package ip

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultIPURL = "https://api.ipify.org"

func GetPublicIP() (string, error) {
	return getPublicIPFrom(defaultIPURL)
}

func getPublicIPFrom(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetching public IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ipify returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty IP response")
	}

	return ip, nil
}
