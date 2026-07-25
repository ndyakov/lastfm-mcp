package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func env(values map[string]string) envLookup { return func(key string) string { return values[key] } }

func TestParseConfigDefaults(t *testing.T) {
	cfg, err := ParseConfig(nil, env(map[string]string{"LASTFM_API_KEY": "key"}), io.Discard, "test")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Transport != "stdio" || cfg.HTTPAddr != "127.0.0.1:8080" || cfg.HTTPPath != "/mcp" {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
	if cfg.RequestTimeout != 30*time.Second || cfg.UserAgent != "lastfm-mcp/test" {
		t.Fatalf("unexpected client defaults: %#v", cfg)
	}
}

func TestParseConfigValidation(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		args []string
	}{
		{"missing key", nil, nil},
		{"transport", map[string]string{"LASTFM_API_KEY": "key"}, []string{"--transport", "tcp"}},
		{"path", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_MCP_HTTP_PATH": "mcp"}, nil},
		{"timeout", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_REQUEST_TIMEOUT": "bad"}, nil},
		{"writes", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_ENABLE_WRITES": "true"}, nil},
		{"auth", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_ENABLE_AUTH_TOOLS": "true"}, nil},
		{"public no token", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_MCP_TRANSPORT": "http", "LASTFM_MCP_HTTP_ADDR": "0.0.0.0:8080"}, nil},
		{"bad address", map[string]string{"LASTFM_API_KEY": "key", "LASTFM_MCP_TRANSPORT": "http", "LASTFM_MCP_HTTP_ADDR": "bad"}, nil},
		{"extra args", map[string]string{"LASTFM_API_KEY": "key"}, []string{"extra"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ParseConfig(test.args, env(test.env), io.Discard, "test"); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestParseConfigFlagsAndPublicToken(t *testing.T) {
	cfg, err := ParseConfig([]string{"--transport", "http", "--http-addr", "0.0.0.0:9090", "--enable-writes", "--enable-auth-tools", "--enable-experimental"}, env(map[string]string{
		"LASTFM_API_KEY": "key", "LASTFM_API_SECRET": "secret", "LASTFM_SESSION_KEY": "session", "LASTFM_MCP_HTTP_TOKEN": "token",
	}), io.Discard, "test")
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.EnableWrites || !cfg.EnableAuthTools || !cfg.EnableExperimental {
		t.Fatalf("flags not applied: %#v", cfg)
	}
}

func TestBearerAuth(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	handler := bearerAuth("secret", next)

	request := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", response.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/mcp", nil)
	request.Header.Set("Authorization", "Bearer secret")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}

	response = httptest.NewRecorder()
	bearerAuth("", next).ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/mcp", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
}
