package app

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	lastfm "github.com/ndyakov/go-lastfm/v2"
	lastfmmcp "github.com/ndyakov/lastfm-mcp/internal/mcpserver"
)

func Run(ctx context.Context, args []string, getenv envLookup, version string) error {
	cfg, err := ParseConfig(args, getenv, os.Stderr, version)
	if err != nil {
		return err
	}
	httpClient := &http.Client{Timeout: cfg.RequestTimeout}
	client, err := lastfm.New(cfg.APIKey, cfg.APISecret, lastfm.WithHTTPClient(httpClient), lastfm.WithUserAgent(cfg.UserAgent))
	if err != nil {
		return err
	}
	client.SetSessionKey(cfg.SessionKey)
	server := lastfmmcp.New(client, lastfmmcp.Options{
		Version:            version,
		EnableWrites:       cfg.EnableWrites,
		EnableAuthTools:    cfg.EnableAuthTools,
		EnableExperimental: cfg.EnableExperimental,
	})

	if cfg.Transport == "stdio" {
		return server.Run(ctx, &mcp.StdioTransport{})
	}
	return runHTTP(ctx, cfg, server)
}

func runHTTP(ctx context.Context, cfg Config, server *mcp.Server) error {
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{
		SessionTimeout: 30 * time.Minute,
	})
	protection := http.NewCrossOriginProtection()
	mux := http.NewServeMux()
	mux.Handle(cfg.HTTPPath, bearerAuth(cfg.HTTPToken, protection.Handler(mcpHandler)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	errCh := make(chan error, 1)
	go func() {
		slog.Info("serving Last.fm MCP", "address", cfg.HTTPAddr, "path", cfg.HTTPPath)
		errCh <- httpServer.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("HTTP shutdown: %w", err)
		}
		return nil
	}
}

func bearerAuth(token string, next http.Handler) http.Handler {
	if token == "" {
		return next
	}
	want := []byte("Bearer " + token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := []byte(r.Header.Get("Authorization"))
		if len(got) != len(want) || subtle.ConstantTimeCompare(got, want) != 1 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="lastfm-mcp"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
