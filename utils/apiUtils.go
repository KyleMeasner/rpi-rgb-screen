package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

func SendGetRequest(url string, rateLimiter *rate.Limiter) ([]byte, error) {
	if rateLimiter != nil {
		rateLimiter.Wait(context.Background())
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func GetAndUnmarshal(url string, responseObject any, rateLimiter *rate.Limiter) error {
	responseBody, err := SendGetRequest(url, rateLimiter)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseObject)
}

func SendPostRequest(url string, contentType string, payload any, rateLimiter *rate.Limiter) ([]byte, error) {
	if rateLimiter != nil {
		rateLimiter.Wait(context.Background())
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, contentType, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func PostAndUnmarshal(url string, contentType string, payload any, responseObject any, rateLimiter *rate.Limiter) error {
	responseBody, err := SendPostRequest(url, contentType, payload, rateLimiter)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseObject)
}
