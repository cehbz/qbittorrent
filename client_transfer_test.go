package qbittorrent

import (
	"net/http"
	"testing"
)

func TestTransferGetInfo(t *testing.T) {
	responseBody := `{"dl_info_speed":1024,"dl_info_data":1048576,"up_info_speed":512,"up_info_data":524288,"dl_rate_limit":0,"up_rate_limit":0,"dht_nodes":100,"connection_status":"connected"}`
	endpointResponses := map[string]mockResponse{
		"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/transfer/info": {statusCode: http.StatusOK, responseBody: responseBody},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/transfer/info"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	info, err := client.TransferGetInfo()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info.DLInfoSpeed != 1024 {
		t.Errorf("Expected DLInfoSpeed 1024, got %d", info.DLInfoSpeed)
	}
	if info.ConnectionStatus != "connected" {
		t.Errorf("Expected ConnectionStatus 'connected', got '%s'", info.ConnectionStatus)
	}
	if info.DHTNodes != 100 {
		t.Errorf("Expected DHTNodes 100, got %d", info.DHTNodes)
	}

	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
