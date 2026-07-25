package mcpserver

import "time"

type EmptyInput struct{}

type PageInput struct {
	Page  int `json:"page,omitempty" jsonschema:"page number starting at 1"`
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of results"`
}

type ArtistInput struct {
	Artist      string `json:"artist,omitempty" jsonschema:"artist name; provide this or mbid"`
	MBID        string `json:"mbid,omitempty" jsonschema:"MusicBrainz artist ID; provide this or artist"`
	Autocorrect bool   `json:"autocorrect,omitempty" jsonschema:"allow Last.fm to correct the artist name"`
}

type AlbumInput struct {
	Artist      string `json:"artist,omitempty" jsonschema:"artist name; required with album when mbid is absent"`
	Album       string `json:"album,omitempty" jsonschema:"album name; required with artist when mbid is absent"`
	MBID        string `json:"mbid,omitempty" jsonschema:"MusicBrainz release ID"`
	Autocorrect bool   `json:"autocorrect,omitempty" jsonschema:"allow Last.fm to correct names"`
}

type TrackInput struct {
	Artist      string `json:"artist,omitempty" jsonschema:"artist name; required with track when mbid is absent"`
	Track       string `json:"track,omitempty" jsonschema:"track name; required with artist when mbid is absent"`
	MBID        string `json:"mbid,omitempty" jsonschema:"MusicBrainz recording ID"`
	Autocorrect bool   `json:"autocorrect,omitempty" jsonschema:"allow Last.fm to correct names"`
}

type ArtistPageInput struct {
	ArtistInput
	PageInput
}

type ArtistInfoInput struct {
	ArtistInput
	Username string `json:"username,omitempty" jsonschema:"Last.fm username for personalized data"`
	Language string `json:"language,omitempty" jsonschema:"ISO 639 language code for biography text"`
}

type ArtistLimitInput struct {
	ArtistInput
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of results"`
}

type ArtistUserInput struct {
	ArtistInput
	User string `json:"user" jsonschema:"Last.fm username"`
}

type AlbumInfoInput struct {
	AlbumInput
	Username string `json:"username,omitempty" jsonschema:"Last.fm username for personalized data"`
}

type AlbumUserInput struct {
	AlbumInput
	User string `json:"user" jsonschema:"Last.fm username"`
}

type TrackInfoInput struct {
	TrackInput
	Username string `json:"username,omitempty" jsonschema:"Last.fm username for personalized data"`
}

type TrackLimitInput struct {
	TrackInput
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of results"`
}

