package qbittorrent

import (
	"net/http"
	"testing"
)

// qBittorrent 5.2 answers 204 No Content to a call that returns no data.
func TestNoContentIsSuccess(t *testing.T) {
	for _, c := range []struct {
		name   string
		status int
		ok     bool
	}{
		{"stop 204", http.StatusNoContent, true},
		{"stop 200", http.StatusOK, true},
		{"stop 404", http.StatusNotFound, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			client, _, err := newMockClient(map[string]mockResponse{
				"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
				"/api/v2/torrents/stop": {statusCode: c.status},
			}, []expectedRequest{
				{method: "POST", url: "/api/v2/auth/login"},
				{method: "POST", url: "/api/v2/torrents/stop"},
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := client.TorrentsPause([]string{"h"}); (err == nil) != c.ok {
				t.Errorf("err = %v, want success %v", err, c.ok)
			}
		})
	}
}

func TestLoginNoContent(t *testing.T) {
	_, _, err := newMockClient(map[string]mockResponse{
		"/api/v2/auth/login": {statusCode: http.StatusNoContent},
	}, []expectedRequest{{method: "POST", url: "/api/v2/auth/login"}})
	if err != nil {
		t.Errorf("login answered 204: %v", err)
	}
}

func TestGetNoContent(t *testing.T) {
	client, _, err := newMockClient(map[string]mockResponse{
		"/api/v2/auth/login":  {statusCode: http.StatusOK, responseBody: "Ok."},
		"/api/v2/app/version": {statusCode: http.StatusNoContent},
	}, []expectedRequest{
		{method: "POST", url: "/api/v2/auth/login"},
		{method: "GET", url: "/api/v2/app/version"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.doGet("/api/v2/app/version", nil); err != nil {
		t.Errorf("GET answered 204: %v", err)
	}
}
