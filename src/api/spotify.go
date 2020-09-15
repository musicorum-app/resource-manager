package api

import (
	"encoding/json"
	"fmt"
	"github.com/musicorum-app/resource-manager/utils"
	"github.com/rapito/go-spotify/spotify"
	"strings"
)

type ArtistObject struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Popularity uint8  `json:"popularity"`
}

type AlbumObject struct {
	ID     string
	Name   string
	Artist string
	Image  string
}

type TrackObject struct {
	ID       string
	Name     string
	Image    string
	Artist   string
	Album    string
	Duration int
}

type ImageItem struct {
	Url string
}

var spot spotify.Spotify

func Initialize() {
	spot = spotify.New(utils.GetEnvVar("SPOTIFY_ID"), utils.GetEnvVar("SPOTIFY_SECRET"))
	authorized, err := spot.Authorize()
	if !authorized {
		utils.FailOnError(err[0])
		println("NOT AUTHENTICATED")
	}
}

func SearchArtist(artist string) *ArtistObject {
	spot.Authorize()
	type ArtistResponseItem struct {
		ID         string
		Images     []ImageItem
		Name       string
		Popularity uint8
	}
	type ArtistsResponse struct {
		Items []ArtistResponseItem
	}
	type SpotifyResponse struct {
		Artists ArtistsResponse
	}
	search, err := spot.Get("search?type=artist&q=\"%s\"", nil, artist)
	if len(err) > 0 {
		utils.FailOnError(err[0])
	}
	var response SpotifyResponse
	resErr := json.Unmarshal(search, &response)
	if resErr != nil {
		println(resErr)
	}

	if len(response.Artists.Items) == 0 {
		return nil
	}
	result := new(ArtistObject)

	result.ID = response.Artists.Items[0].ID
	result.Name = response.Artists.Items[0].Name
	result.Image = response.Artists.Items[0].Images[0].Url
	result.Popularity = response.Artists.Items[0].Popularity

	return result
}

func SearchAlbum(album string, artist string) *AlbumObject {
	spot.Authorize()
	type ArtistItem struct {
		Name string
	}
	type AlbumResponseItem struct {
		ID      string       `json:"id"`
		Artists []ArtistItem `json:"artists"`
		Images  []ImageItem  `json:"images"`
		Name    string       `json:"name"`
	}
	type AlbumsResponse struct {
		Items []AlbumResponseItem `json:"items"`
	}
	type SpotifyResponse struct {
		Albums AlbumsResponse `json:"albums"`
	}
	search, err := spot.Get("search?type=album&q=\"%s\" artist:\"%s\"", nil, album, artist)
	if len(err) > 0 {
		fmt.Print(err)
		utils.FailOnError(err[0])
		return nil
	}
	var response SpotifyResponse
	resErr := json.Unmarshal(search, &response)
	if resErr != nil {
		fmt.Println(resErr)
		return nil
	}

	if len(response.Albums.Items) == 0 {
		return nil
	}
	result := new(AlbumObject)

	item := response.Albums.Items[0]
	result.ID = item.ID
	result.Name = item.Name
	result.Artist = item.Artists[0].Name
	result.Image = item.Images[0].Url

	return result
}

func SearchTrack(name string, album string, artist string) *TrackObject {
	spot.Authorize()
	type ArtistItem struct {
		Name string
	}
	type AlbumItem struct {
		Name   string      `json:"name"`
		Images []ImageItem `json:"images"`
	}
	type TrackResponseItem struct {
		ID       string       `json:"id"`
		Artists  []ArtistItem `json:"artists"`
		Album    AlbumItem    `json:"album"`
		Name     string       `json:"name"`
		Duration int          `json:"duration_ms"`
	}
	type TracksResponse struct {
		Items []TrackResponseItem `json:"items"`
	}
	type SpotifyResponse struct {
		Tracks TracksResponse `json:"tracks"`
	}
	albumSearch := ""
	if album != "" {
		albumSearch = fmt.Sprintf("album:\"%s\"", album)
	}
	search, err := spot.Get("search?type=track&q=\"%s\" artist:\"%s\" %s", nil, name, artist, albumSearch)
	if len(err) > 0 {
		fmt.Print(err)
		utils.FailOnError(err[0])
		return nil
	}
	var response SpotifyResponse
	resErr := json.Unmarshal(search, &response)
	if resErr != nil {
		fmt.Println(resErr)
		return nil
	}

	if len(response.Tracks.Items) == 0 {
		return nil
	}
	result := new(TrackObject)

	item := response.Tracks.Items[0]
	result.ID = item.ID
	result.Name = item.Name
	result.Artist = item.Artists[0].Name
	result.Image = item.Album.Images[0].Url
	result.Album = item.Album.Name
	result.Duration = item.Duration

	return result
}

func GetArtists(ids []string) []ArtistObject {
	spot.Authorize()
	type ArtistResponseItem struct {
		ID         string
		Images     []ImageItem
		Name       string
		Popularity uint8
	}
	type SpotifyResponse struct {
		Artists []ArtistResponseItem
	}
	search, err := spot.Get("artists?ids=%s", nil, strings.Join(ids, ","))
	if len(err) > 0 {
		utils.FailOnError(err[0])
	}
	var response SpotifyResponse
	resErr := json.Unmarshal(search, &response)
	if resErr != nil {
		println(resErr)
	}

	if len(response.Artists) == 0 {
		return nil
	}
	var result []ArtistObject

	for _, a := range response.Artists {
		artist := ArtistObject{
			ID:         a.ID,
			Name:       a.Name,
			Image:      a.Images[0].Url,
			Popularity: a.Popularity,
		}
		result = append(result, artist)
	}

	return result
}
