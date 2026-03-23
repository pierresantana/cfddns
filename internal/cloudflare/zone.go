package cloudflare

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetZoneID returns the Cloudflare zone ID for the given zone name.
func (c *Client) GetZoneID(zoneName string) (string, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/zones?name="+zoneName, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}

	apiResp, err := c.do(req)
	if err != nil {
		return "", err
	}

	var zones []zone
	if err := json.Unmarshal(apiResp.Result, &zones); err != nil {
		return "", fmt.Errorf("parsing zones: %w", err)
	}

	if len(zones) == 0 {
		return "", fmt.Errorf("zone %q not found", zoneName)
	}

	return zones[0].ID, nil
}
