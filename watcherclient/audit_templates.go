package watcherclient

import (
	"fmt"
	"net/http"
)

// CreateAuditTemplate creates a new audit template
func (c *Client) CreateAuditTemplate(template *AuditTemplate) (*AuditTemplate, error) {
	resp, err := c.doRequest(http.MethodPost, "/audit_templates", template)
	if err != nil {
		return nil, err
	}

	var result AuditTemplate
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAuditTemplate retrieves an audit template by UUID
func (c *Client) GetAuditTemplate(uuid string) (*AuditTemplate, error) {
	path := fmt.Sprintf("/audit_templates/%s", uuid)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result AuditTemplate
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListAuditTemplates lists all audit templates
func (c *Client) ListAuditTemplates(opts *ListOptions) ([]AuditTemplate, error) {
	path := "/audit_templates"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result AuditTemplatesResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.AuditTemplates, nil
}

// UpdateAuditTemplate updates an existing audit template
func (c *Client) UpdateAuditTemplate(uuid string, updates map[string]interface{}) (*AuditTemplate, error) {
	path := fmt.Sprintf("/audit_templates/%s", uuid)

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

	var result AuditTemplate
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAuditTemplate deletes an audit template
func (c *Client) DeleteAuditTemplate(uuid string) error {
	path := fmt.Sprintf("/audit_templates/%s", uuid)
	resp, err := c.doRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
