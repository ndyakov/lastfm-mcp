# lastfm-mcp

An MCP server for the Last.fm API, built on [`go-lastfm/v2`](https://github.com/ndyakov/go-lastfm) and the [official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).

It exposes the complete supported Last.fm surface as discoverable MCP tools over:

- stdio, for local desktop and CLI clients;
- Streamable HTTP, for local services and remote deployments.

Read-only tools are enabled by default. Writes, authentication helpers, and unsupported legacy methods require explicit opt-in.

## Status

This project is prepared for its initial repository import. The module path is `github.com/ndyakov/lastfm-mcp`; change it before the first release if the repository will use a different owner or name.

## Requirements

- Go 1.25 or later
- A [Last.fm API account](https://www.last.fm/api/account/create)
- `LASTFM_API_KEY`
- `LASTFM_API_SECRET` and `LASTFM_SESSION_KEY` only when write tools are enabled

## Run with stdio

```sh
export LASTFM_API_KEY="your-api-key"
go run ./cmd/lastfm-mcp
```

Example MCP client configuration:

```json
{
  "mcpServers": {
    "lastfm": {
      "command": "/absolute/path/to/lastfm-mcp",
      "env": {
        "LASTFM_API_KEY": "your-api-key"
      }
    }
  }
}
```

Never place credentials in a configuration file that will be committed.

## Run with Streamable HTTP

```sh
export LASTFM_API_KEY="your-api-key"
export LASTFM_MCP_TRANSPORT="http"
export LASTFM_MCP_HTTP_ADDR="127.0.0.1:8080"
go run ./cmd/lastfm-mcp
```

The MCP endpoint is `http://127.0.0.1:8080/mcp`; health checks use `GET /healthz`.

Non-loopback binds are rejected unless `LASTFM_MCP_HTTP_TOKEN` is set. Clients then send:

```text
Authorization: Bearer <token>
```

The server applies Go's cross-origin protection around the MCP endpoint and sets conservative HTTP timeouts. Put TLS and production identity-aware authentication in a trusted reverse proxy.

## Configuration

| Environment variable | Default | Purpose |
| --- | --- | --- |
| `LASTFM_API_KEY` | required | Last.fm API key |
| `LASTFM_API_SECRET` | empty | Required for signed/authenticated methods |
| `LASTFM_SESSION_KEY` | empty | Required for write tools |
| `LASTFM_USER_AGENT` | `lastfm-mcp/<version>` | Outbound Last.fm user agent |
| `LASTFM_REQUEST_TIMEOUT` | `30s` | Last.fm request timeout |
| `LASTFM_MCP_TRANSPORT` | `stdio` | `stdio` or `http` |
| `LASTFM_MCP_HTTP_ADDR` | `127.0.0.1:8080` | HTTP listen address |
| `LASTFM_MCP_HTTP_PATH` | `/mcp` | Streamable HTTP endpoint |
| `LASTFM_MCP_HTTP_TOKEN` | empty | HTTP bearer token; mandatory on non-loopback binds |
| `LASTFM_ENABLE_WRITES` | `false` | Register mutating tools |
| `LASTFM_ENABLE_AUTH_TOOLS` | `false` | Register credential/session tools |
| `LASTFM_ENABLE_EXPERIMENTAL` | `false` | Register legacy methods absent from Last.fm's public index |

Equivalent command-line flags are available for transport, address, path, and feature gates. Run `lastfm-mcp -h` for details.

## Tool coverage

The default server registers every stable read endpoint supported by `go-lastfm/v2`:

- album: info, user tags, top tags, search;
- artist: correction, info, similar, user tags, top albums/tracks/tags, search;
- chart: top artists, tracks, tags;
- geo: top artists and tracks;
- library: artists;
- tag: info, similar, top albums/artists/tracks/tags, weekly chart ranges;
- track: correction, info, similar, user tags, top tags, search;
- user: profile, friends, loved and recent tracks, personal tags, top albums/artists/tracks/tags, weekly charts and ranges.

`LASTFM_ENABLE_WRITES=true` adds album/artist/track tagging, love/unlove, now-playing, and scrobbling. These require an API secret and session key.

`LASTFM_ENABLE_AUTH_TOOLS=true` adds browser-token, authorization-URL, session-exchange, and mobile-session tools. These may expose session credentials to the MCP client and model; leave them disabled unless needed. Prefer completing authentication outside the MCP session and providing `LASTFM_SESSION_KEY` at process startup.

`LASTFM_ENABLE_EXPERIMENTAL=true` adds legacy tag search, metros, neighbours, artist/track top fans, and tasteometer comparison. Last.fm may change or disable these without notice.

## Docker

```sh
docker build -t lastfm-mcp .
docker run --rm -p 127.0.0.1:8080:8080 \
  -e LASTFM_API_KEY \
  -e LASTFM_MCP_TRANSPORT=http \
  -e LASTFM_MCP_HTTP_ADDR=0.0.0.0:8080 \
  -e LASTFM_MCP_HTTP_TOKEN \
  lastfm-mcp
```

`docker compose up --build` uses the same HTTP setup and reads values from the environment.

## Development

```sh
go mod tidy
go test -race ./...
go vet ./...
```

Until `go-lastfm/v2` is published, use a temporary Go workspace containing this repository and the adjacent `go-lastfm` checkout. Do not commit a local `replace` directive.

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and [ROADMAP.md](ROADMAP.md).

## License

[MIT](LICENSE)
