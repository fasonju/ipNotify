package requests

import (
	"io"
	"net/http"
	"strings"
)

func GetIP(url string) (string, error) {
	// Create a new HTTP client (using default settings)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(body), " \t\n\r"), nil
}
