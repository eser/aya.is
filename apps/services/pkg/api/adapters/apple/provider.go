package apple

import (
	"context"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

// Provider implements auth.Provider for Apple Sign In.
type Provider struct {
	client *Client
}

// NewProvider creates a new Apple provider.
func NewProvider(client *Client) *Provider {
	return &Provider{
		client: client,
	}
}

// Client returns the underlying Apple client.
func (p *Provider) Client() *Client {
	return p.client
}

// SetFirstAuthUserInfo stores user info from Apple's first authorization POST body.
// Apple only sends user name/email in the POST body on the very first authorization.
// Implements auth.FirstAuthInfoSetter.
func (p *Provider) SetFirstAuthUserInfo(name, email string) {
	p.client.firstAuthUserInfo = &FirstAuthUserInfo{
		Name:  name,
		Email: email,
	}
}

// InitiateOAuth builds the Apple OAuth authorization URL.
func (p *Provider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	authURL := p.client.BuildAuthURL(redirectURI, state, p.client.config.Scope)

	p.client.logger.DebugContext(ctx, "Initiating Apple OAuth",
		slog.String("redirect_uri", redirectURI))

	return authURL, nil
}

// HandleOAuthCallback exchanges the code for tokens and returns account info.
func (p *Provider) HandleOAuthCallback(
	ctx context.Context,
	code string,
	redirectURI string,
) (auth.OAuthCallbackResult, error) {
	p.client.logger.DebugContext(ctx, "Processing Apple OAuth callback")

	// Exchange code for tokens
	tokenResp, err := p.client.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Parse id_token to get user info
	claims, err := p.client.ParseIDToken(ctx, tokenResp.IDToken)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	email := claims.Email
	name := ""

	// Apple only sends user info in the POST body on first authorization.
	// Use it if available, otherwise fall back to id_token claims.
	if p.client.firstAuthUserInfo != nil {
		if p.client.firstAuthUserInfo.Name != "" {
			name = p.client.firstAuthUserInfo.Name
		}

		if p.client.firstAuthUserInfo.Email != "" {
			email = p.client.firstAuthUserInfo.Email
		}

		// Clear after use — it's only valid for this single callback
		p.client.firstAuthUserInfo = nil
	}

	p.client.logger.DebugContext(ctx, "Apple OAuth callback successful",
		slog.String("sub", claims.Sub),
		slog.String("email", email))

	return auth.OAuthCallbackResult{
		RemoteID:             claims.Sub,
		Username:             "", // Apple doesn't provide a username
		Name:                 name,
		Email:                email,
		URI:                  "",
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresAt: nil,
		Scope:                "",
	}, nil
}
