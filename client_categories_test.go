package qbittorrent

import (
	"net/http"
	"testing"
)

func TestTorrentsCategories(t *testing.T) {
	responseBody := `{"movies":{"name":"movies","savePath":"/downloads/movies"},"tv":{"name":"tv","savePath":"/downloads/tv"}}`
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":          {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/categories": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/categories"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	categories, err := client.TorrentsCategories()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(categories))
	}

	movies, ok := categories["movies"]
	if !ok {
		t.Fatal("Expected 'movies' category")
	}
	if movies.Name != "movies" {
		t.Errorf("Expected name 'movies', got '%s'", movies.Name)
	}
	if movies.SavePath != "/downloads/movies" {
		t.Errorf("Expected savePath '/downloads/movies', got '%s'", movies.SavePath)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsCreateCategory(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":              {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/createCategory": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/createCategory"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsCreateCategory("movies", "/downloads/movies")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsEditCategory(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":            {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/editCategory": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/editCategory"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsEditCategory("movies", "/new/path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsRemoveCategories(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":                {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/removeCategories": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/removeCategories"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsRemoveCategories([]string{"movies", "tv"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
