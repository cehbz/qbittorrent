package qbittorrent

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

// TestCookieJarIntegration verifies that the cookie jar is properly integrated
func TestCookieJarIntegration(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/info": {statusCode: http.StatusOK, responseBody: "[]"},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/info"},
	}

	client, _, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify the client has a cookie jar
	if client.client.Jar == nil {
		t.Fatal("Client should have a cookie jar")
	}

	// After login, cookie jar should have the SID cookie
	baseURL, _ := url.Parse(client.baseURL)
	cookies := client.client.Jar.Cookies(baseURL)

	foundSID := false
	for _, cookie := range cookies {
		if cookie.Name == "SID" {
			foundSID = true
			if cookie.Value != "test-session-id" {
				t.Errorf("Expected SID cookie value 'test-session-id', got '%s'", cookie.Value)
			}
		}
	}

	if !foundSID {
		t.Error("SID cookie should be stored in cookie jar after login")
	}

	// Make another request to verify cookies are sent automatically
	ctx := context.Background()
	_, err = client.TorrentsInfoCtx(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestNewClientWithCustomJar verifies that custom cookie jars are respected
func TestNewClientWithCustomJar(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Failed to create cookie jar: %v", err)
	}

	customClient := &http.Client{
		Jar: jar,
		Transport: &mockRoundTripper{
			responses: map[string]mockResponse{
				"/api/v2/auth/login": {statusCode: http.StatusOK, responseBody: "Ok."},
			},
			expectedRequests: []expectedRequest{
				{method: "POST", url: "/api/v2/auth/login"},
			},
			t: &testing.T{},
		},
	}

	client, err := NewClient("user", "pass", "http://localhost:8080", customClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify the custom jar is being used
	if client.client.Jar != jar {
		t.Error("Client should use the provided custom cookie jar")
	}
}

// TestNewClientWithoutJar verifies that a jar is created when client has none
func TestNewClientWithoutJar(t *testing.T) {
	customClient := &http.Client{
		Jar: nil, // No jar provided
		Transport: &mockRoundTripper{
			responses: map[string]mockResponse{
				"/api/v2/auth/login": {statusCode: http.StatusOK, responseBody: "Ok."},
			},
			expectedRequests: []expectedRequest{
				{method: "POST", url: "/api/v2/auth/login"},
			},
			t: &testing.T{},
		},
	}

	client, err := NewClient("user", "pass", "http://localhost:8080", customClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify a jar was automatically created
	if client.client.Jar == nil {
		t.Error("Client should automatically create a cookie jar when none is provided")
	}
}

// TestCookieJarThreadSafety verifies that cookie jar is thread-safe
func TestCookieJarThreadSafety(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/info": {statusCode: http.StatusOK, responseBody: "[]"},
	}

	// Create many expected requests for concurrent access
	var expectedRequests []expectedRequest
	expectedRequests = append(expectedRequests, expectedRequest{method: "POST", url: "/api/v2/auth/login"})
	for i := 0; i < 10; i++ {
		expectedRequests = append(expectedRequests, expectedRequest{method: "GET", url: "/api/v2/torrents/info"})
	}

	client, _, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Make concurrent requests
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.TorrentsInfoCtx(ctx)
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
