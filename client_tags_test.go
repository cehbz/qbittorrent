package qbittorrent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestTorrentInfo_UnmarshalJSON_Timestamps(t *testing.T) {
	tests := []struct {
		name             string
		jsonData         string
		expectAddedOn    time.Time
		expectCompletion time.Time
		expectActivity   time.Time
		expectSeen       time.Time
	}{
		{
			name:             "Valid timestamps",
			jsonData:         `{"added_on":1700000000,"completion_on":1700001000,"last_activity":1700002000,"seen_complete":1700003000,"tags":""}`,
			expectAddedOn:    time.Unix(1700000000, 0),
			expectCompletion: time.Unix(1700001000, 0),
			expectActivity:   time.Unix(1700002000, 0),
			expectSeen:       time.Unix(1700003000, 0),
		},
		{
			name:             "Unknown timestamps (-1)",
			jsonData:         `{"added_on":1700000000,"completion_on":-1,"last_activity":1700002000,"seen_complete":-1,"tags":""}`,
			expectAddedOn:    time.Unix(1700000000, 0),
			expectCompletion: time.Time{},
			expectActivity:   time.Unix(1700002000, 0),
			expectSeen:       time.Time{},
		},
		{
			name:             "Zero timestamps",
			jsonData:         `{"added_on":0,"completion_on":0,"last_activity":0,"seen_complete":0,"tags":""}`,
			expectAddedOn:    time.Unix(0, 0),
			expectCompletion: time.Unix(0, 0),
			expectActivity:   time.Unix(0, 0),
			expectSeen:       time.Unix(0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var info TorrentInfo
			if err := json.Unmarshal([]byte(tt.jsonData), &info); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !info.AddedOn.Equal(tt.expectAddedOn) {
				t.Errorf("AddedOn: expected %v, got %v", tt.expectAddedOn, info.AddedOn)
			}
			if !info.CompletionOn.Equal(tt.expectCompletion) {
				t.Errorf("CompletionOn: expected %v, got %v", tt.expectCompletion, info.CompletionOn)
			}
			if !info.LastActivity.Equal(tt.expectActivity) {
				t.Errorf("LastActivity: expected %v, got %v", tt.expectActivity, info.LastActivity)
			}
			if !info.SeenComplete.Equal(tt.expectSeen) {
				t.Errorf("SeenComplete: expected %v, got %v", tt.expectSeen, info.SeenComplete)
			}
		})
	}
}

func TestTorrentInfo_UnmarshalJSON_TagsAndTimestamps(t *testing.T) {
	jsonData := `{"tags":"tag1, tag2","added_on":1700000000,"completion_on":1700001000,"last_activity":1700002000,"seen_complete":-1}`
	var info TorrentInfo
	if err := json.Unmarshal([]byte(jsonData), &info); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(info.Tags) != 2 || info.Tags[0] != "tag1" || info.Tags[1] != "tag2" {
		t.Errorf("expected [tag1 tag2], got %v", info.Tags)
	}
	if !info.AddedOn.Equal(time.Unix(1700000000, 0)) {
		t.Errorf("AddedOn: expected %v, got %v", time.Unix(1700000000, 0), info.AddedOn)
	}
	if !info.SeenComplete.IsZero() {
		t.Errorf("SeenComplete: expected zero, got %v", info.SeenComplete)
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
