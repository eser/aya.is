package github

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

// Provider implements auth.Provider and profiles.LinkProvider for GitHub.
// Both interfaces share the same methods - the service layer handles the difference.
type Provider struct {
	client *Client
}

// NewProvider creates a new GitHub provider.
func NewProvider(client *Client) *Provider {
	return &Provider{
		client: client,
	}
}

// Client returns the underlying GitHub client.
func (p *Provider) Client() *Client {
	return p.client
}

// InitiateOAuth builds the OAuth URL with the given state.
// Implements auth.Provider interface for login.
// For login: caller passes auth.GenerateRandomState().
func (p *Provider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	authURL := p.client.BuildAuthURL(redirectURI, state, p.client.Config().Scope)

	p.client.Logger().DebugContext(ctx, "Initiating GitHub OAuth",
		slog.String("redirect_uri", redirectURI))

	return authURL, nil
}

// InitiateProfileLinkOAuth builds the OAuth URL with expanded scope for profile linking.
// Uses ProfileLinkScope which includes read:org for organization access.
func (p *Provider) InitiateProfileLinkOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	authURL := p.client.BuildAuthURL(redirectURI, state, p.client.Config().ProfileLinkScope)

	p.client.Logger().DebugContext(ctx, "Initiating GitHub OAuth for profile link",
		slog.String("redirect_uri", redirectURI),
		slog.String("scope", p.client.Config().ProfileLinkScope))

	return authURL, nil
}

// HandleOAuthCallback exchanges the code for tokens and returns account info.
// Implements auth.Provider and profiles.LinkProvider interfaces.
// State validation is handled by the service layer.
func (p *Provider) HandleOAuthCallback(
	ctx context.Context,
	code string,
	redirectURI string,
) (auth.OAuthCallbackResult, error) {
	p.client.Logger().DebugContext(ctx, "Processing GitHub OAuth callback")

	// Exchange code for tokens
	tokenResp, err := p.client.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Fetch user info
	userInfo, err := p.client.FetchUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	p.client.Logger().DebugContext(ctx, "GitHub OAuth callback successful",
		slog.Int64("user_id", userInfo.ID),
		slog.String("login", userInfo.Login))

	name := userInfo.Name
	if name == "" {
		name = userInfo.Login
	}

	return auth.OAuthCallbackResult{
		RemoteID:             strconv.FormatInt(userInfo.ID, 10),
		Username:             userInfo.Login,
		Name:                 name,
		Email:                userInfo.Email,
		URI:                  userInfo.HTMLURL,
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         "", // GitHub doesn't provide refresh tokens by default
		AccessTokenExpiresAt: nil,
		Scope:                p.client.Config().Scope,
	}, nil
}
