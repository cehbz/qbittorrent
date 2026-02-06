package qbittorrent

import (
	"context"
	"net/http"
	"testing"
	"time"
)

// TestContextCancellation tests that context cancellation is properly handled
func TestContextCancellation(t *testing.T) {
	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Mock responses
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/info": {statusCode: http.StatusOK, responseBody: "[]"},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
	}

	client, _, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// This should fail because context is cancelled
	_, err = client.TorrentsInfoCtx(ctx)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}

// TestContextTimeout tests that context timeout is properly handled
func TestContextTimeout(t *testing.T) {
	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give context time to expire
	time.Sleep(2 * time.Millisecond)

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/info": {statusCode: http.StatusOK, responseBody: "[]"},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
	}

	client, _, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.TorrentsInfoCtx(ctx)
	if err == nil {
		t.Error("Expected error due to context timeout, got nil")
	}
}

// TestTorrentsInfoCtx tests the context-aware version of TorrentsInfo
func TestTorrentsInfoCtx(t *testing.T) {
	ctx := context.Background()
	responseBody := `[{"hash":"abc123","name":"Test Torrent","progress":0.5}]`

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/info": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/info"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	torrents, err := client.TorrentsInfoCtx(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(torrents) != 1 {
		t.Errorf("Expected 1 torrent, got %d", len(torrents))
	}

	if torrents[0].Name != "Test Torrent" {
		t.Errorf("Expected 'Test Torrent', got '%s'", torrents[0].Name)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsExportCtx tests the context-aware version of TorrentsExport
func TestTorrentsExportCtx(t *testing.T) {
	ctx := context.Background()
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
		t.Fatalf("Failed to create client: %v", err)
	}

	data, err := client.TorrentsExportCtx(ctx, "testhash")
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

// TestTorrentsAddCtx tests the context-aware version of TorrentsAdd
func TestTorrentsAddCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsAddCtx(ctx, "test.torrent", []byte("torrent data"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsDeleteCtx tests the context-aware version of TorrentsDelete
func TestTorrentsDeleteCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsDeleteCtx(ctx, "testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestSetForceStartCtx tests the context-aware version of SetForceStart
func TestSetForceStartCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetForceStartCtx(ctx, "testhash", true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsTrackersCtx tests the context-aware version of TorrentsTrackers
func TestTorrentsTrackersCtx(t *testing.T) {
	ctx := context.Background()
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
		t.Fatalf("Failed to create client: %v", err)
	}

	trackers, err := client.TorrentsTrackersCtx(ctx, "testhash")
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

// TestTorrentsPropertiesCtx tests the context-aware version of TorrentsProperties
func TestTorrentsPropertiesCtx(t *testing.T) {
	ctx := context.Background()
	responseBody := `{
		"hash": "testhash",
		"name": "Test Torrent",
		"save_path": "/downloads",
		"total_size": 1024,
		"addition_date": 1640000000,
		"creation_date": 1639000000,
		"completion_date": 1641000000,
		"last_seen": 1642000000
	}`

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":           {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/properties": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/properties"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	props, err := client.TorrentsPropertiesCtx(ctx, "testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if props.Name != "Test Torrent" {
		t.Errorf("Expected 'Test Torrent', got '%s'", props.Name)
	}

	if props.TotalSize != 1024 {
		t.Errorf("Expected 1024, got %d", props.TotalSize)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsAddTagsCtx tests the context-aware version of TorrentsAddTags
func TestTorrentsAddTagsCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsAddTagsCtx(ctx, "testhash", "tag1,tag2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsRemoveTagsCtx tests the context-aware version of TorrentsRemoveTags
func TestTorrentsRemoveTagsCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsRemoveTagsCtx(ctx, "testhash", "tag1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsCreateTagsCtx tests the context-aware version of TorrentsCreateTags
func TestTorrentsCreateTagsCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsCreateTagsCtx(ctx, "tag1,tag2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsDeleteTagsCtx tests the context-aware version of TorrentsDeleteTags
func TestTorrentsDeleteTagsCtx(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TorrentsDeleteTagsCtx(ctx, "tag1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestTorrentsGetAllTagsCtx tests the context-aware version of TorrentsGetAllTags
func TestTorrentsGetAllTagsCtx(t *testing.T) {
	ctx := context.Background()
	responseBody := `["tag1","tag2","tag3"]`

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/tags": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/torrents/tags"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tags, err := client.TorrentsGetAllTagsCtx(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestSyncMainDataCtx tests the context-aware version of SyncMainData
func TestSyncMainDataCtx(t *testing.T) {
	ctx := context.Background()
	responseBody := `{"rid":1,"full_update":true,"torrents":{},"server_state":{}}`

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":     {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/sync/maindata": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/sync/maindata"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	mainData, err := client.SyncMainDataCtx(ctx, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mainData.Rid != 1 {
		t.Errorf("Expected rid 1, got %d", mainData.Rid)
	}

	if !mainData.FullUpdate {
		t.Error("Expected full_update to be true")
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestSyncTorrentPeersCtx tests the context-aware version of SyncTorrentPeers
func TestSyncTorrentPeersCtx(t *testing.T) {
	ctx := context.Background()
	responseBody := `{"rid":1,"full_update":true,"peers":{},"show_flags":true}`

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":         {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/sync/torrentPeers": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/sync/torrentPeers"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	peers, err := client.SyncTorrentPeersCtx(ctx, "testhash", 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if peers.Rid != 1 {
		t.Errorf("Expected rid 1, got %d", peers.Rid)
	}

	if !peers.ShowFlags {
		t.Error("Expected show_flags to be true")
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

// TestBackwardCompatibility ensures old methods still work
func TestBackwardCompatibility(t *testing.T) {
	expectedData := "torrent file data"

	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":      {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/torrents/export": {statusCode: http.StatusOK, responseBody: expectedData},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "POST", url: "/api/v2/torrents/export"},
	}

	client, _, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test that the old non-Ctx method still works
	data, err := client.TorrentsExport("testhash")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(data) != expectedData {
		t.Errorf("Expected %s, got %s", expectedData, string(data))
	}
}
