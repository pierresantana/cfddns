package cloudflare

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDo_SetsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer my-token" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer my-token")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type header = %q, want %q", got, "application/json")
		}
		w.Write([]byte(`{"success":true,"result":null}`))
	}))
	defer server.Close()

	client := NewClient("my-token")
	client.baseURL = server.URL

	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	_, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDo_UnsuccessfulNoErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"success":false,"errors":[]}`))
	}))
	defer server.Close()

	client := NewClient("token")
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	_, err := client.do(req)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "cloudflare API returned unsuccessful response" {
		t.Errorf("unexpected error: %s", got)
	}
}
