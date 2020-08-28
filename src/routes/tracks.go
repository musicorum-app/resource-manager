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

type TrackRequestItem struct {
	Name   string  `json:"name"`
	Artist string  `json:"artist"`
	Album  *string `json:"album"`
}

type TracksRequest struct {
	Tracks []TrackRequestItem `json:"tracks"`
}

func TracksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data TracksRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	var promises []*promise.Promise

	for _, track := range data.Tracks {
		promises = append(promises, handleTrack(track))
	}

	results, err := promise.All(promises...).Await()

	fmt.Print(results)

	jsonResult, err := json.Marshal(results)
	utils.FailOnError(err)

	io.WriteString(w, string(jsonResult))
}

func handleTrack(track TrackRequestItem) *promise.Promise {
	return promise.New(func(resolveFinal func(interface{}), reject func(error)) {

		albumHash := ""
		if track.Album != nil {
			albumHash = *track.Album
		}

		hash := utils.HashTrack(track.Name, track.Artist, albumHash)
		redisSearch := redis.FindTrack(hash)
		if redisSearch != nil {
			resolveFinal(redisSearch)
			return
		}

		var result *structs.TrackResponse
		result = new(structs.TrackResponse)

		search := database.FindTrack(hash)
		if search != nil {
			result.Hash = hash
			result.Name = search.Name
			result.Album = search.Album
			result.Artist = search.Artist
			result.Spotify = search.Spotify
			result.Cover = search.Cover
			result.Deezer = search.Deezer
			result.Musixmatch = search.Musixmatch
			result.Duration = search.Duration
			resolveFinal(result)
			redis.SetTrack(result)
		} else {
			p := promise.New(func(resolve func(interface{}), reject func(error)) {
				response := api.SearchTrack(track.Name, albumHash, track.Artist)
				fmt.Print(response)
				if response == nil {
					resolve(nil)
					resolveFinal(nil)
					return
				}

				result.Hash = hash
				result.Name = response.Name
				result.Album = response.Album
				result.Artist = response.Artist
				result.Spotify = &response.ID
				result.Cover = response.Image
				result.Duration = &response.Duration
				resolveFinal(result)
				go func() {
					database.SetTrack(result)
				}()
				fmt.Println("RESULT THIG HERE")
				resolve(response)
				resolveFinal(result)
				redis.SetTrack(result)
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
