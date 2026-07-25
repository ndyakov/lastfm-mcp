package mcpserver

import (
	"context"
	"fmt"

	lastfm "github.com/ndyakov/go-lastfm/v2"
)

func registerWriteTools(s toolServer, c *lastfm.Client) {
	addWriteTool(s, "lastfm_album_add_tags", "Add user tags to an album.", false, true, func(ctx context.Context, in TagsWriteInput) (any, error) {
		return nil, c.Album.AddTags(ctx, in.Artist, in.Album, in.Tags)
	})
	addWriteTool(s, "lastfm_album_remove_tag", "Remove one of the authenticated user's tags from an album.", true, true, func(ctx context.Context, in RemoveTagInput) (any, error) {
		return nil, c.Album.RemoveTag(ctx, in.Artist, in.Album, in.Tag)
	})
	addWriteTool(s, "lastfm_artist_add_tags", "Add user tags to an artist.", false, true, func(ctx context.Context, in TagsWriteInput) (any, error) {
		return nil, c.Artist.AddTags(ctx, in.Artist, in.Tags)
	})
	addWriteTool(s, "lastfm_artist_remove_tag", "Remove one of the authenticated user's tags from an artist.", true, true, func(ctx context.Context, in RemoveTagInput) (any, error) {
		return nil, c.Artist.RemoveTag(ctx, in.Artist, in.Tag)
	})
	addWriteTool(s, "lastfm_track_add_tags", "Add user tags to a track.", false, true, func(ctx context.Context, in TagsWriteInput) (any, error) {
		return nil, c.Track.AddTags(ctx, in.Artist, in.Track, in.Tags)
	})
	addWriteTool(s, "lastfm_track_remove_tag", "Remove one of the authenticated user's tags from a track.", true, true, func(ctx context.Context, in RemoveTagInput) (any, error) {
		return nil, c.Track.RemoveTag(ctx, in.Artist, in.Track, in.Tag)
	})
	addWriteTool(s, "lastfm_track_love", "Love a track for the authenticated user.", false, true, func(ctx context.Context, in TrackInput) (any, error) {
		return nil, c.Track.Love(ctx, in.Artist, in.Track)
	})
	addWriteTool(s, "lastfm_track_unlove", "Remove a loved track from the authenticated user's library.", true, true, func(ctx context.Context, in TrackInput) (any, error) {
		return nil, c.Track.Unlove(ctx, in.Artist, in.Track)
	})
	addWriteTool(s, "lastfm_track_update_now_playing", "Set the authenticated user's currently playing track.", false, false, func(ctx context.Context, in NowPlayingInput) (any, error) {
		return c.Track.UpdateNowPlaying(ctx, lastfm.NowPlayingParams{Artist: in.Artist, Track: in.Track, Album: in.Album, AlbumArtist: in.AlbumArtist, MBID: in.MBID, TrackNumber: in.TrackNumber, Duration: in.Duration})
	})
	addWriteTool(s, "lastfm_track_scrobble", "Scrobble one to fifty completed or in-progress listens.", false, false, func(ctx context.Context, in ScrobbleInput) (any, error) {
		entries := make([]lastfm.Scrobble, len(in.Tracks))
		for i, item := range in.Tracks {
			timestamp, err := item.time()
			if err != nil {
				return nil, fmt.Errorf("tracks[%d].timestamp: %w", i, err)
			}
			entries[i] = lastfm.Scrobble{Artist: item.Artist, Track: item.Track, Album: item.Album, AlbumArtist: item.AlbumArtist, MBID: item.MBID, Timestamp: timestamp, ChosenByUser: item.ChosenByUser, TrackNumber: item.TrackNumber, Duration: item.Duration}
		}
		return c.Track.Scrobble(ctx, entries)
	})
}
