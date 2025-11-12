package watcherclient

import (
	"net/http"
)

// GetDataModel retrieves the infrastructure data model
func (c *Client) GetDataModel(dataModelType string) (*DataModel, error) {
	path := "/data_model"
	if dataModelType != "" {
		path += "?type=" + dataModelType
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result DataModel
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
