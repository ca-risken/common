package githubappauth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const (
	githubAppJWTBackdate = time.Minute
	githubAppJWTLifetime = 9 * time.Minute
)

var defaultAllowedBaseURLHosts = map[string]struct{}{
	"api.github.com": {},
}

var newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
}

// Config is the server-side GitHub App credential set.
type Config struct {
	AppID               string
	PrivateKey          string
	AllowedBaseURLHosts []string
}

type InstallationTokenConfig struct {
	BaseURL        string
	InstallationID uint64
}

type Client struct {
	appID               int64
	privateKey          *rsa.PrivateKey
	allowedBaseURLHosts map[string]struct{}
}

func NewClient(conf *Config) (*Client, error) {
	if conf == nil || (conf.AppID == "" && conf.PrivateKey == "") {
		return &Client{}, nil
	}
	if conf.AppID == "" || conf.PrivateKey == "" {
		return nil, errors.New("github app id and private key are required together")
	}
	appID, err := strconv.ParseInt(conf.AppID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse github app id: %w", err)
	}
	privateKey, err := ParsePrivateKey(conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse github app private key: %w", err)
	}
	return &Client{
		appID:               appID,
		privateKey:          privateKey,
		allowedBaseURLHosts: allowedBaseURLHosts(conf.AllowedBaseURLHosts),
	}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.privateKey != nil
}

func allowedBaseURLHosts(configuredHosts []string) map[string]struct{} {
	allowedHosts := make(map[string]struct{}, len(defaultAllowedBaseURLHosts)+len(configuredHosts))
	for host := range defaultAllowedBaseURLHosts {
		allowedHosts[host] = struct{}{}
	}
	for _, host := range configuredHosts {
		normalized := strings.ToLower(strings.TrimSpace(host))
		if normalized == "" {
			continue
		}
		allowedHosts[normalized] = struct{}{}
	}
	return allowedHosts
}

func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	normalized := strings.ReplaceAll(privateKey, `\n`, "\n")
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(normalized))
}

func (c *Client) CreateJWT(now time.Time) (string, error) {
	if !c.Enabled() {
		return "", errors.New("github app auth is not configured")
	}
	claims := jwt.RegisteredClaims{
		Issuer:    strconv.FormatInt(c.appID, 10),
		IssuedAt:  jwt.NewNumericDate(now.Add(-githubAppJWTBackdate)),
		ExpiresAt: jwt.NewNumericDate(now.Add(githubAppJWTLifetime)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign github app jwt: %w", err)
	}
	return signed, nil
}

func (c *Client) NewGitHubClient(ctx context.Context, baseURL string) (*github.Client, error) {
	base, err := c.ValidateBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	jwtToken, err := c.CreateJWT(time.Now())
	if err != nil {
		return nil, err
	}
	httpClient := newGitHubAppHTTPClient(ctx, &oauth2.Token{AccessToken: jwtToken, TokenType: "Bearer"})
	client := github.NewClient(httpClient)
	if base != nil {
		client.BaseURL = base
	}
	return client, nil
}

func (c *Client) ValidateBaseURL(baseURL string) (*url.URL, error) {
	if !c.Enabled() {
		return nil, errors.New("github app auth is not configured")
	}
	if baseURL == "" {
		return nil, nil
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("github app base_url must use https: %s", baseURL)
	}
	if _, ok := c.allowedBaseURLHosts[strings.ToLower(u.Hostname())]; !ok {
		return nil, fmt.Errorf("github app base_url host is not allowed: %s", u.Hostname())
	}
	if u.Path == "" || !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u, nil
}

func (c *Client) ResolveInstallationToken(ctx context.Context, config *InstallationTokenConfig, repoName string) (string, error) {
	if !c.Enabled() {
		return "", errors.New("github app auth is not configured")
	}
	if config == nil {
		return "", errors.New("github app installation token config is required")
	}
	if config.InstallationID == 0 {
		return "", errors.New("installation_id is required")
	}
	client, err := c.NewGitHubClient(ctx, config.BaseURL)
	if err != nil {
		return "", fmt.Errorf("create github app client: %w", err)
	}
	opts, err := InstallationTokenOptions(repoName)
	if err != nil {
		return "", err
	}
	token, _, err := client.Apps.CreateInstallationToken(ctx, int64(config.InstallationID), opts)
	if err != nil {
		return "", fmt.Errorf("create installation token: %w", err)
	}
	if token.GetToken() == "" {
		return "", errors.New("installation token is empty")
	}
	return token.GetToken(), nil
}

func InstallationTokenOptions(repoName string) (*github.InstallationTokenOptions, error) {
	repoName = strings.TrimSpace(repoName)
	if repoName == "" {
		return nil, nil
	}
	parts := strings.Split(repoName, "/")
	switch len(parts) {
	case 1:
		if parts[0] == "" {
			return nil, fmt.Errorf("invalid repository name: %s", repoName)
		}
	case 2:
		if parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid repository full name: %s", repoName)
		}
		repoName = parts[1]
	default:
		return nil, fmt.Errorf("invalid repository full name: %s", repoName)
	}
	return &github.InstallationTokenOptions{
		Repositories: []string{repoName},
	}, nil
}
