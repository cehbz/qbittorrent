package qbittorrent

import (
	"encoding/json"
	"net/http"
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

func TestTorrentsProperties_IsPrivate(t *testing.T) {
	jsonData := `{"is_private": true, "addition_date": 1770257484, "completion_date": -1, "creation_date": 1483593698, "last_seen": -1, "name": "test", "popularity": 1.5}`
	var props TorrentsProperties
	err := json.Unmarshal([]byte(jsonData), &props)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !props.IsPrivate {
		t.Error("expected IsPrivate to be true")
	}
	if props.Popularity != 1.5 {
		t.Errorf("expected Popularity 1.5, got %f", props.Popularity)
	}
}

func TestTorrentsPause(t *testing.T) {
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":     {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/pause": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/pause"},
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
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/resume": {statusCode: http.StatusOK, responseBody: "Ok."},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/resume"},
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
