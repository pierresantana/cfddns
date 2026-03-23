package cloudflare

import (
	"encoding/json"
	"net/http/httptest"
)

func newTestClient(server *httptest.Server) *Client {
	client := NewClient("test-token")
	client.baseURL = server.URL
	return client
}

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
