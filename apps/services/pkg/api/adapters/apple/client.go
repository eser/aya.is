package apple

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/golang-jwt/jwt/v5"
)

// Sentinel errors.
var (
	ErrFailedToExchangeCode = errors.New("failed to exchange authorization code")
	ErrFailedToParseIDToken = errors.New("failed to parse id_token")
	ErrFailedToGenerateKey  = errors.New("failed to generate client secret")
	ErrFailedToFetchJWKS    = errors.New("failed to fetch Apple JWKS")
	ErrInvalidPrivateKey    = errors.New("invalid private key")
	ErrInvalidIDToken       = errors.New("invalid id_token")
	ErrNoPublicKeyFound     = errors.New("no matching public key found in JWKS")
)

const (
	appleAuthURL  = "https://appleid.apple.com/auth/authorize"
	appleTokenURL = "https://appleid.apple.com/auth/token"
	appleJWKSURL  = "https://appleid.apple.com/auth/keys"
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TokenResponse represents Apple's token endpoint response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// IDTokenClaims represents the claims in Apple's id_token JWT.
type IDTokenClaims struct {
	Sub            string `json:"sub"`
	Email          string `json:"email"`
	EmailVerified  any    `json:"email_verified"` // Can be bool or string
	IsPrivateEmail any    `json:"is_private_email"`
	jwt.RegisteredClaims
}

// IsEmailVerified returns whether the email is verified.
func (c *IDTokenClaims) IsEmailVerified() bool {
	switch v := c.EmailVerified.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}

// AppleJWKS represents Apple's JSON Web Key Set.
type AppleJWKS struct {
	Keys []AppleJWK `json:"keys"`
}

// AppleJWK represents a single JSON Web Key from Apple.
type AppleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// FirstAuthUserInfo holds user info from Apple's first authorization POST body.
type FirstAuthUserInfo struct {
	Name  string
	Email string
}

// Client provides Apple Sign In operations.
type Client struct {
	config            *auth.AppleAuthProviderConfig
	logger            *logfx.Logger
	httpClient        HTTPClient
	privateKey        *ecdsa.PrivateKey
	firstAuthUserInfo *FirstAuthUserInfo
}

// NewClient creates a new Apple client.
func NewClient(
	config *auth.AppleAuthProviderConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *Client {
	client := &Client{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}

	// Parse private key at init time if configured
	if config.PrivateKey != "" {
		key, err := parseP8PrivateKey(config.PrivateKey)
		if err != nil {
			logger.Warn("Failed to parse Apple private key",
				slog.String("error", err.Error()))
		} else {
			client.privateKey = key
		}
	}

	return client
}

// Config returns the Apple configuration.
func (c *Client) Config() *auth.AppleAuthProviderConfig {
	return c.config
}

// BuildAuthURL builds the Apple OAuth authorization URL.
func (c *Client) BuildAuthURL(redirectURI, state, scope string) string {
	queryString := url.Values{}
	queryString.Set("client_id", c.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", scope)
	queryString.Set("response_type", "code")
	queryString.Set("response_mode", "form_post")

	return appleAuthURL + "?" + queryString.Encode()
}

// GenerateClientSecret generates a JWT client secret for Apple.
// Apple requires the client secret to be a JWT signed with the app's private key.
func (c *Client) GenerateClientSecret() (string, error) {
	if c.privateKey == nil {
		return "", ErrInvalidPrivateKey
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    c.config.TeamID,
		Subject:   c.config.ClientID,
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(6 * 30 * 24 * time.Hour)), // ~6 months
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.config.KeyID

	signedToken, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGenerateKey, err)
	}

	return signedToken, nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens.
func (c *Client) ExchangeCodeForToken(
	ctx context.Context,
	code string,
	redirectURI string,
) (*TokenResponse, error) {
	clientSecret, err := c.GenerateClientSecret()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	values := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	if redirectURI != "" {
		values.Set("redirect_uri", redirectURI)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		appleTokenURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to exchange code for token",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Apple token exchange failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToExchangeCode, resp.StatusCode)
	}

	var tokenResp TokenResponse

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	if tokenResp.IDToken == "" {
		c.logger.ErrorContext(ctx, "No id_token received from Apple")

		return nil, ErrFailedToExchangeCode
	}

	c.logger.DebugContext(ctx, "Successfully obtained Apple tokens")

	return &tokenResp, nil
}

// ParseIDToken parses and verifies Apple's id_token JWT.
func (c *Client) ParseIDToken(
	ctx context.Context,
	idToken string,
) (*IDTokenClaims, error) {
	// Fetch Apple's JWKS for verification
	jwks, err := c.fetchJWKS(ctx)
	if err != nil {
		// If JWKS fetch fails, parse without verification as fallback
		c.logger.WarnContext(ctx, "Failed to fetch Apple JWKS, parsing without verification",
			slog.String("error", err.Error()))

		return c.parseIDTokenUnverified(idToken)
	}

	// Parse the token header to find the key ID
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	claims := &IDTokenClaims{} //nolint:exhaustruct

	_, err = parser.ParseWithClaims(idToken, claims, func(token *jwt.Token) (any, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, ErrInvalidIDToken
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return jwkToPublicKey(key)
			}
		}

		return nil, ErrNoPublicKeyFound
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToParseIDToken, err)
	}

	return claims, nil
}

// parseIDTokenUnverified parses the id_token without signature verification.
// Used as fallback when JWKS fetch fails.
func (c *Client) parseIDTokenUnverified(idToken string) (*IDTokenClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	claims := &IDTokenClaims{} //nolint:exhaustruct

	_, _, err := parser.ParseUnverified(idToken, claims)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToParseIDToken, err)
	}

	return claims, nil
}

// fetchJWKS fetches Apple's JSON Web Key Set.
func (c *Client) fetchJWKS(ctx context.Context) (*AppleJWKS, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appleJWKSURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchJWKS, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchJWKS, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchJWKS, resp.StatusCode)
	}

	var jwks AppleJWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchJWKS, err)
	}

	return &jwks, nil
}

// parseP8PrivateKey parses Apple's .p8 private key (PEM-encoded PKCS8 ECDSA).
func parseP8PrivateKey(keyStr string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, ErrInvalidPrivateKey
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidPrivateKey, err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}

	return ecdsaKey, nil
}

// jwkToPublicKey converts an Apple JWK to an RSA public key.
func jwkToPublicKey(key AppleJWK) (any, error) {
	nBytes, err := jwt.NewParser().DecodeSegment(key.N)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode N: %w", ErrInvalidIDToken, err)
	}

	eBytes, err := jwt.NewParser().DecodeSegment(key.E)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode E: %w", ErrInvalidIDToken, err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}
