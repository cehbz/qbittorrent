package qbittorrent

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTorrentsProperties(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		wantErr        bool
		expectSavePath string
	}{
		{
			name:           "Success",
			responseBody:   `{"save_path":"/downloads","piece_size":524288,"total_size":1024}`,
			wantErr:        false,
			expectSavePath: "/downloads",
		},
		{
			name:         "Empty response",
			responseBody: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpointResponses := map[string]mockResponse{
				"/api/v2/auth/login": {statusCode: http.StatusOK, responseBody: "Ok."},
				"/api/v2/torrents/properties": {
					statusCode:   http.StatusOK,
					responseBody: tt.responseBody,
				},
			}
			expectedRequests := []expectedRequest{
				{method: "POST", url: "/api/v2/auth/login"},
				{
					method: "GET",
					url:    "/api/v2/torrents/properties",
					query:  url.Values{"hash": []string{"testhash"}},
				},
			}

			client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			props, err := client.TorrentsProperties("testhash")
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if props.SavePath != tt.expectSavePath {
				t.Fatalf("Expected save path %q, got %q", tt.expectSavePath, props.SavePath)
			}

			if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
				t.Errorf("Not all expected requests were made")
			}
		})
	}
}

func TestIntegration_TorrentsProperties(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/torrents/properties" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("hash") != "testhash" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"save_path":"/downloads","total_size":1024}`))
	}))
	defer ts.Close()

	client := &Client{
		client:  ts.Client(),
		baseURL: ts.URL,
	}

	props, err := client.TorrentsProperties("testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if props.SavePath != "/downloads" {
		t.Fatalf("Expected save path '/downloads', got %q", props.SavePath)
	}
	if props.TotalSize != 1024 {
		t.Fatalf("Expected total size 1024, got %d", props.TotalSize)
	}
}
