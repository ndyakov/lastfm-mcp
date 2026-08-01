# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/lastfm-mcp ./cmd/lastfm-mcp

FROM gcr.io/distroless/static-debian12:nonroot
# Required by the MCP Registry: proves this image belongs to the server
# entry named in server.json.
LABEL io.modelcontextprotocol.server.name="io.github.ndyakov/lastfm-mcp"
COPY --from=build /out/lastfm-mcp /usr/local/bin/lastfm-mcp
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/lastfm-mcp"]
