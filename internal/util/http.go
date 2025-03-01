package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const timeout = 5 * time.Second

func SendFact(url string, urlValues map[string]string, token string) ([]byte, error) {
	request, err := createRequest(url, urlValues, token)
	if err != nil {
		return nil, err
	}
	return sendFact(request)
}

func GetFacts(url string, urlValues map[string]string, token string) ([]byte, error) {
	request, err := createRequest(url, urlValues, token)
	if err != nil {
		return nil, err
	}

	return getFacts(request)
}

func createRequest(urlAddr string, urlValues map[string]string, token string) (*http.Request, error) {
	formData := url.Values{}
	for key, val := range urlValues {
		formData.Set(key, val)
	}

	req, err := http.NewRequest(http.MethodPost, urlAddr, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("utils.createRequest is failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	return req, nil
}

func sendFact(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sendFact.Do is failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getFacts.ReadAll is failed: %v", err)
	}

	return responseBody, fmt.Errorf("status code from server is: %d", resp.StatusCode)
}

func getFacts(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getFacts.Do is failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getFacts.ReadAll is failed: %v", err)
	}

	return responseBody, nil
}
