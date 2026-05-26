package githubappauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

func generateRSAPrivateKeyPEM(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa private key: %v", err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	return privateKey, string(pem.EncodeToMemory(block))
}

func TestNewClient(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	cases := []struct {
		name      string
		conf      *Config
		wantApp   bool
		wantError bool
	}{
		{
			name: "OK no app auth",
		},
		{
			name: "OK empty app auth",
			conf: &Config{},
		},
		{
			name:    "OK app auth",
			conf:    &Config{AppID: "12345", PrivateKey: privateKeyPEM},
			wantApp: true,
		},
		{
			name:      "NG missing private key",
			conf:      &Config{AppID: "12345"},
			wantError: true,
		},
		{
			name:      "NG invalid app id",
			conf:      &Config{AppID: "invalid", PrivateKey: privateKeyPEM},
			wantError: true,
		},
		{
			name:      "NG invalid private key",
			conf:      &Config{AppID: "12345", PrivateKey: "invalid"},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, err := NewClient(c.conf)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got := client.Enabled(); got != c.wantApp {
				t.Fatalf("Unexpected GitHub App support: want=%t, got=%t", c.wantApp, got)
			}
		})
	}
}

func TestParsePrivateKey(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	escapedPEM := strings.ReplaceAll(privateKeyPEM, "\n", `\n`)
	cases := []struct {
		name      string
		input     string
		wantError bool
	}{
		{name: "OK raw pem", input: privateKeyPEM},
		{name: "OK escaped pem", input: escapedPEM},
		{name: "NG invalid pem", input: "invalid", wantError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParsePrivateKey(c.input)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("Expected private key but got nil")
			}
		})
	}
}

func TestCreateJWT(t *testing.T) {
	privateKey, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewClient(&Config{AppID: "12345", PrivateKey: privateKeyPEM})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	now := time.Now().Truncate(time.Second)
	tokenString, err := client.CreateJWT(now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			t.Fatalf("Unexpected signing method: %v", token.Method.Alg())
		}
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("Expected valid token")
	}
	if claims.Issuer != "12345" {
		t.Fatalf("Unexpected issuer: %s", claims.Issuer)
	}
	if claims.IssuedAt == nil || claims.IssuedAt.Time.Unix() != now.Add(-githubAppJWTBackdate).Unix() {
		t.Fatalf("Unexpected issued at: %+v", claims.IssuedAt)
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Unix() != now.Add(githubAppJWTLifetime).Unix() {
		t.Fatalf("Unexpected expires at: %+v", claims.ExpiresAt)
	}
}

func TestResolveInstallationToken(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	var gotAuthorization string
	var gotRepositories []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/app/installations/12345/access_tokens" {
			t.Fatalf("Unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("Unexpected method: %s", r.Method)
		}
		gotAuthorization = r.Header.Get("Authorization")
		var body struct {
			Repositories []string `json:"repositories"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		gotRepositories = body.Repositories
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	origNewGitHubAppHTTPClient := newGitHubAppHTTPClient
	newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
		client := server.Client()
		client.Transport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base:   client.Transport,
		}
		return client
	}
	defer func() {
		newGitHubAppHTTPClient = origNewGitHubAppHTTPClient
	}()

	client, err := NewClient(&Config{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got, err := client.ResolveInstallationToken(context.Background(), &InstallationTokenConfig{
		BaseURL:        server.URL + "/",
		InstallationID: 12345,
	}, "owner/repo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "installation-token" {
		t.Fatalf("Unexpected token: %s", got)
	}
	if !strings.HasPrefix(gotAuthorization, "Bearer ") {
		t.Fatalf("Unexpected authorization header: %s", gotAuthorization)
	}
	if len(gotRepositories) != 1 || gotRepositories[0] != "repo" {
		t.Fatalf("Unexpected repositories: %+v", gotRepositories)
	}
}

func TestResolveInstallationTokenError(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := client.ResolveInstallationToken(context.Background(), &InstallationTokenConfig{InstallationID: 12345}, ""); err == nil {
		t.Fatal("Expected error but got none")
	}
}

func TestResolveInstallationTokenRejectsUntrustedBaseURL(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewClient(&Config{AppID: "12345", PrivateKey: privateKeyPEM})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_, err = client.ResolveInstallationToken(context.Background(), &InstallationTokenConfig{
		BaseURL:        "https://attacker.example/",
		InstallationID: 12345,
	}, "")
	if err == nil {
		t.Fatal("Expected error but got none")
	}
	if !strings.Contains(err.Error(), "base_url host is not allowed") {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateBaseURLAddsTrailingSlash(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewClient(&Config{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{"ghes.example.com"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	got, err := client.ValidateBaseURL("https://ghes.example.com/api/v3")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got.String() != "https://ghes.example.com/api/v3/" {
		t.Fatalf("Unexpected base URL: %s", got.String())
	}
}

func TestInstallationTokenOptions(t *testing.T) {
	cases := []struct {
		name      string
		repoName  string
		wantRepos []string
		wantError bool
	}{
		{name: "OK empty", repoName: ""},
		{name: "OK repo", repoName: "repo", wantRepos: []string{"repo"}},
		{name: "OK owner repo", repoName: "owner/repo", wantRepos: []string{"repo"}},
		{name: "NG empty repo", repoName: "/", wantError: true},
		{name: "NG empty owner", repoName: "/repo", wantError: true},
		{name: "NG empty repository", repoName: "owner/", wantError: true},
		{name: "NG too many parts", repoName: "a/b/c", wantError: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := InstallationTokenOptions(c.repoName)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got == nil {
				if len(c.wantRepos) != 0 {
					t.Fatalf("Unexpected nil options")
				}
				return
			}
			if len(got.Repositories) != len(c.wantRepos) {
				t.Fatalf("Unexpected repositories: %+v", got.Repositories)
			}
			for i := range got.Repositories {
				if got.Repositories[i] != c.wantRepos[i] {
					t.Fatalf("Unexpected repositories: %+v", got.Repositories)
				}
			}
		})
	}
}
