package routes

import (
	"encoding/json"
	"fmt"
	"github.com/chebyrash/promise"
	"github.com/musicorum-app/resource-manager/api"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/queue"
	"github.com/musicorum-app/resource-manager/redis"
	"github.com/musicorum-app/resource-manager/structs"
	"github.com/musicorum-app/resource-manager/utils"
	"io"
	"net/http"
)

type AlbumRequestItem struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type AlbumsRequest struct {
	Albums []AlbumRequestItem `json:"albums"`
}

func AlbumsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data AlbumsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	var promises []*promise.Promise

	for _, album := range data.Albums {
		promises = append(promises, handleAlbum(album))
	}

	results, err := promise.All(promises...).Await()

	fmt.Print(results)

	jsonResult, err := json.Marshal(results)
	utils.FailOnError(err)

	io.WriteString(w, string(jsonResult))
}

func handleAlbum(album AlbumRequestItem) *promise.Promise {
	return promise.New(func(resolveFinal func(interface{}), reject func(error)) {
		fmt.Println("Searching for " + album.Name)
		hash := utils.HashAlbum(album.Name, album.Artist)
		redisSearch := redis.FindAlbum(hash)
		if redisSearch != nil {
			resolveFinal(redisSearch)
			return
		}
		var result *structs.AlbumResponse
		result = new(structs.AlbumResponse)

		search := database.FindAlbum(album.Name, album.Artist)
		if search != nil {
			result.Hash = hash
			result.Name = search.Name
			result.Artist = search.Artist
			result.Cover = search.Cover
			result.Spotify = search.Spotify
			resolveFinal(result)
			redis.SetAlbum(result)
		} else {
			p := promise.New(func(resolve func(interface{}), reject func(error)) {
				response := api.SearchAlbum(album.Name, album.Artist)
				fmt.Print(response)
				if response == nil {
					resolve(nil)
					resolveFinal(nil)
					return
				}
				result.Hash = hash
				result.Name = response.Name
				result.Artist = response.Artist
				result.Cover = response.Image
				result.Spotify = response.ID
				resolveFinal(result)
				go func() {
					database.InsertAlbum(result)
				}()
				fmt.Println("RESULT THIG HERE")
				resolve(response)
				resolveFinal(result)
				redis.SetAlbum(result)
			})
			action := queue.AddItem("spotify", p)
			response, err := action.Await()
			if err != nil {
				resolveFinal(nil)
				return
			}
			if response == nil {
				resolveFinal(nil)
				return
			}
		}
	})
}
