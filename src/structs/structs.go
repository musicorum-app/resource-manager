package structs

import "github.com/mattn/go-nulltype"

type ArtistResponse struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
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
	Name     string             `json:"key"`
	Artist   string             `json:"artist"`
	Album    string             `json:"album"`
	Spotify  string             `json:"spotify"`
	Duration nulltype.NullInt64 `json:"duration"` // Milliseconds
	CachedAt nulltype.NullInt64 `json:"cachedAt"`
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
	Name    string `json:"key"`
	Artist  string `json:"artist"`
	Cover   string `json:"cover"`
	Spotify string `json:"spotify"`
}
