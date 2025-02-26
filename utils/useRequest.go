package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RequestInit struct {
	Header map[string]string
	Data   map[string]interface{}
}

type RequestHandler interface {
	Get(string, RequestInit, *interface{}) error
	Post(string, RequestInit, *interface{}) error
}

type RequestHelper struct {
}

func (rh *RequestHelper) Get(url string, init RequestInit, result *interface{}) error {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	
	// Set headers
	for key, value := range init.Header {
		req.Header.Set(key, value)
	}
	
	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// Unmarshal the response data into the result
	return json.Unmarshal(body, result)
}

func (rh *RequestHelper) Post(url string, init RequestInit, result *interface{}) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(init.Data)
	if err != nil {
		return err
	}
	
	// Create a new request with JSON body
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	// Set headers
	for key, value := range init.Header {
		req.Header.Set(key, value)
	}
	
	// Set default Content-Type if not provided
	if _, exists := init.Header["Content-Type"]; !exists {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// Unmarshal the response data into the result
	return json.Unmarshal(body, result)
}