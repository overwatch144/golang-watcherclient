package watcherclient

import (
	"fmt"
	"net/http"
)

// GetGoal retrieves a goal by UUID or name
func (c *Client) GetGoal(identifier string) (*Goal, error) {
	path := fmt.Sprintf("/goals/%s", identifier)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result Goal
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListGoals lists all available goals
func (c *Client) ListGoals(opts *ListOptions) ([]Goal, error) {
	path := "/goals"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result GoalsResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Goals, nil
}
