package watcherclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

const (
	DefaultAPIVersion = "v1"
	DefaultTimeout    = 30 * time.Second
)

// Client represents the Watcher API client
type Client struct {
	endpoint      string
	authenticator interface{ GetToken() (string, error) }
	httpClient    *http.Client
	apiVersion    string
}

// ClientOptions represents client configuration options
type ClientOptions struct {
	AuthURL         string
	Username        string
	Password        string
	ProjectName     string
	ProjectDomainID string
	UserDomainID    string
	Region          string
	Timeout         time.Duration
	AllowReauth     bool // Enable automatic re-authentication
}

// NewClient creates a new Watcher client with Keystone authentication
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}

	// Build authentication options
	authOpts := &AuthOptions{
		IdentityEndpoint: opts.AuthURL,
		Username:         opts.Username,
		Password:         opts.Password,
		DomainID:         opts.UserDomainID,
		AllowReauth:      opts.AllowReauth,
		Scope: &gophercloud.AuthScope{
			ProjectName: opts.ProjectName,
			DomainID:    opts.ProjectDomainID,
		},
	}

	// Validate auth options
	if err := ValidateAuthOptions(authOpts); err != nil {
		return nil, fmt.Errorf("invalid auth options: %w", err)
	}

	// Create authenticator
	auth, err := NewAuthenticator(authOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}

	client := &Client{
		endpoint:      auth.GetEndpoint() + "/" + DefaultAPIVersion,
		authenticator: auth,
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		apiVersion: DefaultAPIVersion,
	}

	return client, nil
}

// NewClientWithToken creates a client with existing token and endpoint
func NewClientWithToken(endpoint, token string) *Client {
	return &Client{
		endpoint:      endpoint + "/" + DefaultAPIVersion,
		authenticator: NewTokenAuthenticator(endpoint, token),
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiVersion: DefaultAPIVersion,
	}
}

// NewClientWithAuthenticator creates a client with a custom authenticator
func NewClientWithAuthenticator(auth interface {
	GetToken() (string, error)
	GetEndpoint() string
}) *Client {
	return &Client{
		endpoint:      auth.GetEndpoint() + "/" + DefaultAPIVersion,
		authenticator: auth,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiVersion: DefaultAPIVersion,
	}
}

// doRequest performs an HTTP request with automatic token handling
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	// Get current valid token
	token, err := c.authenticator.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.endpoint+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "golang-watcherclient/1.0")

	// Set API version header if needed
	req.Header.Set("OpenStack-API-Version", "infra-optim "+c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Handle authentication errors
	if resp.StatusCode == http.StatusUnauthorized {
		// Try to re-authenticate if using full authenticator
		if auth, ok := c.authenticator.(*Authenticator); ok && auth.autoReauth {
			resp.Body.Close()
			if err := auth.Reauth(); err != nil {
				return nil, fmt.Errorf("re-authentication failed: %w", err)
			}
			// Retry the request with new token
			return c.doRequest(method, path, body)
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("authentication failed: token expired or invalid")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes),
			URL:        req.URL.String(),
			Method:     method,
		}
	}

	return resp, nil
}

// parseResponse parses JSON response into the provided interface
func parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if len(bodyBytes) == 0 {
		// Empty response is valid for some operations (like DELETE)
		return nil
	}

	if err := json.Unmarshal(bodyBytes, v); err != nil {
		return fmt.Errorf("failed to parse response: %w (body: %s)", err, string(bodyBytes))
	}

	return nil
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// GetTimeout returns the current HTTP client timeout
func (c *Client) GetTimeout() time.Duration {
	return c.httpClient.Timeout
}

// GetEndpoint returns the Watcher API endpoint
func (c *Client) GetEndpoint() string {
	return c.endpoint
}

// GetAPIVersion returns the current API version
func (c *Client) GetAPIVersion() string {
	return c.apiVersion
}

// SetAPIVersion sets the API version
func (c *Client) SetAPIVersion(version string) {
	c.apiVersion = version
}

// GetAuthInfo returns authentication information (if using full authenticator)
func (c *Client) GetAuthInfo() (*AuthInfo, error) {
	if auth, ok := c.authenticator.(*Authenticator); ok {
		info := auth.GetAuthInfo()
		return &info, nil
	}
	return nil, fmt.Errorf("auth info not available with token authenticator")
}

// Ping checks if the Watcher API is accessible
func (c *Client) Ping() error {
	resp, err := c.doRequest(http.MethodGet, "/", nil)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	resp.Body.Close()
	return nil
}

// GetVersion returns the API version information
func (c *Client) GetVersion() (map[string]interface{}, error) {
	// Remove /v1 from endpoint for version query
	baseEndpoint := c.endpoint[:len(c.endpoint)-len("/"+c.apiVersion)]

	req, err := http.NewRequest(http.MethodGet, baseEndpoint+"/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// APIError represents an API error response
type APIError struct {
	StatusCode int
	Message    string
	URL        string
	Method     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s %s returned %d: %s",
		e.Method, e.URL, e.StatusCode, e.Message)
}

// IsNotFound checks if the error is a 404 Not Found error
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsConflict checks if the error is a 409 Conflict error
func IsConflict(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusConflict
	}
	return false
}

// IsUnauthorized checks if the error is a 401 Unauthorized error
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a 403 Forbidden error
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// ClientConfig holds client configuration for advanced usage
type ClientConfig struct {
	Endpoint      string
	Timeout       time.Duration
	RetryCount    int
	RetryWaitTime time.Duration
	MaxRetryWait  time.Duration
	Debug         bool
	CustomHeaders map[string]string
}

// NewClientWithConfig creates a client with advanced configuration
func NewClientWithConfig(opts ClientOptions, config ClientConfig) (*Client, error) {
	client, err := NewClient(opts)
	if err != nil {
		return nil, err
	}

	if config.Timeout > 0 {
		client.SetTimeout(config.Timeout)
	}

	// TODO: Implement retry logic and custom headers if needed

	return client, nil
}
