package structs

import "github.com/mattn/go-nulltype"

type ArtistResponse struct {
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Spotify string `json:"spotify"`
}

type ArtistCache struct {
	Key      string             `json:"key"`
	Spotify  string             `json:"spotify"`
	Url      string             `json:"url"`
	CachedAt nulltype.NullInt64 `json:"cachedAt"`
}

type TrackCache struct {
	Name       string  `json:"key"`
	Artist     string  `json:"artist"`
	Album      string  `json:"album"`
	Cover      string  `json:"cover"`
	Spotify    *string `json:"spotify"`
	Musixmatch *string `json:"musixmatch"`
	Deezer     *string `json:"deezer"`
	Duration   *int    `json:"duration"` // Milliseconds
	CachedAt   int     `json:"cachedAt"`
}

type TrackResponse struct {
	Hash       string  `json:"hash"`
	Name       string  `json:"name"`
	Artist     string  `json:"artist"`
	Album      string  `json:"album"`
	Cover      string  `json:"cover"`
	Spotify    *string `json:"spotify"`
	Musixmatch *string `json:"musixmatch"`
	Deezer     *string `json:"deezer"`
	Duration   *int    `json:"duration"`
}

type AlbumCache struct {
	Name     string             `json:"key"`
	Artist   string             `json:"artist"`
	Cover    string             `json:"cover"`
	Spotify  string             `json:"spotify"`
	CachedAt nulltype.NullInt64 `json:"cachedAt"`
}

type AlbumResponse struct {
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	Artist  string `json:"artist"`
	Cover   string `json:"cover"`
	Spotify string `json:"spotify"`
}
