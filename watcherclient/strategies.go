package watcherclient

import (
	"fmt"
	"net/http"
)

// GetStrategy retrieves a strategy by UUID or name
func (c *Client) GetStrategy(identifier string) (*Strategy, error) {
	path := fmt.Sprintf("/strategies/%s", identifier)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result Strategy
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStrategies lists all available strategies
func (c *Client) ListStrategies(opts *ListOptions) ([]Strategy, error) {
	path := "/strategies"
	if opts != nil {
		path += buildQueryString(opts)
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result StrategiesResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Strategies, nil
}

// ListStrategiesByGoal lists strategies for a specific goal
func (c *Client) ListStrategiesByGoal(goalIdentifier string) ([]Strategy, error) {
	path := fmt.Sprintf("/goals/%s/strategies", goalIdentifier)
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result StrategiesResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Strategies, nil
}
