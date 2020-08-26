package routes

import (
	"encoding/json"
	"fmt"
	"github.com/chebyrash/promise"
	"github.com/musicorum-app/resource-manager/redis"
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
		//var result *structs.TrackResponse
		//result = new(structs.TrackResponse)
		//
		//search := database.FindAlbum(track.Name, track.Artist)
		// TODO: finish

	})
}
