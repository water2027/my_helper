package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"fmt"
	"strconv"
)

type RequestInit struct {
	Header map[string]string
	Query  map[string]interface{}
	Body   interface{}
}

type RequestHandler interface {
	Get(string, RequestInit, *interface{}) error
	Post(string, RequestInit, *interface{}) error
}

type RequestHelper struct {
}

func UseRequest() *RequestHelper {
	return &RequestHelper{}
}

// Get performs an HTTP GET request
func (r *RequestHelper) Get(url string, init RequestInit, response interface{}) error {
	return r.doRequest(http.MethodGet, url, init, response)
}

// Post performs an HTTP POST request
func (r *RequestHelper) Post(url string, init RequestInit, response interface{}) error {
	return r.doRequest(http.MethodPost, url, init, response)
}

// doRequest is a helper method to perform the actual HTTP request
func (r *RequestHelper) doRequest(method, url string, init RequestInit, response interface{}) error {
	var bodyReader io.Reader

	// Process body if provided (mainly for POST)
	if init.Body != nil && method == http.MethodPost {
		jsonBody, err := json.Marshal(init.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range init.Header {
		req.Header.Set(key, value)
	}

	// Set content type if not set and body exists for POST
	if init.Body != nil && method == http.MethodPost && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set query parameters
	if init.Query != nil {
		q := req.URL.Query()
		for key, value := range init.Query {
			switch v := value.(type) {
			case string:
				q.Add(key, v)
			case int:
				q.Add(key, strconv.Itoa(v))
			case float64:
				q.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
			case bool:
				q.Add(key, strconv.FormatBool(v))
			case []string:
				for _, item := range v {
					q.Add(key, item)
				}
			default:
				// For any other type, convert to string using fmt.Sprint
				q.Add(key, fmt.Sprint(v))
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s, Body: %s", resp.StatusCode, resp.Status, string(body))
	}

	// If response pointer is nil or body is empty, don't try to unmarshal
	if response == nil || len(body) == 0 {
		return nil
	}

	// Parse JSON response
	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}