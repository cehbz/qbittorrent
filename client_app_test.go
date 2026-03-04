package qbittorrent

import (
	"net/http"
	"testing"
)

func TestAppVersion(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":  {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/app/version": {statusCode: http.StatusOK, responseBody: "v4.6.0"},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/app/version"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	version, err := client.AppVersion()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if version != "v4.6.0" {
		t.Errorf("Expected version 'v4.6.0', got '%s'", version)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestAppDefaultSavePath(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":         {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/app/defaultSavePath": {statusCode: http.StatusOK, responseBody: "/downloads"},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/app/defaultSavePath"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	path, err := client.AppDefaultSavePath()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if path != "/downloads" {
		t.Errorf("Expected path '/downloads', got '%s'", path)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestAppPreferences(t *testing.T) {
	responseBody := `{"save_path":"/downloads","max_connec":500}`
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/app/preferences": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/app/preferences"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	prefs, err := client.AppPreferences()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if prefs["save_path"] != "/downloads" {
		t.Errorf("Expected save_path '/downloads', got '%v'", prefs["save_path"])
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestSetAppPreferences(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":         {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/app/setPreferences": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/app/setPreferences"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	prefs := map[string]any{"save_path": "/new/downloads"}
	err = client.SetAppPreferences(prefs)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestAuthLogout(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":  {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/auth/logout": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/auth/logout"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.AuthLogout()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
