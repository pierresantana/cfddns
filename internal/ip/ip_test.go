package ip

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	t.Run("returns trimmed IP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("  203.0.113.42\n"))
		}))
		defer server.Close()

		got, err := getPublicIPFrom(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "203.0.113.42" {
			t.Errorf("got %q, want %q", got, "203.0.113.42")
		}
	})

	t.Run("error on empty response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(""))
		}))
		defer server.Close()

		_, err := getPublicIPFrom(server.URL)
		if err == nil {
			t.Fatal("expected error for empty response")
		}
	})

	t.Run("error on non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, err := getPublicIPFrom(server.URL)
		if err == nil {
			t.Fatal("expected error for 500 status")
		}
	})
}
