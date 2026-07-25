package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	lastfm "github.com/ndyakov/go-lastfm/v2"
)

type Options struct {
	Version            string
	EnableWrites       bool
	EnableAuthTools    bool
	EnableExperimental bool
}

func New(client *lastfm.Client, options Options) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "lastfm-mcp", Version: options.Version}, nil)
	registerReadTools(server, client)
	if options.EnableWrites {
		registerWriteTools(server, client)
	}
	if options.EnableAuthTools {
		registerAuthTools(server, client)
	}
	if options.EnableExperimental {
		registerExperimentalTools(server, client)
	}
	return server
}

func addReadTool[In any](server *mcp.Server, name, description string, call func(context.Context, In) (any, error)) {
	openWorld := true
	addTool(server, name, description, &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}, call)
}

func addWriteTool[In any](server *mcp.Server, name, description string, destructive, idempotent bool, call func(context.Context, In) (any, error)) {
	openWorld := true
	addTool(server, name, description, &mcp.ToolAnnotations{
		DestructiveHint: &destructive,
		IdempotentHint:  idempotent,
		OpenWorldHint:   &openWorld,
	}, call)
}

func addTool[In any](server *mcp.Server, name, description string, annotations *mcp.ToolAnnotations, call func(context.Context, In) (any, error)) {
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: description, Annotations: annotations},
		func(ctx context.Context, _ *mcp.CallToolRequest, input In) (*mcp.CallToolResult, any, error) {
			value, err := call(ctx, input)
			if err != nil {
				return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}}}, nil, nil
			}
			if value == nil {
				value = map[string]any{"ok": true}
			}
			encoded, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				return nil, nil, fmt.Errorf("encode tool result: %w", err)
			}
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(encoded)}}}, nil, nil
		})
}
