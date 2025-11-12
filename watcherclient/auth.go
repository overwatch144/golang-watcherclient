package watcherclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// AuthOptions contains authentication configuration
type AuthOptions struct {
	IdentityEndpoint            string
	Username                    string
	Password                    string
	UserID                      string
	ApplicationCredentialID     string
	ApplicationCredentialName   string
	ApplicationCredentialSecret string
	DomainID                    string
	DomainName                  string
	TenantID                    string
	TenantName                  string
	AllowReauth                 bool
	TokenID                     string
	Scope                       *gophercloud.AuthScope
}

// Authenticator handles authentication and token management
type Authenticator struct {
	authOptions *AuthOptions
	provider    *gophercloud.ProviderClient
	token       string
	tokenExpiry time.Time
	endpoint    string
	mutex       sync.RWMutex
	autoReauth  bool
}

// NewAuthenticator creates a new authenticator instance
func NewAuthenticator(opts *AuthOptions) (*Authenticator, error) {
	if opts == nil {
		return nil, fmt.Errorf("auth options cannot be nil")
	}

	// Validate required fields
	if opts.IdentityEndpoint == "" {
		return nil, fmt.Errorf("identity endpoint is required")
	}

	auth := &Authenticator{
		authOptions: opts,
		autoReauth:  opts.AllowReauth,
	}

	// Perform initial authentication
	if err := auth.Authenticate(); err != nil {
		return nil, fmt.Errorf("initial authentication failed: %w", err)
	}

	return auth, nil
}

// Authenticate performs authentication against Keystone
func (a *Authenticator) Authenticate() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Build gophercloud auth options
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint:            a.authOptions.IdentityEndpoint,
		Username:                    a.authOptions.Username,
		UserID:                      a.authOptions.UserID,
		Password:                    a.authOptions.Password,
		DomainID:                    a.authOptions.DomainID,
		DomainName:                  a.authOptions.DomainName,
		TenantID:                    a.authOptions.TenantID,
		TenantName:                  a.authOptions.TenantName,
		AllowReauth:                 a.authOptions.AllowReauth,
		TokenID:                     a.authOptions.TokenID,
		Scope:                       a.authOptions.Scope,
		ApplicationCredentialID:     a.authOptions.ApplicationCredentialID,
		ApplicationCredentialName:   a.authOptions.ApplicationCredentialName,
		ApplicationCredentialSecret: a.authOptions.ApplicationCredentialSecret,
	}

	// Create authenticated client with context
	ctx := context.Background()
	provider, err := openstack.AuthenticatedClient(ctx, authOpts)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	a.provider = provider
	a.token = provider.TokenID

	// Get token expiry time if possible
	if err := a.updateTokenExpiry(); err != nil {
		// Non-fatal error, log but continue
		// Token expiry will be checked on demand
	}

	// Get Watcher endpoint
	endpointOpts := gophercloud.EndpointOpts{
		Type: "infra-optim",
	}

	endpoint, err := provider.EndpointLocator(endpointOpts)
	if err != nil {
		return fmt.Errorf("failed to locate Watcher endpoint: %w", err)
	}

	a.endpoint = endpoint

	return nil
}

// updateTokenExpiry extracts and updates token expiry time
func (a *Authenticator) updateTokenExpiry() error {
	if a.provider == nil {
		return fmt.Errorf("provider not initialized")
	}

	// Create identity v3 client
	identityClient, err := openstack.NewIdentityV3(a.provider, gophercloud.EndpointOpts{})
	if err != nil {
		return fmt.Errorf("failed to create identity client: %w", err)
	}

	// Get token details with context
	ctx := context.Background()
	tokenDetails, err := tokens.Get(ctx, identityClient, a.token).Extract()
	if err != nil {
		return fmt.Errorf("failed to get token details: %w", err)
	}

	a.tokenExpiry = tokenDetails.ExpiresAt
	return nil
}

// GetToken returns the current valid token
func (a *Authenticator) GetToken() (string, error) {
	a.mutex.RLock()
	token := a.token
	expiry := a.tokenExpiry
	a.mutex.RUnlock()

	// Check if token is expired or about to expire (5 min buffer)
	if !expiry.IsZero() && time.Until(expiry) < 5*time.Minute {
		if a.autoReauth {
			// Token expired or expiring soon, re-authenticate
			if err := a.Authenticate(); err != nil {
				return "", fmt.Errorf("failed to re-authenticate: %w", err)
			}
			// Get new token
			a.mutex.RLock()
			token = a.token
			a.mutex.RUnlock()
		} else {
			return "", fmt.Errorf("token expired and auto-reauth is disabled")
		}
	}

	return token, nil
}

// GetEndpoint returns the Watcher service endpoint
func (a *Authenticator) GetEndpoint() string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.endpoint
}

// GetProvider returns the gophercloud provider client
func (a *Authenticator) GetProvider() *gophercloud.ProviderClient {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.provider
}

