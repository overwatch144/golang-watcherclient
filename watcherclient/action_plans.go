package watcherclient

import (
	"fmt"
	"net/http"
)

// GetActionPlan retrieves an action plan by UUID
func (c *Client) GetActionPlan(uuid string) (*ActionPlan, error) {
	path := fmt.Sprintf("/action_plans/%s", uuid)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ActionPlan
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListActionPlans lists all action plans
func (c *Client) ListActionPlans(opts *ListOptions) ([]ActionPlan, error) {
	path := "/action_plans"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ActionPlansResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.ActionPlans, nil
}

// UpdateActionPlan updates an existing action plan
func (c *Client) UpdateActionPlan(uuid string, updates map[string]interface{}) (*ActionPlan, error) {
	path := fmt.Sprintf("/action_plans/%s", uuid)

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

	var result ActionPlan
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteActionPlan deletes an action plan
func (c *Client) DeleteActionPlan(uuid string) error {
	path := fmt.Sprintf("/action_plans/%s", uuid)
	resp, err := c.doRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// StartActionPlan starts execution of an action plan
func (c *Client) StartActionPlan(uuid string) (*ActionPlan, error) {
	return c.UpdateActionPlan(uuid, map[string]interface{}{
		"state": "TRIGGERED",
	})
}

// CancelActionPlan cancels an action plan
func (c *Client) CancelActionPlan(uuid string) (*ActionPlan, error) {
	return c.UpdateActionPlan(uuid, map[string]interface{}{
		"state": "CANCELLED",
	})
}