type TrackUserInput struct {
	TrackInput
	User string `json:"user" jsonschema:"Last.fm username"`
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"search query"`
	PageInput
}

type TrackSearchInput struct {
	Track  string `json:"track" jsonschema:"track name to search for"`
	Artist string `json:"artist,omitempty" jsonschema:"optional artist filter"`
	PageInput
}

type UserInput struct {
	User string `json:"user" jsonschema:"Last.fm username"`
}

type UserPageInput struct {
	User string `json:"user" jsonschema:"Last.fm username"`
	PageInput
	IncludeRecentTracks bool `json:"include_recent_tracks,omitempty" jsonschema:"include recent tracks for friends"`
}

type UserRecentInput struct {
	User string `json:"user" jsonschema:"Last.fm username"`
	PageInput
	From     int64 `json:"from,omitempty" jsonschema:"start Unix timestamp, inclusive"`
	To       int64 `json:"to,omitempty" jsonschema:"end Unix timestamp, inclusive"`
	Extended bool  `json:"extended,omitempty" jsonschema:"include extended artist data"`
}

type UserTopInput struct {
	User   string `json:"user" jsonschema:"Last.fm username"`
	Period string `json:"period,omitempty" jsonschema:"overall, 7day, 1month, 3month, 6month, or 12month"`
	PageInput
}

type UserTopTagsInput struct {
	User  string `json:"user" jsonschema:"Last.fm username"`
	Limit int    `json:"limit,omitempty" jsonschema:"maximum number of tags"`
}

type UserWeeklyInput struct {
	User string `json:"user" jsonschema:"Last.fm username"`
	From int64  `json:"from,omitempty" jsonschema:"chart start Unix timestamp"`
	To   int64  `json:"to,omitempty" jsonschema:"chart end Unix timestamp"`
}

type PersonalTagsInput struct {
	User        string `json:"user" jsonschema:"Last.fm username"`
	Tag         string `json:"tag" jsonschema:"tag to filter by"`
	TaggingType string `json:"tagging_type" jsonschema:"artists, albums, or tracks"`
	PageInput
}

type TagInput struct {
	Tag string `json:"tag" jsonschema:"Last.fm tag"`
}

type TagInfoInput struct {
	Tag      string `json:"tag" jsonschema:"Last.fm tag"`
	Language string `json:"language,omitempty" jsonschema:"ISO 639 language code"`
}

type TagPageInput struct {
	Tag string `json:"tag" jsonschema:"Last.fm tag"`
	PageInput
}

type GeoInput struct {
	Country  string `json:"country" jsonschema:"country name as defined by ISO 3166-1"`
	Location string `json:"location,omitempty" jsonschema:"optional location within the country"`
	PageInput
}

type TagsWriteInput struct {
	Artist string   `json:"artist" jsonschema:"artist name"`
	Album  string   `json:"album,omitempty" jsonschema:"album name when tagging an album"`
	Track  string   `json:"track,omitempty" jsonschema:"track name when tagging a track"`
	Tags   []string `json:"tags" jsonschema:"one to ten tags"`
}

type RemoveTagInput struct {
	Artist string `json:"artist" jsonschema:"artist name"`
	Album  string `json:"album,omitempty" jsonschema:"album name when untagging an album"`
	Track  string `json:"track,omitempty" jsonschema:"track name when untagging a track"`
	Tag    string `json:"tag" jsonschema:"tag to remove"`
}

type NowPlayingInput struct {
	Artist      string `json:"artist" jsonschema:"artist name"`
	Track       string `json:"track" jsonschema:"track name"`
	Album       string `json:"album,omitempty"`
	AlbumArtist string `json:"album_artist,omitempty"`
	MBID        string `json:"mbid,omitempty"`
	TrackNumber int    `json:"track_number,omitempty"`
	Duration    int    `json:"duration,omitempty" jsonschema:"duration in seconds"`
}

type ScrobbleInput struct {
	Tracks []ScrobbleTrack `json:"tracks" jsonschema:"one to fifty tracks"`
}

type ScrobbleTrack struct {
	Artist       string `json:"artist"`
	Track        string `json:"track"`
	Album        string `json:"album,omitempty"`
	AlbumArtist  string `json:"album_artist,omitempty"`
	MBID         string `json:"mbid,omitempty"`
	Timestamp    string `json:"timestamp" jsonschema:"RFC 3339 timestamp when listening began"`
	ChosenByUser *bool  `json:"chosen_by_user,omitempty"`
	TrackNumber  int    `json:"track_number,omitempty"`
	Duration     int    `json:"duration,omitempty" jsonschema:"duration in seconds"`
}

func (s ScrobbleTrack) time() (time.Time, error) { return time.Parse(time.RFC3339, s.Timestamp) }

type TokenInput struct {
	Token string `json:"token" jsonschema:"authorized Last.fm request token"`
}

type AuthorizationURLInput struct {
	Token    string `json:"token,omitempty"`
	Callback string `json:"callback,omitempty" jsonschema:"absolute callback URL"`
}

type MobileSessionInput struct {
	Username string `json:"username"`
	Password string `json:"password" jsonschema:"Last.fm password; avoid this tool when browser authentication is possible"`
}

type CountryInput struct {
	Country string `json:"country,omitempty"`
}
type NeighboursInput struct {
	User  string `json:"user"`
	Limit int    `json:"limit,omitempty"`
}
type TopFansInput struct {
	Artist string `json:"artist,omitempty"`
	Track  string `json:"track,omitempty"`
	MBID   string `json:"mbid,omitempty"`
}
type TasteInput struct {
	Type1  string `json:"type1"`
	Value1 string `json:"value1"`
	Type2  string `json:"type2"`
	Value2 string `json:"value2"`
	Limit  int    `json:"limit,omitempty"`
}
