package qbittorrent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTorrentInfo_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected []string
	}{
		{
			name:     "Empty tags",
			jsonData: `{"tags": ""}`,
			expected: []string{},
		},
		{
			name:     "One tag",
			jsonData: `{"tags": "tag1"}`,
			expected: []string{"tag1"},
		},
		{
			name:     "Multiple tags no spaces",
			jsonData: `{"tags": "tag1,tag2,tag3"}`,
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "Multiple tags with spaces",
			jsonData: `{"tags": "tag1, tag2, tag3"}`,
			expected: []string{"tag1", "tag2", "tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var torrentInfo TorrentInfo
			err := json.Unmarshal([]byte(tt.jsonData), &torrentInfo)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(torrentInfo.Tags) != len(tt.expected) {
				t.Fatalf("expected %d tags, got %d", len(tt.expected), len(torrentInfo.Tags))
			}

			for i, tag := range tt.expected {
				if torrentInfo.Tags[i] != tag {
					t.Errorf("expected tag %v, got %v", tag, torrentInfo.Tags[i])
				}
			}
		})
	}
}

func TestTorrentsAddTags(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":       {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/addTags": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/addTags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsAddTags([]string{"hash1", "hash2"}, "tag1,tag2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsRemoveTags(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":          {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/removeTags": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/removeTags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsRemoveTags([]string{"hash1"}, "tag1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsGetAllTags(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":   {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/tags": {statusCode: http.StatusOK, responseBody: `["tag1","tag2","tag3"]`},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/tags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tags, err := client.TorrentsGetAllTags()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags, got %d", len(tags))
	}
	if tags[0] != "tag1" || tags[1] != "tag2" || tags[2] != "tag3" {
		t.Errorf("Unexpected tags: %v", tags)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsCreateTags(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":          {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/createTags": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/createTags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsCreateTags("newtag1,newtag2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsDeleteTags(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":          {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/deleteTags": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/deleteTags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsDeleteTags("oldtag")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestClient_TorrentsGetTags(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"tags": "tag1, tag2"},{"tags": "tag2, tag3"}]`))
	}))
	defer mockServer.Close()

	client := &Client{
		baseURL: mockServer.URL,
		client:  mockServer.Client(),
	}

	tags, err := client.TorrentsGetTags([]string{"somehash1", "somehash2"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTags := []string{"tag1", "tag2", "tag3"}
	if len(tags) != len(expectedTags) {
		t.Fatalf("expected %d tags, got %d", len(expectedTags), len(tags))
	}

	tagSet := make(map[string]struct{})
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	for _, expectedTag := range expectedTags {
		if _, exists := tagSet[expectedTag]; !exists {
			t.Errorf("expected tag %v to be present", expectedTag)
		}
	}
}
