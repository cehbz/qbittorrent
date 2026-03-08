package qbittorrent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestTorrentsExport(t *testing.T) {
	expectedData := "torrent file data"
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/export": {statusCode: http.StatusOK, responseBody: expectedData},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/export"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := client.TorrentsExport("testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(data) != expectedData {
		t.Errorf("Expected %s, got %s", expectedData, string(data))
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsAdd(t *testing.T) {
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

	err = client.TorrentsAdd([]byte("torrent data"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsDelete(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/delete": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/delete"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsDelete([]string{"testhash"}, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestSetForceStart(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":             {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setForceStart": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setForceStart"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.SetForceStart([]string{"testhash"}, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsTrackers(t *testing.T) {
	responseBody := `[{"url":"tracker1","status":1},{"url":"tracker2","status":0}]`
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":        {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/trackers": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/trackers"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	trackers, err := client.TorrentsTrackers("testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(trackers) != 2 {
		t.Errorf("Expected 2 trackers, got %d", len(trackers))
	}

	if trackers[0].URL != "tracker1" {
		t.Errorf("Expected tracker URL 'tracker1', got '%s'", trackers[0].URL)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsProperties_Private(t *testing.T) {
	t.Run("private true", func(t *testing.T) {
		jsonData := `{"private": true, "addition_date": 1770257484, "completion_date": -1, "creation_date": 1483593698, "last_seen": -1, "name": "test", "popularity": 1.5}`
		var props TorrentsProperties
		err := json.Unmarshal([]byte(jsonData), &props)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if props.Private == nil || !*props.Private {
			t.Error("expected Private to be true")
		}
		if props.Popularity != 1.5 {
			t.Errorf("expected Popularity 1.5, got %f", props.Popularity)
		}
	})

	t.Run("private null (no metadata)", func(t *testing.T) {
		jsonData := `{"private": null, "addition_date": 0, "completion_date": -1, "creation_date": 0, "last_seen": -1}`
		var props TorrentsProperties
		err := json.Unmarshal([]byte(jsonData), &props)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if props.Private != nil {
			t.Errorf("expected Private to be nil, got %v", *props.Private)
		}
	})
}

func TestCookieJar_SIDPersisted(t *testing.T) {
	var requestCount atomic.Int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		switch {
		case r.URL.Path == "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "SID", Value: "test-session-id", Path: "/"})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))

		case r.URL.Path == "/api/v2/app/version":
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value != "test-session-id" {
				t.Errorf("request %d: expected SID cookie 'test-session-id', got err=%v cookie=%v", n, err, cookie)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v5.0.3"))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client, err := NewClient("user", "pass", ts.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	version, err := client.AppVersion()
	if err != nil {
		t.Fatalf("AppVersion error: %v", err)
	}
	if version != "v5.0.3" {
		t.Errorf("expected 'v5.0.3', got %q", version)
	}
}

func TestCookieJar_ReauthOn403(t *testing.T) {
	var loginCount atomic.Int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			n := loginCount.Add(1)
			http.SetCookie(w, &http.Cookie{
				Name:  "SID",
				Value: "session-" + string(rune('0'+n)),
				Path:  "/",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))

		case "/api/v2/app/version":
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value == "session-1" {
				// First session is "expired"
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v5.0.3"))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client, err := NewClient("user", "pass", ts.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// First real request gets 403, triggers re-auth, retries with new session
	version, err := client.AppVersion()
	if err != nil {
		t.Fatalf("AppVersion error: %v", err)
	}
	if version != "v5.0.3" {
		t.Errorf("expected 'v5.0.3', got %q", version)
	}
	if loginCount.Load() != 2 {
		t.Errorf("expected 2 logins (initial + re-auth), got %d", loginCount.Load())
	}
}

func TestCookieJar_LogoutClearsCookies(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "SID", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))

		case "/api/v2/auth/logout":
			w.WriteHeader(http.StatusOK)

		case "/api/v2/app/version":
			_, err := r.Cookie("SID")
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v5.0.3"))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client, err := NewClient("user", "pass", ts.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.AuthLogout()
	if err != nil {
		t.Fatalf("AuthLogout error: %v", err)
	}

	// After logout, cookie jar should be empty; next request should get 403
	// which triggers re-auth (login again), so it should still succeed
	version, err := client.AppVersion()
	if err != nil {
		t.Fatalf("AppVersion after logout error: %v", err)
	}
	if version != "v5.0.3" {
		t.Errorf("expected 'v5.0.3', got %q", version)
	}
}

func TestTorrentsPause(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/stop": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/stop"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsPause([]string{"hash1", "hash2"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsResume(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":     {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/start": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/start"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsResume([]string{"hash1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsRecheck(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":       {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/recheck": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/recheck"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsRecheck([]string{"hash1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsRename(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/rename": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/rename"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsRename("hash1", "new name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsFiles(t *testing.T) {
	responseBody := `[{"index":0,"name":"file1.txt","size":1024,"progress":1.0,"priority":1,"is_seed":false,"piece_range":[0,10],"availability":1.0}]`
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":     {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/files": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/files"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files, err := client.TorrentsFiles("hash1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}
	if files[0].Name != "file1.txt" {
		t.Errorf("Expected name 'file1.txt', got '%s'", files[0].Name)
	}
	if files[0].Size != 1024 {
		t.Errorf("Expected size 1024, got %d", files[0].Size)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetLocation(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":           {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setLocation": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setLocation"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetLocation([]string{"hash1"}, "/new/path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsReannounce(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":           {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/reannounce":  {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/reannounce"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsReannounce([]string{"hash1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetCategory(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":            {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setCategory":  {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setCategory"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetCategory([]string{"hash1", "hash2"}, "movies")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetAutoTMM(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":                 {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setAutoManagement": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setAutoManagement"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetAutoTMM([]string{"hash1"}, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetDownloadLimit(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":                 {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setDownloadLimit":  {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setDownloadLimit"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetDownloadLimit([]string{"hash1"}, 1048576)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetUploadLimit(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":               {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setUploadLimit":  {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setUploadLimit"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetUploadLimit([]string{"hash1"}, 524288)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsSetShareLimits(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":               {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/setShareLimits":  {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/setShareLimits"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsSetShareLimits([]string{"hash1"}, 2.0, 1440, 720)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTorrentsFilePrio(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":        {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/filePrio": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/filePrio"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TorrentsFilePrio("hash1", []int{0, 1, 3}, 7)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
