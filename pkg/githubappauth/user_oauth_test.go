package githubappauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewOAuthClient(t *testing.T) {
	cases := []struct {
		name      string
		conf      *OAuthConfig
		wantOAuth bool
		wantError bool
	}{
		{
			name: "OK no oauth",
		},
		{
			name: "OK empty oauth",
			conf: &OAuthConfig{},
		},
		{
			name:      "OK oauth",
			conf:      &OAuthConfig{ClientID: "client-id", ClientSecret: "client-secret"},
			wantOAuth: true,
		},
		{
			name:      "NG missing client secret",
			conf:      &OAuthConfig{ClientID: "client-id"},
			wantError: true,
		},
		{
			name:      "NG untrusted oauth host",
			conf:      &OAuthConfig{ClientID: "client-id", ClientSecret: "client-secret", OAuthBaseURL: "https://attacker.example"},
			wantError: true,
		},
		{
			name:      "NG untrusted api host",
			conf:      &OAuthConfig{ClientID: "client-id", ClientSecret: "client-secret", APIBaseURL: "https://attacker.example"},
			wantError: true,
		},
		{
			name:      "NG relative redirect url",
			conf:      &OAuthConfig{ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "/api/v1/code/github-app/oauth/callback"},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, err := NewOAuthClient(c.conf)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got := client.Enabled(); got != c.wantOAuth {
				t.Fatalf("Unexpected GitHub App OAuth support: want=%t, got=%t", c.wantOAuth, got)
			}
		})
	}
}

func TestAuthorizationURL(t *testing.T) {
	client, err := NewOAuthClient(&OAuthConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://risken.example/api/v1/code/github-app/oauth/callback",
		Scopes:       []string{"read:org"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	got, err := client.AuthorizationURL("state-value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse auth url: %v", err)
	}
	if u.String() != "https://github.com/login/oauth/authorize?client_id=client-id&redirect_uri=https%3A%2F%2Frisken.example%2Fapi%2Fv1%2Fcode%2Fgithub-app%2Foauth%2Fcallback&response_type=code&scope=read%3Aorg&state=state-value" {
		t.Fatalf("Unexpected auth url: %s", u.String())
	}
}

func TestExchangeCodeAndGetAuthenticatedUser(t *testing.T) {
	var gotCode string
	var gotRedirectURL string
	var gotAuthorization string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/oauth/access_token":
			if r.Method != http.MethodPost {
				t.Fatalf("Unexpected token method: %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse form: %v", err)
			}
			gotCode = r.Form.Get("code")
			gotRedirectURL = r.Form.Get("redirect_uri")
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"access_token":"user-token","token_type":"bearer"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		case "/user":
			gotAuthorization = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{"login": "octocat"}); err != nil {
				t.Fatalf("write user response: %v", err)
			}
		default:
			t.Fatalf("Unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}

	client, err := NewOAuthClient(&OAuthConfig{
		ClientID:                 "client-id",
		ClientSecret:             "client-secret",
		OAuthBaseURL:             server.URL,
		APIBaseURL:               server.URL,
		RedirectURL:              "https://risken.example/api/v1/code/github-app/oauth/callback",
		AllowedOAuthBaseURLHosts: []string{serverURL.Hostname()},
		AllowedAPIBaseURLHosts:   []string{serverURL.Hostname()},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	token, err := client.ExchangeCode(ctx, "oauth-code")
	if err != nil {
		t.Fatalf("Unexpected exchange error: %v", err)
	}
	if gotCode != "oauth-code" {
		t.Fatalf("Unexpected oauth code: %s", gotCode)
	}
	if gotRedirectURL != "https://risken.example/api/v1/code/github-app/oauth/callback" {
		t.Fatalf("Unexpected redirect_uri: %s", gotRedirectURL)
	}

	origNewGitHubOAuthHTTPClient := newGitHubOAuthHTTPClient
	newGitHubOAuthHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
		client := server.Client()
		client.Transport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base:   client.Transport,
		}
		return client
	}
	defer func() {
		newGitHubOAuthHTTPClient = origNewGitHubOAuthHTTPClient
	}()

	user, err := client.GetAuthenticatedUser(context.Background(), token)
	if err != nil {
		t.Fatalf("Unexpected user error: %v", err)
	}
	if user.GetLogin() != "octocat" {
		t.Fatalf("Unexpected user login: %s", user.GetLogin())
	}
	if !strings.HasPrefix(gotAuthorization, "Bearer ") {
		t.Fatalf("Unexpected authorization header: %s", gotAuthorization)
	}
}

func TestOAuthClientErrors(t *testing.T) {
	client, err := NewOAuthClient(nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := client.AuthorizationURL("state"); err == nil {
		t.Fatal("Expected authorization url error but got none")
	}

	client, err = NewOAuthClient(&OAuthConfig{ClientID: "client-id", ClientSecret: "client-secret"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := client.AuthorizationURL(""); err == nil {
		t.Fatal("Expected empty state error but got none")
	}
	if _, err := client.ExchangeCode(context.Background(), ""); err == nil {
		t.Fatal("Expected empty code error but got none")
	}
	if _, err := client.GetAuthenticatedUser(context.Background(), nil); err == nil {
		t.Fatal("Expected empty token error but got none")
	}
}
