package qbittorrent

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestTorrentsAddParams(t *testing.T) {
	torrentData := []byte("test torrent data")

	tests := []struct {
		name           string
		params         *TorrentsAddParams
		checkRequest   func(*testing.T, *http.Request)
		responseStatus int
		responseBody   string
		wantErr        bool
	}{
		{
			name: "Basic torrent with category",
			params: &TorrentsAddParams{
				Torrents:  [][]byte{torrentData},
				SavePath:  "/downloads",
				Category:  "movies",
				SkipCheck: true,
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				// Parse multipart form
				err := req.ParseMultipartForm(32 << 20)
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				// Check savepath
				if savepath := req.FormValue("savepath"); savepath != "/downloads" {
					t.Errorf("Expected savepath '/downloads', got '%s'", savepath)
				}

				// Check category
				if category := req.FormValue("category"); category != "movies" {
					t.Errorf("Expected category 'movies', got '%s'", category)
				}

				// Check skip_checking
				if skip := req.FormValue("skip_checking"); skip != "true" {
					t.Errorf("Expected skip_checking 'true', got '%s'", skip)
				}

				// Check torrent file
				if req.MultipartForm == nil || len(req.MultipartForm.File["torrents"]) == 0 {
					t.Error("No torrent file found in multipart form")
				}
			},
			responseStatus: http.StatusOK,
			responseBody:   "Ok.",
			wantErr:        false,
		},
		{
			name: "Torrent with all parameters",
			params: &TorrentsAddParams{
				Torrents:    [][]byte{torrentData},
				SavePath:    "/downloads",
				Category:    "tv",
				Tags:        "hd,x265",
				SkipCheck:   true,
				Paused:      true,
				RootFolder:  boolPtr(false),
				Rename:      "renamed_torrent",
				UpLimit:     1048576, // 1MB/s
				DlLimit:     2097152, // 2MB/s
				RatioLimit:  2.5,
				SeedingTime: 1440, // 24 hours
				AutoTMM:     true,
				Sequential:  true,
				FirstLast:   true,
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				err := req.ParseMultipartForm(32 << 20)
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				// Check all parameters
				checks := map[string]string{
					"savepath":           "/downloads",
					"category":           "tv",
					"tags":               "hd,x265",
					"skip_checking":      "true",
					"paused":             "true",
					"root_folder":        "false",
					"rename":             "renamed_torrent",
					"upLimit":            "1048576",
					"dlLimit":            "2097152",
					"ratioLimit":         "2.50",
					"seedingTimeLimit":   "1440",
					"autoTMM":            "true",
					"sequentialDownload": "true",
					"firstLastPiecePrio": "true",
				}

				for field, expected := range checks {
					if value := req.FormValue(field); value != expected {
						t.Errorf("Expected %s '%s', got '%s'", field, expected, value)
					}
				}
			},
			responseStatus: http.StatusOK,
			responseBody:   "Ok.",
			wantErr:        false,
		},
		{
			name: "Multiple torrents",
			params: &TorrentsAddParams{
				Torrents: [][]byte{
					[]byte("torrent1"),
					[]byte("torrent2"),
					[]byte("torrent3"),
				},
				Category: "movies",
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				err := req.ParseMultipartForm(32 << 20)
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				// Check we have 3 torrent files
				if req.MultipartForm == nil {
					t.Fatal("No multipart form found")
				}

				files := req.MultipartForm.File["torrents"]
				if len(files) != 3 {
					t.Errorf("Expected 3 torrent files, got %d", len(files))
				}

				// Verify file contents
				for i, file := range files {
					f, err := file.Open()
					if err != nil {
						t.Errorf("Failed to open file %d: %v", i, err)
						continue
					}
					defer f.Close()

					content, err := io.ReadAll(f)
					if err != nil {
						t.Errorf("Failed to read file %d: %v", i, err)
						continue
					}

					expected := []byte("torrent" + string(rune('1'+i)))
					if !bytes.Equal(content, expected) {
						t.Errorf("File %d content mismatch: got %s, want %s", i, content, expected)
					}
				}
			},
			responseStatus: http.StatusOK,
			responseBody:   "Ok.",
			wantErr:        false,
		},
		{
			name: "URLs instead of torrents",
			params: &TorrentsAddParams{
				URLs: []string{
					"magnet:?xt=urn:btih:123",
					"http://example.com/torrent.torrent",
				},
				Category: "music",
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				err := req.ParseMultipartForm(32 << 20)
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				// Check URLs
				urls := req.Form["urls"]
				if len(urls) != 2 {
					t.Errorf("Expected 2 URLs, got %d", len(urls))
				}

				expectedURLs := []string{
					"magnet:?xt=urn:btih:123",
					"http://example.com/torrent.torrent",
				}

				for i, url := range urls {
					if i < len(expectedURLs) && url != expectedURLs[i] {
						t.Errorf("URL %d mismatch: got %s, want %s", i, url, expectedURLs[i])
					}
				}
			},
			responseStatus: http.StatusOK,
			responseBody:   "Ok.",
			wantErr:        false,
		},
		{
			name: "Mixed torrents and URLs",
			params: &TorrentsAddParams{
				Torrents: [][]byte{torrentData},
				URLs:     []string{"magnet:?xt=urn:btih:456"},
				Category: "books",
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				err := req.ParseMultipartForm(32 << 20)
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				// Check torrent file
				if req.MultipartForm == nil || len(req.MultipartForm.File["torrents"]) != 1 {
					t.Error("Expected 1 torrent file")
				}

				// Check URL
				urls := req.Form["urls"]
				if len(urls) != 1 || urls[0] != "magnet:?xt=urn:btih:456" {
					t.Error("Expected 1 URL with correct value")
				}
			},
			responseStatus: http.StatusOK,
			responseBody:   "Ok.",
			wantErr:        false,
		},
		{
			name: "Error response",
			params: &TorrentsAddParams{
				Torrents: [][]byte{torrentData},
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				// Request is valid, but server returns error
			},
			responseStatus: http.StatusForbidden,
			responseBody:   "Forbidden",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock successful AuthLogin and custom response for TorrentsAdd
			endpointResponses := map[string]mockResponse{
				"/api/v2/auth/login":   {statusCode: http.StatusOK, responseBody: "Ok."},
				"/api/v2/torrents/add": {statusCode: tt.responseStatus, responseBody: tt.responseBody},
			}

			expectedRequests := []expectedRequest{
				{method: "POST", url: "/api/v2/auth/login"},
				{method: "POST", url: "/api/v2/torrents/add"},
			}

			// Create custom handler for request inspection
			customHandler := map[string]func(*http.Request){
				"/api/v2/torrents/add": func(req *http.Request) {
					tt.checkRequest(t, req)
				},
			}

			client, mockTransport, err := newMockClientWithHandler(endpointResponses, expectedRequests, customHandler)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.TorrentsAddParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TorrentsAddParams() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check the request made
			if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
				t.Errorf("Not all expected requests were made")
			}
		})
	}
}

func TestTorrentsAddParams_BackwardCompatibility(t *testing.T) {
	// Test that the original TorrentsAdd method still works
	torrentData := []byte("test torrent data")

	// Mock successful AuthLogin and TorrentsAdd responses
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":   {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/add": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/add"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsAdd("test.torrent", torrentData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}

func TestTorrentsAddParams_EmptyParams(t *testing.T) {
	// Test with empty params struct
	params := &TorrentsAddParams{}

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":   {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/add": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/add"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should succeed even with no torrents or URLs
	err = client.TorrentsAddParams(params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
