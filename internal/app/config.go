package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	APIKey             string
	APISecret          string
	SessionKey         string
	Transport          string
	HTTPAddr           string
	HTTPPath           string
	HTTPToken          string
	UserAgent          string
	EnableWrites       bool
	EnableAuthTools    bool
	EnableExperimental bool
	RequestTimeout     time.Duration
}

type envLookup func(string) string

func ParseConfig(args []string, getenv envLookup, stderr io.Writer, version string) (Config, error) {
	cfg := Config{
		APIKey:             getenv("LASTFM_API_KEY"),
		APISecret:          getenv("LASTFM_API_SECRET"),
		SessionKey:         getenv("LASTFM_SESSION_KEY"),
		Transport:          valueOr(getenv("LASTFM_MCP_TRANSPORT"), "stdio"),
		HTTPAddr:           valueOr(getenv("LASTFM_MCP_HTTP_ADDR"), "127.0.0.1:8080"),
		HTTPPath:           valueOr(getenv("LASTFM_MCP_HTTP_PATH"), "/mcp"),
		HTTPToken:          getenv("LASTFM_MCP_HTTP_TOKEN"),
		UserAgent:          valueOr(getenv("LASTFM_USER_AGENT"), "lastfm-mcp/"+version),
		EnableWrites:       envBool(getenv("LASTFM_ENABLE_WRITES")),
		EnableAuthTools:    envBool(getenv("LASTFM_ENABLE_AUTH_TOOLS")),
		EnableExperimental: envBool(getenv("LASTFM_ENABLE_EXPERIMENTAL")),
		RequestTimeout:     envDuration(getenv("LASTFM_REQUEST_TIMEOUT"), 30*time.Second),
	}

	fs := flag.NewFlagSet("lastfm-mcp", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&cfg.Transport, "transport", cfg.Transport, "transport: stdio or http")
	fs.StringVar(&cfg.HTTPAddr, "http-addr", cfg.HTTPAddr, "HTTP listen address")
	fs.StringVar(&cfg.HTTPPath, "http-path", cfg.HTTPPath, "Streamable HTTP endpoint path")
	fs.BoolVar(&cfg.EnableWrites, "enable-writes", cfg.EnableWrites, "enable tools that modify Last.fm data")
	fs.BoolVar(&cfg.EnableAuthTools, "enable-auth-tools", cfg.EnableAuthTools, "enable session authentication tools")
	fs.BoolVar(&cfg.EnableExperimental, "enable-experimental", cfg.EnableExperimental, "enable unsupported legacy methods")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if fs.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if cfg.APIKey == "" {
		return Config{}, errors.New("LASTFM_API_KEY is required")
	}
	if cfg.Transport != "stdio" && cfg.Transport != "http" {
		return Config{}, fmt.Errorf("unsupported transport %q; use stdio or http", cfg.Transport)
	}
	if cfg.HTTPPath == "" || !strings.HasPrefix(cfg.HTTPPath, "/") {
		return Config{}, errors.New("HTTP path must start with /")
	}
	if cfg.RequestTimeout <= 0 {
		return Config{}, errors.New("request timeout must be positive")
	}
	if cfg.EnableWrites && (cfg.APISecret == "" || cfg.SessionKey == "") {
		return Config{}, errors.New("writes require LASTFM_API_SECRET and LASTFM_SESSION_KEY")
	}
	if cfg.EnableAuthTools && cfg.APISecret == "" {
		return Config{}, errors.New("authentication tools require LASTFM_API_SECRET")
	}
	if cfg.Transport == "http" {
		host, _, err := net.SplitHostPort(cfg.HTTPAddr)
		if err != nil {
			return Config{}, fmt.Errorf("invalid HTTP address: %w", err)
		}
		if !isLoopbackHost(host) && cfg.HTTPToken == "" {
			return Config{}, errors.New("LASTFM_MCP_HTTP_TOKEN is required for non-loopback HTTP binds")
		}
	}
	return cfg, nil
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func envBool(value string) bool {
	parsed, _ := strconv.ParseBool(value)
	return parsed
}

func envDuration(value string, fallback time.Duration) time.Duration {
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return parsed
}

func isLoopbackHost(host string) bool {
	host = strings.Trim(host, "[]")
	return host == "localhost" || net.ParseIP(host).IsLoopback()
}
