package cloudflare

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListDNSRecords(t *testing.T) {
	t.Run("returns matching records", func(t *testing.T) {
		records := []DNSRecord{
			{ID: "rec-1", Type: "A", Name: "home.example.com", Content: "1.2.3.4", TTL: 300},
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if q := r.URL.Query().Get("type"); q != "A" {
				t.Errorf("expected type=A, got %s", q)
			}
			if q := r.URL.Query().Get("name"); q != "home.example.com" {
				t.Errorf("expected name=home.example.com, got %s", q)
			}
			resp := apiResponse{Success: true, Result: mustMarshal(records)}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := newTestClient(server)
		got, err := client.ListDNSRecords("zone-1", "home.example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 record, got %d", len(got))
		}
		if got[0].Content != "1.2.3.4" {
			t.Errorf("expected content 1.2.3.4, got %s", got[0].Content)
		}
	})

	t.Run("returns empty for no records", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := apiResponse{Success: true, Result: mustMarshal([]DNSRecord{})}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := newTestClient(server)
		got, err := client.ListDNSRecords("zone-1", "missing.example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected 0 records, got %d", len(got))
		}
	})
}

func TestCreateDNSRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload dnsRecordPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse payload: %v", err)
		}
		if payload.Type != "A" {
			t.Errorf("expected type A, got %s", payload.Type)
		}
		if payload.Name != "new.example.com" {
			t.Errorf("expected name new.example.com, got %s", payload.Name)
		}
		if payload.Content != "5.6.7.8" {
			t.Errorf("expected content 5.6.7.8, got %s", payload.Content)
		}
		if payload.TTL != 120 {
			t.Errorf("expected TTL 120, got %d", payload.TTL)
		}

		resp := apiResponse{Success: true, Result: mustMarshal(map[string]string{"id": "new-rec"})}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	err := client.CreateDNSRecord("zone-1", "new.example.com", "5.6.7.8", 120)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateDNSRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload dnsRecordPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse payload: %v", err)
		}
		if payload.Content != "9.10.11.12" {
			t.Errorf("expected content 9.10.11.12, got %s", payload.Content)
		}

		resp := apiResponse{Success: true, Result: mustMarshal(map[string]string{"id": "rec-1"})}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	err := client.UpdateDNSRecord("zone-1", "rec-1", "home.example.com", "9.10.11.12", 300)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := apiResponse{
			Success: false,
			Errors:  []apiError{{Code: 9021, Message: "TTL must be between 60 and 86400"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	err := client.CreateDNSRecord("zone-1", "test.example.com", "1.2.3.4", 30)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "cloudflare API error: 9021 - TTL must be between 60 and 86400" {
		t.Errorf("unexpected error message: %s", got)
	}
}
