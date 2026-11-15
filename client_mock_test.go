package qbittorrent

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

// recorder is a testing.TB that keeps the errors reported to it.
type recorder struct {
	testing.TB
	errors []string
}

func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

// The fake server fails the test when a request differs from the one expected.
func TestMockReportsMismatch(t *testing.T) {
	for _, c := range []struct {
		name string
		want expectedRequest
	}{
		{"method", expectedRequest{method: "GET", url: "/api/v2/torrents/stop"}},
		{"path", expectedRequest{method: "POST", url: "/api/v2/torrents/pause"}},
		{"form", expectedRequest{method: "POST", url: "/api/v2/torrents/stop", params: url.Values{"hashes": {"other"}}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			rec := &recorder{TB: t}
			client, _, err := newMockClient(rec, map[string]mockResponse{
				"/api/v2/auth/login":    {statusCode: http.StatusOK, responseBody: "Ok."},
				"/api/v2/torrents/stop": {statusCode: http.StatusOK},
			}, []expectedRequest{{method: "POST", url: "/api/v2/auth/login"}, c.want})
			if err != nil {
				t.Fatal(err)
			}
			client.TorrentsPause([]string{"h"})
			if len(rec.errors) == 0 {
				t.Errorf("a %s mismatch went unreported", c.name)
			}
		})
	}
}
