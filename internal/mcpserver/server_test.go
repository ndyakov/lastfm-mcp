package mcpserver

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	lastfm "github.com/ndyakov/go-lastfm/v2"
)

func TestToolSetsAndReadCall(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripper(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("method"); got != "user.getInfo" {
			t.Errorf("method = %q", got)
		}
		return jsonResponse(`{"user":{"name":"alice"}}`), nil
	})}

	client, err := lastfm.New("key", "secret", lastfm.WithBaseURL("https://example.invalid/"), lastfm.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	client.SetSessionKey("session")

	for _, test := range []struct {
		name  string
		opts  Options
		count int
	}{
		{"read only", Options{}, 44},
		{"writes", Options{EnableWrites: true}, 54},
		{"auth", Options{EnableAuthTools: true}, 48},
		{"experimental", Options{EnableExperimental: true}, 50},
		{"all", Options{EnableWrites: true, EnableAuthTools: true, EnableExperimental: true}, 64},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := New(client, test.opts)
			session := connect(t, server)
			tools, err := session.ListTools(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(tools.Tools) != test.count {
				t.Fatalf("got %d tools, want %d", len(tools.Tools), test.count)
			}
			if test.name == "read only" {
				result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "lastfm_user_get_info", Arguments: map[string]any{"user": "alice"}})
				if err != nil {
					t.Fatal(err)
				}
				if result.IsError || len(result.Content) != 1 {
					t.Fatalf("unexpected result: %#v", result)
				}
			}
		})
	}
}

func TestToolReturnsAPIErrorAsToolError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripper(func(*http.Request) (*http.Response, error) {
		return jsonResponse(`{"error":6,"message":"not found"}`), nil
	})}
	client, _ := lastfm.New("key", "", lastfm.WithBaseURL("https://example.invalid/"), lastfm.WithHTTPClient(httpClient))
	session := connect(t, New(client, Options{}))
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "lastfm_user_get_info", Arguments: map[string]any{"user": "missing"}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected MCP tool error")
	}
}

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(request *http.Request) (*http.Response, error) { return r(request) }

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func connect(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "test"}, nil)
	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = clientSession.Close()
		_ = serverSession.Close()
	})
	return clientSession
}
