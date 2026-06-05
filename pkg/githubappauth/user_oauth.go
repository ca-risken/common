package githubappauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const (
	defaultGitHubOAuthBaseURL = "https://github.com/"
)

var defaultAllowedOAuthBaseURLHosts = map[string]struct{}{
	"github.com": {},
}

var newGitHubOAuthHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
}

// OAuthConfig is the GitHub App user-to-server OAuth credential set.
type OAuthConfig struct {
	ClientID                 string
	ClientSecret             string
	OAuthBaseURL             string
	APIBaseURL               string
	Scopes                   []string
	AllowedOAuthBaseURLHosts []string
	AllowedAPIBaseURLHosts   []string
}

type OAuthClient struct {
	oauthConfig       *oauth2.Config
	apiBaseURL        *url.URL
	allowedAPIHosts   map[string]struct{}
	allowedOAuthHosts map[string]struct{}
}

func NewOAuthClient(conf *OAuthConfig) (*OAuthClient, error) {
	if conf == nil || (conf.ClientID == "" && conf.ClientSecret == "") {
		return &OAuthClient{}, nil
	}
	if conf.ClientID == "" || conf.ClientSecret == "" {
		return nil, errors.New("github app oauth client id and client secret are required together")
	}
	client := &OAuthClient{
		allowedAPIHosts:   allowedBaseURLHosts(conf.AllowedAPIBaseURLHosts),
		allowedOAuthHosts: allowedOAuthBaseURLHosts(conf.AllowedOAuthBaseURLHosts),
	}
	oauthBaseURL, err := client.validateOAuthBaseURL(conf.OAuthBaseURL)
	if err != nil {
		return nil, err
	}
	apiBaseURL, err := client.validateAPIBaseURL(conf.APIBaseURL)
	if err != nil {
		return nil, err
	}
	client.apiBaseURL = apiBaseURL
	client.oauthConfig = &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauthBaseURL.ResolveReference(&url.URL{Path: "login/oauth/authorize"}).String(),
			TokenURL: oauthBaseURL.ResolveReference(&url.URL{Path: "login/oauth/access_token"}).String(),
		},
	}
	return client, nil
}

func (c *OAuthClient) Enabled() bool {
	return c != nil && c.oauthConfig != nil
}

func allowedOAuthBaseURLHosts(configuredHosts []string) map[string]struct{} {
	allowedHosts := make(map[string]struct{}, len(defaultAllowedOAuthBaseURLHosts)+len(configuredHosts))
	for host := range defaultAllowedOAuthBaseURLHosts {
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

func (c *OAuthClient) validateOAuthBaseURL(baseURL string) (*url.URL, error) {
	if baseURL == "" {
		baseURL = defaultGitHubOAuthBaseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("github app oauth_base_url must use https: %s", baseURL)
	}
	if _, ok := c.allowedOAuthHosts[strings.ToLower(u.Hostname())]; !ok {
		return nil, fmt.Errorf("github app oauth_base_url host is not allowed: %s", u.Hostname())
	}
	if u.Path == "" || !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u, nil
}

func (c *OAuthClient) validateAPIBaseURL(baseURL string) (*url.URL, error) {
	if baseURL == "" {
		return nil, nil
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("github app api_base_url must use https: %s", baseURL)
	}
	if _, ok := c.allowedAPIHosts[strings.ToLower(u.Hostname())]; !ok {
		return nil, fmt.Errorf("github app api_base_url host is not allowed: %s", u.Hostname())
	}
	if u.Path == "" || !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u, nil
}

func (c *OAuthClient) AuthorizationURL(state string) (string, error) {
	if !c.Enabled() {
		return "", errors.New("github app oauth is not configured")
	}
	if state == "" {
		return "", errors.New("github app oauth state is required")
	}
	return c.oauthConfig.AuthCodeURL(state), nil
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	if !c.Enabled() {
		return nil, errors.New("github app oauth is not configured")
	}
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("github app oauth code is required")
	}
	token, err := c.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange github app oauth code: %w", err)
	}
	if token == nil || token.AccessToken == "" {
		return nil, errors.New("github app oauth access token is empty")
	}
	return token, nil
}

func (c *OAuthClient) NewUserClient(ctx context.Context, token *oauth2.Token) (*github.Client, error) {
	if !c.Enabled() {
		return nil, errors.New("github app oauth is not configured")
	}
	if token == nil || token.AccessToken == "" {
		return nil, errors.New("github app oauth access token is required")
	}
	httpClient := newGitHubOAuthHTTPClient(ctx, token)
	client := github.NewClient(httpClient)
	if c.apiBaseURL != nil {
		client.BaseURL = c.apiBaseURL
	}
	return client, nil
}

func (c *OAuthClient) GetAuthenticatedUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	client, err := c.NewUserClient(ctx, token)
	if err != nil {
		return nil, err
	}
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("get authenticated github user: %w", err)
	}
	if user == nil || user.GetLogin() == "" {
		return nil, errors.New("authenticated github user login is empty")
	}
	return user, nil
}
