package cloudflare

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetZoneID(t *testing.T) {
	t.Run("returns zone ID when found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("name"); got != "sstec.com.br" {
				t.Errorf("expected name=sstec.com.br, got %s", got)
			}
			zones := []zone{{ID: "zone-123", Name: "sstec.com.br"}}
			resp := apiResponse{Success: true, Result: mustMarshal(zones)}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := newTestClient(server)
		id, err := client.GetZoneID("sstec.com.br")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "zone-123" {
			t.Errorf("got %q, want %q", id, "zone-123")
		}
	})

	t.Run("returns error when zone not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := apiResponse{Success: true, Result: mustMarshal([]zone{})}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := newTestClient(server)
		_, err := client.GetZoneID("unknown.com")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
