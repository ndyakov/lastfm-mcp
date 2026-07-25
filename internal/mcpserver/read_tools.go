package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	lastfm "github.com/ndyakov/go-lastfm/v2"
)

func registerReadTools(s toolServer, c *lastfm.Client) {
	registerAlbumTools(s, c)
	registerArtistTools(s, c)
	registerChartTools(s, c)
	registerGeoTools(s, c)
	registerLibraryTools(s, c)
	registerTagTools(s, c)
	registerTrackTools(s, c)
	registerUserTools(s, c)
}

// toolServer is an alias kept private so registration signatures remain compact.
type toolServer = *mcp.Server

func page(in PageInput) lastfm.Pagination { return lastfm.Pagination{Page: in.Page, Limit: in.Limit} }
func artist(in ArtistInput) lastfm.ArtistRef {
	return lastfm.ArtistRef{Artist: in.Artist, MBID: in.MBID, Autocorrect: in.Autocorrect}
}
func album(in AlbumInput) lastfm.AlbumRef {
	return lastfm.AlbumRef{Artist: in.Artist, Album: in.Album, MBID: in.MBID, Autocorrect: in.Autocorrect}
}
func track(in TrackInput) lastfm.TrackRef {
	return lastfm.TrackRef{Artist: in.Artist, Track: in.Track, MBID: in.MBID, Autocorrect: in.Autocorrect}
}

func registerAlbumTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_album_get_info", "Get album metadata, wiki, tags, and listener statistics.", func(ctx context.Context, in AlbumInfoInput) (any, error) {
		return c.Album.GetInfo(ctx, lastfm.AlbumInfoParams{AlbumRef: album(in.AlbumInput), Username: in.Username})
	})
	addReadTool(s, "lastfm_album_get_tags", "Get a user's tags for an album.", func(ctx context.Context, in AlbumUserInput) (any, error) {
		return c.Album.GetTags(ctx, lastfm.AlbumTagsParams{AlbumRef: album(in.AlbumInput), User: in.User})
	})
	addReadTool(s, "lastfm_album_get_top_tags", "Get the most popular tags for an album.", func(ctx context.Context, in AlbumInput) (any, error) {
		return c.Album.GetTopTags(ctx, album(in))
	})
	addReadTool(s, "lastfm_album_search", "Search Last.fm albums by name.", func(ctx context.Context, in SearchInput) (any, error) {
		return c.Album.Search(ctx, lastfm.AlbumSearchParams{Album: in.Query, Pagination: page(in.PageInput)})
	})
}

func registerArtistTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_artist_get_correction", "Get Last.fm's canonical correction for an artist name.", func(ctx context.Context, in ArtistInput) (any, error) {
		return c.Artist.GetCorrection(ctx, in.Artist)
	})
	addReadTool(s, "lastfm_artist_get_info", "Get artist metadata, biography, tags, and statistics.", func(ctx context.Context, in ArtistInfoInput) (any, error) {
		return c.Artist.GetInfo(ctx, lastfm.ArtistInfoParams{ArtistRef: artist(in.ArtistInput), Username: in.Username, Language: in.Language})
	})
	addReadTool(s, "lastfm_artist_get_similar", "Get artists similar to an artist.", func(ctx context.Context, in ArtistLimitInput) (any, error) {
		return c.Artist.GetSimilar(ctx, lastfm.ArtistSimilarParams{ArtistRef: artist(in.ArtistInput), Limit: in.Limit})
	})
	addReadTool(s, "lastfm_artist_get_tags", "Get a user's tags for an artist.", func(ctx context.Context, in ArtistUserInput) (any, error) {
		return c.Artist.GetTags(ctx, lastfm.ArtistTagsParams{ArtistRef: artist(in.ArtistInput), User: in.User})
	})
	addReadTool(s, "lastfm_artist_get_top_albums", "Get an artist's most popular albums.", func(ctx context.Context, in ArtistPageInput) (any, error) {
		return c.Artist.GetTopAlbums(ctx, lastfm.ArtistPagedParams{ArtistRef: artist(in.ArtistInput), Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_artist_get_top_tracks", "Get an artist's most popular tracks.", func(ctx context.Context, in ArtistPageInput) (any, error) {
		return c.Artist.GetTopTracks(ctx, lastfm.ArtistPagedParams{ArtistRef: artist(in.ArtistInput), Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_artist_get_top_tags", "Get the most popular tags for an artist.", func(ctx context.Context, in ArtistInput) (any, error) {
		return c.Artist.GetTopTags(ctx, artist(in))
	})
	addReadTool(s, "lastfm_artist_search", "Search Last.fm artists by name.", func(ctx context.Context, in SearchInput) (any, error) {
		return c.Artist.Search(ctx, lastfm.ArtistSearchParams{Artist: in.Query, Pagination: page(in.PageInput)})
	})
}

func registerChartTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_chart_get_top_artists", "Get the global top artists chart.", func(ctx context.Context, in PageInput) (any, error) { return c.Chart.GetTopArtists(ctx, page(in)) })
	addReadTool(s, "lastfm_chart_get_top_tags", "Get the global top tags chart.", func(ctx context.Context, in PageInput) (any, error) { return c.Chart.GetTopTags(ctx, page(in)) })
	addReadTool(s, "lastfm_chart_get_top_tracks", "Get the global top tracks chart.", func(ctx context.Context, in PageInput) (any, error) { return c.Chart.GetTopTracks(ctx, page(in)) })
}

func registerGeoTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_geo_get_top_artists", "Get the most popular artists in a country.", func(ctx context.Context, in GeoInput) (any, error) {
		return c.Geo.GetTopArtists(ctx, lastfm.GeoParams{Country: in.Country, Location: in.Location, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_geo_get_top_tracks", "Get the most popular tracks in a country or location.", func(ctx context.Context, in GeoInput) (any, error) {
		return c.Geo.GetTopTracks(ctx, lastfm.GeoParams{Country: in.Country, Location: in.Location, Pagination: page(in.PageInput)})
	})
}

func registerLibraryTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_library_get_artists", "Get all artists in a user's Last.fm library.", func(ctx context.Context, in UserPageInput) (any, error) {
		return c.Library.GetArtists(ctx, lastfm.LibraryArtistsParams{User: in.User, Pagination: page(in.PageInput)})
	})
}

func registerTagTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_tag_get_info", "Get tag metadata and wiki text.", func(ctx context.Context, in TagInfoInput) (any, error) {
		return c.Tag.GetInfo(ctx, lastfm.TagInfoParams{Tag: in.Tag, Language: in.Language})
	})
	addReadTool(s, "lastfm_tag_get_similar", "Get tags similar to a tag.", func(ctx context.Context, in TagInput) (any, error) { return c.Tag.GetSimilar(ctx, in.Tag) })
	addReadTool(s, "lastfm_tag_get_top_albums", "Get the top albums for a tag.", func(ctx context.Context, in TagPageInput) (any, error) {
		return c.Tag.GetTopAlbums(ctx, lastfm.TagPagedParams{Tag: in.Tag, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_tag_get_top_artists", "Get the top artists for a tag.", func(ctx context.Context, in TagPageInput) (any, error) {
		return c.Tag.GetTopArtists(ctx, lastfm.TagPagedParams{Tag: in.Tag, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_tag_get_top_tracks", "Get the top tracks for a tag.", func(ctx context.Context, in TagPageInput) (any, error) {
		return c.Tag.GetTopTracks(ctx, lastfm.TagPagedParams{Tag: in.Tag, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_tag_get_top_tags", "Get the globally most popular tags.", func(ctx context.Context, _ EmptyInput) (any, error) { return c.Tag.GetTopTags(ctx) })
	addReadTool(s, "lastfm_tag_get_weekly_chart_list", "Get available weekly chart date ranges for a tag.", func(ctx context.Context, in TagInput) (any, error) { return c.Tag.GetWeeklyChartList(ctx, in.Tag) })
}

func registerTrackTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_track_get_correction", "Get Last.fm's canonical correction for an artist and track name.", func(ctx context.Context, in TrackInput) (any, error) {
		return c.Track.GetCorrection(ctx, in.Artist, in.Track)
	})
	addReadTool(s, "lastfm_track_get_info", "Get track metadata, wiki, tags, and statistics.", func(ctx context.Context, in TrackInfoInput) (any, error) {
		return c.Track.GetInfo(ctx, lastfm.TrackInfoParams{TrackRef: track(in.TrackInput), Username: in.Username})
	})
	addReadTool(s, "lastfm_track_get_similar", "Get tracks similar to a track.", func(ctx context.Context, in TrackLimitInput) (any, error) {
		return c.Track.GetSimilar(ctx, lastfm.TrackSimilarParams{TrackRef: track(in.TrackInput), Limit: in.Limit})
	})
	addReadTool(s, "lastfm_track_get_tags", "Get a user's tags for a track.", func(ctx context.Context, in TrackUserInput) (any, error) {
		return c.Track.GetTags(ctx, lastfm.TrackTagsParams{TrackRef: track(in.TrackInput), User: in.User})
	})
	addReadTool(s, "lastfm_track_get_top_tags", "Get the most popular tags for a track.", func(ctx context.Context, in TrackInput) (any, error) { return c.Track.GetTopTags(ctx, track(in)) })
	addReadTool(s, "lastfm_track_search", "Search Last.fm tracks by name and optional artist.", func(ctx context.Context, in TrackSearchInput) (any, error) {
		return c.Track.Search(ctx, lastfm.TrackSearchParams{Track: in.Track, Artist: in.Artist, Pagination: page(in.PageInput)})
	})
}

func registerUserTools(s toolServer, c *lastfm.Client) {
	addReadTool(s, "lastfm_user_get_friends", "Get a user's Last.fm friends.", func(ctx context.Context, in UserPageInput) (any, error) {
		return c.User.GetFriends(ctx, lastfm.UserPagedParams{User: in.User, IncludeRecentTracks: in.IncludeRecentTracks, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_user_get_info", "Get a Last.fm user's profile and statistics.", func(ctx context.Context, in UserInput) (any, error) { return c.User.GetInfo(ctx, in.User) })
	addReadTool(s, "lastfm_user_get_loved_tracks", "Get tracks loved by a user.", func(ctx context.Context, in UserPageInput) (any, error) {
		return c.User.GetLovedTracks(ctx, lastfm.UserPagedParams{User: in.User, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_user_get_personal_tags", "Get items to which a user applied a tag.", func(ctx context.Context, in PersonalTagsInput) (any, error) {
		return c.User.GetPersonalTags(ctx, lastfm.PersonalTagsParams{User: in.User, Tag: in.Tag, TaggingType: in.TaggingType, Pagination: page(in.PageInput)})
	})
	addReadTool(s, "lastfm_user_get_recent_tracks", "Get a user's recent listening history.", func(ctx context.Context, in UserRecentInput) (any, error) {
		return c.User.GetRecentTracks(ctx, lastfm.RecentTracksParams{User: in.User, Pagination: page(in.PageInput), From: in.From, To: in.To, Extended: in.Extended})
	})
	addReadTool(s, "lastfm_user_get_top_albums", "Get a user's top albums for a period.", func(ctx context.Context, in UserTopInput) (any, error) { return c.User.GetTopAlbums(ctx, userTop(in)) })
	addReadTool(s, "lastfm_user_get_top_artists", "Get a user's top artists for a period.", func(ctx context.Context, in UserTopInput) (any, error) { return c.User.GetTopArtists(ctx, userTop(in)) })
	addReadTool(s, "lastfm_user_get_top_tracks", "Get a user's top tracks for a period.", func(ctx context.Context, in UserTopInput) (any, error) { return c.User.GetTopTracks(ctx, userTop(in)) })
	addReadTool(s, "lastfm_user_get_top_tags", "Get a user's most frequently used tags.", func(ctx context.Context, in UserTopTagsInput) (any, error) {
		return c.User.GetTopTags(ctx, lastfm.UserTopTagsParams{User: in.User, Limit: in.Limit})
	})
	addReadTool(s, "lastfm_user_get_weekly_album_chart", "Get a user's album chart for a weekly range.", func(ctx context.Context, in UserWeeklyInput) (any, error) {
		return c.User.GetWeeklyAlbumChart(ctx, weekly(in))
	})
	addReadTool(s, "lastfm_user_get_weekly_artist_chart", "Get a user's artist chart for a weekly range.", func(ctx context.Context, in UserWeeklyInput) (any, error) {
		return c.User.GetWeeklyArtistChart(ctx, weekly(in))
	})
	addReadTool(s, "lastfm_user_get_weekly_track_chart", "Get a user's track chart for a weekly range.", func(ctx context.Context, in UserWeeklyInput) (any, error) {
		return c.User.GetWeeklyTrackChart(ctx, weekly(in))
	})
	addReadTool(s, "lastfm_user_get_weekly_chart_list", "Get available weekly chart date ranges for a user.", func(ctx context.Context, in UserInput) (any, error) { return c.User.GetWeeklyChartList(ctx, in.User) })
}

func userTop(in UserTopInput) lastfm.UserTopParams {
	return lastfm.UserTopParams{User: in.User, Period: lastfm.Period(in.Period), Pagination: page(in.PageInput)}
}

func weekly(in UserWeeklyInput) lastfm.WeeklyChartParams {
	return lastfm.WeeklyChartParams{User: in.User, From: in.From, To: in.To}
}
