package watcherclient

import (
	"fmt"
	"net/http"
)

// GetAction retrieves an action by UUID
func (c *Client) GetAction(uuid string) (*Action, error) {
	path := fmt.Sprintf("/actions/%s", uuid)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result Action
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListActions lists all actions
func (c *Client) ListActions(opts *ListOptions) ([]Action, error) {
	path := "/actions"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ActionsResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Actions, nil
}

// ListActionsByActionPlan lists actions for a specific action plan
func (c *Client) ListActionsByActionPlan(actionPlanUUID string) ([]Action, error) {
	path := fmt.Sprintf("/action_plans/%s/actions", actionPlanUUID)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ActionsResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Actions, nil
}