// IsTokenExpired checks if the token is expired
func (a *Authenticator) IsTokenExpired() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if a.tokenExpiry.IsZero() {
		return false
	}

	return time.Now().After(a.tokenExpiry)
}

// GetTokenExpiry returns the token expiry time
func (a *Authenticator) GetTokenExpiry() time.Time {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.tokenExpiry
}

// Reauth forces re-authentication
func (a *Authenticator) Reauth() error {
	return a.Authenticate()
}

// TokenAuthenticator creates an authenticator with existing token
type TokenAuthenticator struct {
	endpoint string
	token    string
}

// NewTokenAuthenticator creates authenticator with existing token and endpoint
func NewTokenAuthenticator(endpoint, token string) *TokenAuthenticator {
	return &TokenAuthenticator{
		endpoint: endpoint,
		token:    token,
	}
}

// GetToken returns the token (no validation)
func (t *TokenAuthenticator) GetToken() (string, error) {
	if t.token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return t.token, nil
}

// GetEndpoint returns the endpoint
func (t *TokenAuthenticator) GetEndpoint() string {
	return t.endpoint
}

// ValidateAuthOptions validates authentication options
func ValidateAuthOptions(opts *AuthOptions) error {
	if opts == nil {
		return fmt.Errorf("auth options cannot be nil")
	}

	if opts.IdentityEndpoint == "" {
		return fmt.Errorf("identity endpoint is required")
	}

	// Check if we have valid authentication method
	hasPassword := opts.Username != "" && opts.Password != ""
	hasToken := opts.TokenID != ""
	hasAppCred := opts.ApplicationCredentialID != "" && opts.ApplicationCredentialSecret != ""
	hasAppCredName := opts.ApplicationCredentialName != "" && opts.ApplicationCredentialSecret != ""

	if !hasPassword && !hasToken && !hasAppCred && !hasAppCredName {
		return fmt.Errorf("no valid authentication method provided")
	}

	// Validate scope
	if opts.Scope != nil {
		if opts.Scope.ProjectID == "" && opts.Scope.ProjectName == "" &&
			opts.Scope.DomainID == "" && opts.Scope.DomainName == "" {
			return fmt.Errorf("scope must specify either project or domain")
		}
	}

	return nil
}

// BuildAuthOptions builds AuthOptions from ClientOptions
func BuildAuthOptions(opts ClientOptions) *AuthOptions {
	authOpts := &AuthOptions{
		IdentityEndpoint: opts.AuthURL,
		Username:         opts.Username,
		Password:         opts.Password,
		DomainID:         opts.UserDomainID,
		AllowReauth:      true,
	}

	// Build scope
	if opts.ProjectName != "" || opts.ProjectDomainID != "" {
		authOpts.Scope = &gophercloud.AuthScope{
			ProjectName: opts.ProjectName,
			DomainID:    opts.ProjectDomainID,
		}
	}

	return authOpts
}

// AuthInfo contains authentication information for debugging
type AuthInfo struct {
	Username        string
	UserID          string
	ProjectName     string
	ProjectID       string
	DomainName      string
	DomainID        string
	TokenExpiry     time.Time
	IsExpired       bool
	TimeUntilExpiry time.Duration
}

// GetAuthInfo returns current authentication information
func (a *Authenticator) GetAuthInfo() AuthInfo {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	info := AuthInfo{
		Username:    a.authOptions.Username,
		UserID:      a.authOptions.UserID,
		DomainName:  a.authOptions.DomainName,
		DomainID:    a.authOptions.DomainID,
		TokenExpiry: a.tokenExpiry,
	}

	if a.authOptions.Scope != nil {
		info.ProjectName = a.authOptions.Scope.ProjectName
		info.ProjectID = a.authOptions.Scope.ProjectID
	}

	if !a.tokenExpiry.IsZero() {
		info.IsExpired = time.Now().After(a.tokenExpiry)
		info.TimeUntilExpiry = time.Until(a.tokenExpiry)
	}

	return info
}

// SessionManager manages multiple authentication sessions
type SessionManager struct {
	sessions map[string]*Authenticator
	mutex    sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Authenticator),
	}
}

// AddSession adds a new authentication session
func (sm *SessionManager) AddSession(name string, auth *Authenticator) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.sessions[name] = auth
}

// GetSession retrieves an authentication session
func (sm *SessionManager) GetSession(name string) (*Authenticator, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	auth, exists := sm.sessions[name]
	if !exists {
		return nil, fmt.Errorf("session '%s' not found", name)
	}

	return auth, nil
}

// RemoveSession removes an authentication session
func (sm *SessionManager) RemoveSession(name string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.sessions, name)
}

// ListSessions returns all session names
func (sm *SessionManager) ListSessions() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	names := make([]string, 0, len(sm.sessions))
	for name := range sm.sessions {
		names = append(names, name)
	}
	return names
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	removed := 0
	for name, auth := range sm.sessions {
		if auth.IsTokenExpired() {
			delete(sm.sessions, name)
			removed++
		}
	}
	return removed
}
