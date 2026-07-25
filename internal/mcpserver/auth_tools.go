package mcpserver

import (
	"context"

	lastfm "github.com/ndyakov/go-lastfm/v2"
)

func registerAuthTools(s toolServer, c *lastfm.Client) {
	addWriteTool(s, "lastfm_auth_get_token", "Create a temporary Last.fm browser-authentication token.", false, false, func(ctx context.Context, _ EmptyInput) (any, error) {
		return c.Auth.GetToken(ctx)
	})
	addWriteTool(s, "lastfm_auth_authorization_url", "Build the Last.fm browser authorization URL for a token.", false, true, func(_ context.Context, in AuthorizationURLInput) (any, error) {
		return map[string]string{"url": c.AuthorizationURL(in.Token, in.Callback)}, nil
	})
	addWriteTool(s, "lastfm_auth_get_session", "Exchange an authorized request token for a session and install it in this server process.", false, false, func(ctx context.Context, in TokenInput) (any, error) {
		return c.Auth.GetSession(ctx, in.Token)
	})
	addWriteTool(s, "lastfm_auth_get_mobile_session", "Authenticate with a Last.fm username and password and install the session in this process. Prefer browser authentication.", false, false, func(ctx context.Context, in MobileSessionInput) (any, error) {
		return c.Auth.GetMobileSession(ctx, in.Username, in.Password)
	})
}
