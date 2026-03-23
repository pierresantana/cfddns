package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type dnsRecordPayload struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// ListDNSRecords returns all A records matching the given hostname in the zone.
func (c *Client) ListDNSRecords(zoneID, hostname string) ([]DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?type=A&name=%s", c.baseURL, zoneID, hostname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	apiResp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var records []DNSRecord
	if err := json.Unmarshal(apiResp.Result, &records); err != nil {
		return nil, fmt.Errorf("parsing records: %w", err)
	}

	return records, nil
}

// CreateDNSRecord creates a new A record in the zone.
func (c *Client) CreateDNSRecord(zoneID, hostname, ip string, ttl int) error {
	payload := dnsRecordPayload{
		Type:    "A",
		Name:    hostname,
		Content: ip,
		TTL:     ttl,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling payload: %w", err)
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records", c.baseURL, zoneID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	_, err = c.do(req)
	return err
}

// UpdateDNSRecord updates an existing A record with a new IP.
func (c *Client) UpdateDNSRecord(zoneID, recordID, hostname, ip string, ttl int) error {
	payload := dnsRecordPayload{
		Type:    "A",
		Name:    hostname,
		Content: ip,
		TTL:     ttl,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling payload: %w", err)
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.baseURL, zoneID, recordID)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	_, err = c.do(req)
	return err
}
