package api

import (
	"encoding/json"
	"fmt"
	"github.com/musicorum-app/resource-manager/utils"
	"github.com/rapito/go-spotify/spotify"
)

type ArtistObject struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type AlbumObject struct {
	ID     string
	Name   string
	Artist string
	Image  string
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
		ID     string
		Images []ImageItem
		Name   string
	}
	type ArtistsResponse struct {
		Items []ArtistResponseItem
	}
	type SpotifyResponse struct {
		Artists ArtistsResponse
	}
	search, err := spot.Get("search?type=artist&q=%s", nil, artist)
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
	search, err := spot.Get("search?type=album&q=%s artist:%s", nil, album, artist)
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
