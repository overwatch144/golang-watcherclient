package watcherclient

import (
	"fmt"
	"net/http"
)

// CreateAudit creates a new audit
func (c *Client) CreateAudit(audit *Audit) (*Audit, error) {
	resp, err := c.doRequest(http.MethodPost, "/audits", audit)
	if err != nil {
		return nil, err
	}

	var result Audit
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAudit retrieves an audit by UUID
func (c *Client) GetAudit(uuid string) (*Audit, error) {
	path := fmt.Sprintf("/audits/%s", uuid)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result Audit
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListAudits lists all audits
func (c *Client) ListAudits(opts *ListOptions) ([]Audit, error) {
	path := "/audits"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result AuditsResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Audits, nil
}

// UpdateAudit updates an existing audit
func (c *Client) UpdateAudit(uuid string, updates map[string]interface{}) (*Audit, error) {
	path := fmt.Sprintf("/audits/%s", uuid)

	// Watcher API uses PATCH with RFC 6902 JSON Patch format
	patches := []map[string]interface{}{}
	for key, value := range updates {
		patches = append(patches, map[string]interface{}{
			"op":    "replace",
			"path":  "/" + key,
			"value": value,
		})
	}

	resp, err := c.doRequest(http.MethodPatch, path, patches)
	if err != nil {
		return nil, err
	}

	var result Audit
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAudit deletes an audit
func (c *Client) DeleteAudit(uuid string) error {
	path := fmt.Sprintf("/audits/%s", uuid)
	resp, err := c.doRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// StartAudit starts an audit (changes state to ONGOING)
func (c *Client) StartAudit(uuid string) (*Audit, error) {
	return c.UpdateAudit(uuid, map[string]interface{}{
		"state": "ONGOING",
	})
}
