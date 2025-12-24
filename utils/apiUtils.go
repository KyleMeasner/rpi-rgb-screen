package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"
)

func SendGetRequest(url string, headers map[string]string, rateLimiter *rate.Limiter) ([]byte, error) {
	if rateLimiter != nil {
		rateLimiter.Wait(context.Background())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func GetAndUnmarshal(url string, headers map[string]string, responseObject any, rateLimiter *rate.Limiter) error {
	responseBody, err := SendGetRequest(url, headers, rateLimiter)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseObject)
}

func SendPostRequest(url string, headers map[string]string, contentType string, payload any, rateLimiter *rate.Limiter) ([]byte, error) {
	if rateLimiter != nil {
		rateLimiter.Wait(context.Background())
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func PostAndUnmarshal(url string, headers map[string]string, contentType string, payload any, responseObject any, rateLimiter *rate.Limiter) error {
	responseBody, err := SendPostRequest(url, headers, contentType, payload, rateLimiter)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseObject)
}

// Returns 0 if conversion fails
func StringToInt(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return result
}
