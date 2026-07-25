package mcpserver

import (
	"context"

	lastfm "github.com/ndyakov/go-lastfm/v2"
)

func registerExperimentalTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_experimental_tag_search", "Search tags using a legacy method that Last.fm may change or disable.", func(ctx context.Context, in SearchInput) (any, error) {
		return c.Experimental.SearchTags(ctx, lastfm.TagSearchParams{Tag: in.Query, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_experimental_geo_get_metros", "Get legacy Last.fm metro areas, optionally filtered by country.", func(ctx context.Context, in CountryInput) (any, error) {
		return c.Experimental.GetMetros(ctx, in.Country)
	})
	addReadTool(s, "lastfm_experimental_user_get_neighbours", "Get a user's musical neighbours using a legacy method.", func(ctx context.Context, in NeighboursInput) (any, error) {
		return c.Experimental.GetNeighbours(ctx, lastfm.NeighboursParams{User: in.User, Limit: in.Limit})
	})
	addReadTool(s, "lastfm_experimental_artist_get_top_fans", "Get an artist's top fans using a legacy method.", func(ctx context.Context, in TopFansInput) (any, error) {
		return c.Experimental.GetArtistTopFans(ctx, lastfm.TopFansParams{Artist: in.Artist, MBID: in.MBID})
	})
	addReadTool(s, "lastfm_experimental_track_get_top_fans", "Get a track's top fans using a legacy method.", func(ctx context.Context, in TopFansInput) (any, error) {
		return c.Experimental.GetTrackTopFans(ctx, lastfm.TopFansParams{Artist: in.Artist, Track: in.Track, MBID: in.MBID})
	})
	addReadTool(s, "lastfm_experimental_tasteometer_compare", "Compare two users or artists using the legacy tasteometer method.", func(ctx context.Context, in TasteInput) (any, error) {
		return c.Experimental.CompareTaste(ctx, lastfm.TasteometerParams{Type1: in.Type1, Value1: in.Value1, Type2: in.Type2, Value2: in.Value2, Limit: in.Limit})
	})
}
